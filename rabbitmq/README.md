# Sample RabbitMQ App

For arm64 user see `ARM` section at the end of the document.

## Docker build image

Build a new docker image locally with the sample postgresql app:

```shell title="Build docker image"
docker buildx build . --platform linux/amd64 --tag <IMAG NAME>:<IMAGE TAG>
```

And then push this new image or use a CI system to build and push based on whateveer trigger.

```shell title="Push docker image"
docker push <IMAG NAME>:<IMAGE TAG>
```

## Out Of the Box images

Github Actions automate the build of the sample_apps-rabbitmq app. All images can be found and pull from:

```text
https://github.com/bzhtux/sample_apps/pkgs/container/sample_apps-rabbitmq/versions
```

Pull a specific version like this:

```shell
docker pull ghcr.io/bzhtux/sample_apps-rabbitmq:<version>
```

## Test it locally

### Create a kind cluster

Create a kind cluster with extraPortMappings and node-labels.

* extraPortMappings allow the local host to make requests to the Ingress controller over ports 80/443
* node-labels only allow the ingress controller to run on a specific node(s) matching the label selector

```yaml title="Create a kind cluster"
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: rabbitmq-service-binding
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
```

Now use this new cluster changing the  kubernetes context as below:

```shell title="switch k8s context"
kubectl cluster-info --context kind-rabbitmq-service-binding
```

### Namespace

Create a new `namespace` to deploy the sample_apps-rabbitmq and a RabbitMQ messages queueing system:

```shell title="Create a new namepsace"
kubectl create namespace rabbit-app
```

Update kubernetes conntext to use this new namespace:

```shell title="Use this new namespace with the current k8s context"
kubectl config set-context --current --namespace=rabbit-app
```

### Deploy RabbitMQ using helm

Add bitnami helm repo:

```shell title="Add bitnami help repo"
helm repo add bitnami https://charts.bitnami.com/bitnami
```

Then install RabbitMQ:

```shell title="Install bitnami help chart for postgresql"
helm install rmq bitnami/rabbitmq
```

RabbitMQ can be accessed within the cluster on port 5672 at rmq-rabbitmq.qsa.svc.cluster.local

To get the password for `RabbitMQ` run the following:

```shell title="Get password"
kubectl get secret --namespace rabbit-app rmq-rabbitmq -o jsonpath="{.data.rabbitmq-password}" | base64 -d
```

The user is set to `user`.

To get the `Erlang Cookie`:

```shell title="Get Erlang cookie"
kubectl get secret --namespace rabbit-app rmq-rabbitmq -o jsonpath="{.data.rabbitmq-erlang-cookie}" | base64 -d
```

### Use Contour as the Ingress controller

Deploy Contour components:

```shell
kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
```

Apply kind specific patches to forward the hostPorts to the ingress controller, set taint tolerations and schedule it to the custom labelled node.

```json
{
  "spec": {
    "template": {
      "spec": {
        "nodeSelector": {
          "ingress-ready": "true"
        },
        "tolerations": [
          {
            "key": "node-role.kubernetes.io/control-plane",
            "operator": "Equal",
            "effect": "NoSchedule"
          },
          {
            "key": "node-role.kubernetes.io/master",
            "operator": "Equal",
            "effect": "NoSchedule"
          }
        ]
      }
    }
  }
}
```

```shell
kubectl patch daemonsets -n projectcontour envoy -p '{"spec":{"template":{"spec":{"nodeSelector":{"ingress-ready":"true"},"tolerations":[{"key":"node-role.kubernetes.io/control-plane","operator":"Equal","effect":"NoSchedule"},{"key":"node-role.kubernetes.io/master","operator":"Equal","effect":"NoSchedule"}]}}}}'
```

### Define RabbitMQ configuration

TO create the secret for this RabbitMQ app, run the following commands:

