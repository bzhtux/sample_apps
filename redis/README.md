# Sample Redis App

## Docker build image

Build a new docker image locally with the sample redis app:

```shell
docker buildx build . --platform linux/amd64 --tag <IMAG NAME>:<IMAGE TAG>
```

And then push this new image or use a CI system to build and push based on whateveer trigger.

```shell
docker push <IMAG NAME>:<IMAGE TAG>
```

## Out Of the Box images

Github Actions automate the build of the sample_apps-redis app. All images can be found and pull from:

```text
https://github.com/bzhtux/sample_apps/pkgs/container/sample_apps-redis/versions
```

```shell
docker pull ghcr.io/bzhtux/sample_apps-redis:<version>
```

## Test it locally

### Create a kind cluster

```yaml
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
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
EOF
```

Now use this new cluster changing the  kubernetes context as below:

```shell
kubectl cluster-info --context kind-kind
```

### Namespace

Create a new `namespace` to deploy the sample_apps-posstgres and a postgreeSQL DB:

```shell
create namespace redis-app
```

Update kubernetes conntext to use this new namespace:

```shell
kubectl config set-context --current --namespace=redis-app
```

### Deploy Redis using helm

Add bitnami helm repo:

```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
```

Then install Redis:

```shell
helm install redis bitnami/redis
```

Redis can be accessed on the following DNS names from within your cluster:

* `redis-master.redis-app.svc.cluster.local for read/write operations (port 6379)`
* `redis-replicas.redis-app.svc.cluster.local for read-only operations (port 6379)`

To get your password run the following command:

```shell
kubectl get secret --namespace redis-app redis -o jsonpath="{.data.redis-password}" | base64 -d
```

Write down the Redis host from peevious output:

```text
[...]
redis-master.redis.svc.cluster.local
[...]
```

Get the Redis password:

```shell
kubectl get secret --namespace redis-app redis -o jsonpath="{.data.redis-password}" | base64 -d
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

### Ingress usage

The following example creates a simple http service and an Ingress object to route to this services.

```yaml
---
kind: Service
apiVersion: v1
metadata:
  name: redis-app-svc
spec:
  selector:
    app: redis-app
    app.kubernetes.io/name: redis-app
  ports:
  # Default port used by the image
  - port: 8080
```

```yaml
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: sample-apps-ingress
spec:
  ingressClassName: contour
  rules:
  - host: app-redis.127.0.0.1.nip.io
    http:
      paths:
      - backend: 
          service:
            name: redis-app-svc
            port:
              number: 8080
        pathType: Prefix
        path: /
```

### Define Redis configuration

Define connection informations and crededentials within the k8s/01.secret.yaml as below:

```shell
export REDIS_HOST=$(echo -n "redis-master.redis-demo.svc.cluster.local" | base64)
export REDIS_USER=$(echo -n "default" | base64)
export REDIS_PASS=$(kubectl get secret --namespace redis-demo redis -o jsonpath="{.data.redis-password}")
export REDIS_PORT=$(echo -n "6379" |  base64)
export REDIS_DB=$(echo -n "0" | base64)
export REDIS_SSL=$(echo -n false | base64)
export REDIS_TYPE=$(echo -n redis | base64)
```

```yaml
cat <<EOF | kubectl apply -f-
apiVersion: v1
kind: Secret
metadata:
  name: goredis
data:
  host: $REDIS_HOST
  port: $REDIS_PORT
  username: $REDIS_USER
  password: $REDIS_PASS
  database: $REDIS_DB
  sslenabled: $REDIS_SSL
  type: $REDIS_TYPE
```

### Deploy redis sample app

```shell
cat <<EOF | kubectl apply -f-
apiVersion: apps/v1
kind: Deployment
metadata:
  name: goredis
spec:
  strategy:
    type: RollingUpdate
  selector:
    matchLabels:
      app: goredis
      app.kubernetes.io/name: goredis
  template:
    metadata:
      labels:
        app: goredis
        app.kubernetes.io/name: goredis
    spec:
      containers:
      - name: goredis
        # image: ghcr.io/bzhtux/sample_apps-redis:v0.0.5
        image: bzhtux/redis-app:test
        imagePullPolicy: Always
        volumeMounts:
        - name: services-bindings
          mountPath: /bindings
          readOnly: true
        env:
        - name: SERVICE_BINDING_ROOT
          value: /bindings
      volumes:
      - name: services-bindings
        projected:
          sources:
          - secret:
              name: goredis
              items:
              - key: host
                path: redis/host
              - key: port
                path: redis/port
              - key: username
                path: redis/username
              - key: password
                path: redis/password
              - key: database
                path: redis/database
              - key: sslenabled
                path: redis/sslenabled
              - key: type
                path: mongodb/type
EOF
```

```shell
cat <<EOF | kubectl apply -f-
apiVersion: v1
kind: Service
metadata:
  name: goredis
spec:
  ports:
  - name: goredis
    port: 8080
    targetPort: 8080
  selector:
    app: goredis
    app.kubernetes.io/name: goredis
EOF
```

```shell
cat <<EOF | kubectl apply -f-
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: goredis
spec:
  ingressClassName: contour
  rules:
  - host: goredis.127.0.0.1.nip.io
    http:
      paths:
      - backend: 
          service:
            name: goredis
            port:
              number: 8080
        pathType: Prefix
        path: /
EOF
```

Test the deployment:

```shell
curl -sL http://app-redis.127.0.0.1.nip.io/ | jq .
{
  "message": "Alive",
  "status": "Up"
}
```

Update  test.sh with the hotname `app-redis.127.0.0.1.nip.io` and run the tests:

```shell
./test.sh
Adding a new key: key1=val1
{
  "data": {
    "key": "key1",
    "value": "val1"
  },
  "message": "New key has been recorder successfuly",
  "status": "OK"
}

Adding twice the same key, expecting a conflict
{
  "message": "Key already exists: key1",
  "status": "Conflict"
}

Getting the previous key: key1
{
  "data": {
    "key": "key1",
    "value": "val1"
  },
  "message": "Key was found",
  "status": "Ok"
}

Deleting key1
{
  "message": "Key was successfuly deleted: key1",
  "status": "Ok"
}
```

It works !
