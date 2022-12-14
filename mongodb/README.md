# Sample MongoDB App

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

Github Actions automate the build of the sample_apps-mongo app. All images can be found and pull from:

```text
https://github.com/bzhtux/sample_apps/pkgs/container/sample_apps-redis/versions
```

```shell
docker pull ghcr.io/bzhtux/sample_apps-mongo:<version>
```

## Test it locally

### Create a kind cluster

```yaml
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: mongodb-service-binding
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

```shell
kubectl cluster-info --context kind-mongodb-service-binding
```

### Namespace

Create a new `namespace` to deploy the sample_apps-mongo :

```shell
kubectl create namespace mongo-app
```

Update kubernetes conntext to use this new namespace:

```shell
kubectl config set-context --current --namespace=mongo-app
```

### Deploy Redis using helm

Add bitnami helm repo:

```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
```

Then install Redis:

```shell
helm install mongodb bitnami/mongodb
```

MongoDB can be accessed on the following DNS names from within your cluster:

* `mongodb.mongo-app.svc.cluster.local`

To get your password run the following command:

```shell
export MONGODB_ROOT_PASSWORD=$(kubectl get secret --namespace mongo-app mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 -d)
```

Write down the MongoDB host from previous output:

```text
[...]
mongodb.mongo.svc.cluster.local
[...]
```

Get the MongoDB password:

```shell
kubectl get secret --namespace mongo-app mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 -d
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
apiVersion: v1
kind: Service
metadata:
  name: mongo-app-svc
spec:
  ports:
  - name: mongo-app
    port: 8080
    targetPort: 8080
  selector:
    app: mongo-app
    app.kubernetes.io/name: mongo-app
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
  - host: gomongo.127.0.0.1.nip.io
    http:
      paths:
      - backend:
          service:
            name: mongo-app-svc
            port:
              number: 8080
        pathType: Prefix
        path: /
```

### Define Redis configuration

Define connection informations and crededentials within the k8s/01.secret.yaml as below:

```shell
echo -n "mongodb.mongo-app.svc.cluster.local" | base64
echo -n "root" | base64
echo -n "${MONGODB_ROOT_PASSWORD}" | base64
echo -n "27017" |  base64
echo -n "mongodb" | base64 -d
```

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: mongo
data:
  host: bW9uZ29kYi5tb25nby5zdmMuY2x1c3Rlci5sb2NhbA==
  port: MjcwMTc=
  username: cm9vdA==
  password: OWxWeHlNOWp3MA==
  type: bW9uZ29kYg==
```

### Deploy in k8s kind

```shell
kubectl create -f k8s/
secret/mongo created
deployment.apps/mongo-app created
service/mongo-app created
ingress.networking.k8s.io/sample-apps-ingress created
```

Test the deployment:

```shell
curl -sL http://gomongo.127.0.0.1.nip.io/ | jq .
{
  "message": "Alive",
  "status": "Up"
}
```

Update  test.sh with the hotname `gomongo.127.0.0.1.nip.io` and run the tests:

```json
curl -sL -X POST -d '{"Title": "Hello world ", "Author":"bzhtux"}' http://gomongo.127.0.0.1.nip.io:8080/add | jq .
{
  "data": {
    "Book Author": "bzhtux",
    "Book title": "Hello world",
    "ID": "63515b321a0c3cb17aa08a5b",
    "result": {
      "InsertedID": "63515b321a0c3cb17aa08a5b"
    }
  },
  "message": "New book added to books' collection",
  "status": "OK"
}
curl -sL -X POST -d '{"Title": "Hello world", "Author":"bzhtux"}' http://127.0.0.1.nip.io:8080/add | jq .
{
  "message": "Book Hello world already exists.",
  "status": "Conflict"
}
```
