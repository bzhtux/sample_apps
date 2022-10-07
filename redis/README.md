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
echo -n "redis-master.redis-app.svc.cluster.local" | base64
echo -n "default" | base64
echo -n "redswriter" | base64
echo -n "6379" |  base64
echo -n "0" | base64
echo -n "redis-master.redis-app.svc.cluster.local" | base64
```

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: redis
data:
  host: cmVkaXMtbWFzdGVyLnJlZGlzLWFwcC5zdmMuY2x1c3Rlci5sb2NhbA==
  port: NjM3OQ==
  username: ZGVmYXVsdA==
  password: a2d3SGQ0NFB6Yw==
  database: MA==
  sslenabled: ZmFsc2U=
```

### Deploy in k8s kind

```shell
kubectl create -f k8s/
secret/redis created
deployment.apps/redis-app created
service/redis-app created
ingress.networking.k8s.io/sample-apps-ingress created
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

Logs from the container:

```text
redis-app-5c7699b6bd-vclsm redis-app Lauching sample_app-redis v0.0.1...
redis-app-5c7699b6bd-vclsm redis-app 2022/10/07 14:38:34 Setting a new key key1 with value val1
redis-app-5c7699b6bd-vclsm redis-app [GIN] 2022/10/07 - 14:38:34 | 200 |   12.603373ms |      172.18.0.1 | POST     "/add"
redis-app-5c7699b6bd-vclsm redis-app [GIN] 2022/10/07 - 14:38:34 | 409 |     2.00916ms |      172.18.0.1 | POST     "/add"
redis-app-5c7699b6bd-vclsm redis-app [GIN] 2022/10/07 - 14:38:34 | 200 |    1.650912ms |      172.18.0.1 | GET      "/get/key1"
redis-app-5c7699b6bd-vclsm redis-app [GIN] 2022/10/07 - 14:38:34 | 200 |    2.493383ms |      172.18.0.1 | DELETE   "/del/key1"
```
