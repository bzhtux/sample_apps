#!/usr/bin/env bash

set -euo pipefail

NAMESPACE="mongo-app"
APP_NAME="mongo-app"
# SVC_NAME represent the bitnami package to install e.g mongodb
SVC_NAME="mongodb"
# example for hub.docker.io: 
# ${docker_username}/${image_name}
# => bzhtux/myapp 
DOCKER_IMG="bzhtux/mongo-app"
DOCKER_TAG="v0.0.4"

tearDown(){
    kind delete cluster --name=${SVC_NAME}-service-binding
}

# shellcheck source=/dev/null
source ./helper.sh

# trap "tearDown" SIGTERM
trap "tearDown" EXIT

echo -e "\033[32m*** Build a new docker image\033[0m"
docker buildx build . --platform linux/amd64 --tag "${DOCKER_IMG}:${DOCKER_TAG}"

echo -e "\033[32m*** Push new docker image\033[0m"
docker push "${DOCKER_IMG}:${DOCKER_TAG}"

echo -e "\033[32m*** Create a kind cluster\033[0m"

cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: ${SVC_NAME}-service-binding
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
- role: worker
EOF

echo -e  "\033[32m*** Switch kubectl context\033[0m"
kubectl cluster-info --context kind-${SVC_NAME}-service-binding

echo -e "\033[32m*** Create a dedicated namespace\033[0m"
kubectl create namespace ${NAMESPACE}

echo -e "\033[32m*** Use this namespace with current context\033[0m"
kubectl config set-context --current --namespace=${NAMESPACE}

echo -e "\033[32m*** Use bitnami repo to deploy PostgSQL\033[0m"
helm repo add bitnami https://charts.bitnami.com/bitnami

echo -e "\033[32m*** Install MongoDB\033[0m"
helm install "${SVC_NAME}" bitnami/"${SVC_NAME}"

echo -e "\033[32m*** Get Hostname and Password\033[0m"
# MONGO_HOST=${pg-MongoDB.pg-app.svc.cluster.local}
MONGO_HOST="${SVC_NAME}.${NAMESPACE}.svc.cluster.local"
MONGO_PASS=$(kubectl get secret --namespace ${NAMESPACE} ${SVC_NAME} -o jsonpath="{.data.${SVC_NAME}-root-password}")
echo "> Done"

echo -e "\033[32m*** Deploy Contour components\033[0m"
kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

echo -e "\033[32m*** Patch Contour for kind specific to forward the hostPorts to the ingress controller, set taint tolerations and schedule it to the custom labelled node\033[0m"
kubectl patch daemonsets -n projectcontour envoy -p '{"spec":{"template":{"spec":{"nodeSelector":{"ingress-ready":"true"},"tolerations":[{"key":"node-role.kubernetes.io/control-plane","operator":"Equal","effect":"NoSchedule"},{"key":"node-role.kubernetes.io/master","operator":"Equal","effect":"NoSchedule"}]}}}}'

echo -e "\033[32m*** Prepare k8s manifests\033[0m"
MONGOPASS=${MONGO_PASS}
MONGOUSER=$(echo -n "root" | base64)
MONGOHOST=$(echo -n ${MONGO_HOST} | base64)
MONGOPORT=$(echo -n 27017 | base64)
MONGOTYPE=$(echo -n mongodb | base64)

echo -e "\033[32m*** Create Secret\033[0m"
cat <<EOF | kubectl apply -f -
---
apiVersion: v1
kind: Secret
metadata:
  name: ${APP_NAME}
data:
  host: $MONGOHOST
  port: $MONGOPORT
  username: $MONGOUSER
  password: $MONGOPASS
  type: $MONGOTYPE
EOF

SVC_POD=$(GetPod $SVC_NAME)
if ! IsPodRunning "${SVC_POD}";then
  WaitForPod "${SVC_POD}" 10
fi

echo -e "\033[32m*** Create Deployment\033[0m"
# kubectl apply -f k8s/02.deployment.yaml
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${APP_NAME}
spec:
  selector:
    matchLabels:
      app: ${APP_NAME}
      app.kubernetes.io/name: ${APP_NAME}
  template:
    metadata:
      labels:
        app: ${APP_NAME}
        app.kubernetes.io/name: ${APP_NAME}
    spec:
      containers:
      - name: ${APP_NAME}
        image: bzhtux/mongo-app:v0.0.4
        imagePullPolicy: Always
        volumeMounts:
        - name: services-bindings
          mountPath: "/bindings"
          readOnly: true
        env:
          - name: SERVICE_BINDING_ROOT
            value: "/bindings"
      - name: sidecar-container
        image: busybox
        command: ["sh","-c","sleep 360000"]
        volumeMounts:
        - name: services-bindings
          mountPath: "/bindings"
          readOnly: true
        env:
          - name: SERVICE_BINDING_ROOT
            value: "/bindings"
      volumes:
      - name: services-bindings
        projected:
          sources:
          - secret:
              name: ${APP_NAME}
              items:
                - key: host
                  path: mongodb/host
                - key: port
                  path: mongodb/port
                - key: username
                  path: mongodb/username
                - key: password
                  path: mongodb/password
                - key: type
                  path: mongodb/type
EOF

echo -e "\033[32m*** Create Service\033[0m"
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: ${APP_NAME}-svc
spec:
  ports:
  - name: ${APP_NAME}
    port: 8080
    targetPort: 8080
  selector:
    app: ${APP_NAME}
    app.kubernetes.io/name: ${APP_NAME}
EOF

echo -e "\033[32m*** Create Ingress\033[0m"
cat <<EOF | kubectl apply -f -
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: sample-apps-ingress
spec:
  ingressClassName: contour
  rules:
  - host: gomongo.127.0.0.1.nip.io
    http:
      paths:
      - backend: 
          service:
            name: ${APP_NAME}-svc
            port:
              number: 8080
        pathType: Prefix
        path: /
EOF



echo -e "\033[32m*** Waiting for MongoDB sample App to be running ...\033[0m"

APP_POD=$(GetPod $APP_NAME)
if ! IsPodRunning "${APP_POD}";then
  WaitForPod "${APP_POD}" 5
fi
WaitForPod "${APP_POD}" 10
echo "> Done"

echo -e "\033[32m*** Running End to End tests...\033[0m"
./test.sh

echo "> Done"
