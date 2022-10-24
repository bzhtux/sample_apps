# Sample PostgreSQL App

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

Github Actions automate the build of the sample_apps-postgres app. All images can be found and pull from:

```text
https://github.com/bzhtux/sample_apps/pkgs/container/sample_apps-postgres/versions
```

Pull a specific version like this:

```shell
docker pull ghcr.io/bzhtux/sample_apps-postgres:<version>
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
name: postgresql-service-binding
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
kubectl cluster-info --context postgresql-service-binding
```

### Namespace

Create a new `namespace` to deploy the sample_apps-posstgres and a postgreeSQL DB:

```shell title="Create a new namepsace"
kubectl create namespace pg-app
```

Update kubernetes conntext to use this new namespace:

```shell title="Use this new namespace with the current k8s context"
kubectl config set-context --current --namespace=pg-app
```

### Deploy PostgreSQL using helm

Add bitnami helm repo:

```shell title="Add bitnami help repo"
helm repo add bitnami https://charts.bitnami.com/bitnami
```

Then install PostgreSQL:

```shell title="Install bitnami help chart for postgresql"
helm install pg bitnami/postgresql
```

PostgreSQL can be accessed via port 5432 on the following DNS names from within your cluster:

`pg-postgresql.pg-app.svc.cluster.local - Read/Write connection`

To get the password for "postgres" run:

```shell
export POSTGRES_PASSWORD=$(kubectl get secret --namespace pg-app pg-postgresql -o jsonpath="{.data.postgres-password}" | base64 -d)
```

To connect to your database run the following command:

```shell
kubectl get secret --namespace pg-app pg-postgresql -o jsonpath="{.data.postgres-password}" | base64 -d
kubectl exec -ti pg-postgresql-0 -- psql --host pg-postgresql -U postgres -d postgres -p 5432 -W
```

### PostgreSQL requirements

Prepare PostgreSQL with username, password, database and extension required by the `sample_apps-postgres` application:

```shell
postgres=# CREATE USER sample_user ;
CREATE ROLE
postgres=# ALTER ROLE sample_user WITH PASSWORD 'sample_password' ;
ALTER ROLE
postgres=# CREATE DATABASE sampledb WITH OWNER sample_user ;
CREATE DATABASE
postgres=# GRANT ALL PRIVILEGES ON DATABASE sampledb TO sample_user ;
GRANT
postgres=# CREATE EXTENSION IF NOT EXISTS "uuid-ossp" ;
CREATE EXTENSION
postgres=# \dx
                            List of installed extensions
   Name    | Version |   Schema   |                   Description
-----------+---------+------------+-------------------------------------------------
 plpgsql   | 1.0     | pg_catalog | PL/pgSQL procedural language
 uuid-ossp | 1.1     | public     | generate universally unique identifiers (UUIDs)
(2 rows)
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

```yaml title="PG service manifest"
---
apiVersion: v1
kind: Service
metadata:
  name: pg-app-svc
spec:
  ports:
  - name: pg-app
    port: 8080
    targetPort: 8080
  selector:
    app: pg-app
    app.kubernetes.io/name: pg-app
```

```yaml title="PG ingress manifest"
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: sample-apps-ingress
spec:
  ingressClassName: contour
  rules:
  - host: app-pg.127.0.0.1.nip.io
    http:
      paths:
      - backend: 
          service:
            name: pg-app-svc
            port:
              number: 8080
        pathType: Prefix
        path: /
```

### Define PostgreSQL configuration

Provide the correct informations in the `k8s/01.secret.yaml` file :

```shell
echo -n "sample_user" | base64
echo -n "sample_password" | base64
echo -n "sampledb" | base64
echo -n "5432" |  base64
echo -n "pg-postgresql.pg-app.svc.cluster.local" | base64
echo -n "postgresql" | base64
```

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: postgres
data:
  host: cGctcG9zdGdyZXNxbC5wZy1hcHAuc3ZjLmNsdXN0ZXIubG9jYWw=
  port: NTQzMg==
  username: c2FtcGxlX3VzZXI=
  password: c2FtcGxlX3Bhc3N3b3Jk
  database: c2FtcGxlZGI=
  sslenabled: dHJ1ZQ==
  type: cG9zdGdyZXNxbA==
```

### Deploy in k8s kind

```shell
kubectl create -f k8s/
secret/postgres created
deployment.apps/pg-app created
service/pg-app created
ingress.networking.k8s.io/sample-apps-ingress created
```

Test the deployment:

```shell
curl -sL http://app-pg.127.0.0.1.nip.io/ | jq .
{
  "message": "Alive",
  "status": "Up"
}
```

Update test.sh with the hostname `app-pg.127.0.0.1.nip.io` and run the tests:

```shell
./test.sh
Adding a new Book: The Hitchhiker's Guide to the Galaxy by Douglas Adams
New book has ID 1

Adding twice the same book, expecting a conflict
{
  "data": {
    "ID": 1
  },
  "message": "A Book already exists with title: The Hitchhiker's Guide to the Galaxy",
  "status": "Conflict"
}

Getting the book with ID: 1
{
  "Author": "Douglas Adams",
  "ID": "1",
  "Title": "The Hitchhiker's Guide to the Galaxy"
}

Deleting book with ID 1
{
  "message": "Book with ID 1 was successfuly deleted",
  "status": "Deleted"
}
```

Logs from the container:

```text
pg-app-6c9bdb469d-xbbjw pg-app 2022/10/07 12:38:46 /go/src/github.com/bzhtux/postgres/pkg/postgresql/AddNewBook.go:19 record not found
pg-app-6c9bdb469d-xbbjw pg-app [1.228ms] [rows:0] SELECT * FROM "books" WHERE Title = 'The Hitchhiker\'s Guide to the Galaxy' ORDER BY "books"."id" LIMIT 1
pg-app-6c9bdb469d-xbbjw pg-app [GIN] 2022/10/07 - 12:38:46 | 202 |    6.753476ms |      172.18.0.1 | POST     "/add"
pg-app-6c9bdb469d-xbbjw pg-app [GIN] 2022/10/07 - 12:38:46 | 409 |      867.94µs |      172.18.0.1 | POST     "/add"
pg-app-6c9bdb469d-xbbjw pg-app [GIN] 2022/10/07 - 12:38:46 | 200 |    1.905347ms |      172.18.0.1 | GET      "/get/1"
pg-app-6c9bdb469d-xbbjw pg-app [GIN] 2022/10/07 - 12:38:46 | 200 |     3.80335ms |      172.18.0.1 | DELETE   "/del/1"
```