```shell title="create variable for secret"
export RMQ_USER=$(echo -n "user" | base64)
export RMQ_PASSWORD=$(kubectl get secret --namespace rabbit-app rmq-rabbitmq -o jsonpath="{.data.rabbitmq-password}")
export RMQ_HOST=$(echo -n "rmq-rabbitmq.rabbit-app.svc.cluster.local" | base64)
export RMQ_PORT=$(echo -n 5672 | base64)
export RMQ_QUEUE=$(echo -n "vmware" | base64)
export RMQ_TYPE=$(echo -n "rabbitmq" | base64)
```

Once the secret variables are availables, create the secret for the RabbitMQ app:

```yaml title="create k8S secret"
cat <<EOF | kubectl apply -f-
apiVersion: v1
kind: Secret
metadata:
  name: rabbitmq
data:
  host: $RMQ_HOST
  port: $RMQ_PORT
  username: $RMQ_USER
  password: $RMQ_PASSWORD
  queue: $RMQ_QUEUE
  type: $RMQ_TYPE
EOF
```

Then create the deployment as below:

```yaml title="create deployment"
cat <<EOF | kubectl apply -f-
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rabbitmq-app
spec:
  selector:
    matchLabels:
      app: rabbitmq-app
      app.kubernetes.io/name: rabbitmq-app
  template:
    metadata:
      labels:
        app: rabbitmq-app
        app.kubernetes.io/name: rabbitmq-app
    spec:
      containers:
      - name: rabbitmq-app
        # image: ghcr.io/bzhtux/sample_apps-rabbitmq:v0.0.7
        image: bzhtux/rabbitmq-app:test
        imagePullPolicy: Always
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
                name: rabbitmq
                items:
                  - key: host
                    path: rmq/host
                  - key: port
                    path: rmq/port
                  - key: username
                    path: rmq/username
                  - key: password
                    path: rmq/password
                  - key: queue
                    path: rmq/queue
                  - key: type
                    path: rmq/type
EOF
```

Create a service for ingress to access the app:

```yaml title="Create k8s service"
cat <<EOF | kubectl apply -f-
apiVersion: v1
kind: Service
metadata:
  name: rabbitmq-app-svc
spec:
  ports:
  - name: rabbitmq-app
    port: 8080
    targetPort: 8080
  selector:
    app: rabbitmq-app
    app.kubernetes.io/name: rabbitmq-app
EOF
```

And then create an ingress object to access the app from outside the cluster:

```yaml title="create k8s ingress"
cat <<EOF | kubectl apply -f-
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: rabbit-app-ingress
spec:
  ingressClassName: contour
  rules:
  - host: gormq.127.0.0.1.nip.io
    http:
      paths:
      - backend:
          service:
            name: rabbitmq-app-svc
            port:
              number: 8080
        pathType: Prefix
        path: /
EOF
```

## ARM64

If like me you run a M1 chipset or other `arm64` chipset, the `helm` charts used to deploy `RabbitMQ` won't work. So as a workaround on my arm64 laptop, I deployed a docker image with `port mapping` as below:

### Create a local docker network

```shell title="create docker network"
docker network create -d bridge rabbitmq
```

 - `-d bridge` is a driver that let you access the docker image using port mapping
 - `rabbitmq` is the network name

Now run the docker image for `RabbitMQ` arm64 compatible:

```shell title="run docker image"
docker run --network rabbitmq -p 5672:5672 -p 9090:15672 -d --hostname rabbit --name rabbitmq arm64v8/rabbitmq:3.9-management-alpine
```

 - `arm64v8/rabbitmq:3.9-management-alpine` is the arm64 compatible image for RabbitMQ
 - `-p 5672:5672` is the port mapping to access rabbitmq locally on port 5672
 - `-p 9090:15672` is the port mapping to access the rabbitmq management console locally on port 9090 (e.g http://localhost:9090/)
 - `username` is set to `guest` by default
 - `password` is set to `guest` by default

Now fil in the config file with these values as below:

```yaml title="config file"
---
host: "0.0.0.0"
port: 5672
username: "guest"
password: "guest"
queue: "queue-for-the-app"
tls: false
```

`RabbitMQ` is now reachable locally using the `5672` port:

```shell title="test with netcat"
$ nc -z -w 1 0.0.0.0 5672
Connection to 0.0.0.0 port 5672 [tcp/amqp] succeeded!
```
