# Sample Redis App

## Docker build image

Build a new docker image with the sample redis app:

```shell
docker buildx build . --platform linux/amd64 --tag <IMAG NAME>:<IMAGE TAG>
```

And then push this new image or use a CI system to build and push based on whateveer trigger.

```shell
docker push <IMAG NAME>:<IMAGE TAG>
```

## Test it locally

### Create a kind cluster:

```shell
kind create cluster --name redis
Creating cluster "redis" ...
 ✓ Ensuring node image (kindest/node:v1.25.2) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-redis"
You can now use your cluster with:

kubectl cluster-info --context kind-redis

Not sure what to do next? 😅  Check out https://kind.sigs.k8s.io/docs/user/quick-start/
```

Now create a namespace to deploy redis sample app:

```shell
kubectl create namespace redis
```

Update kubernetes context to use this new namespace:

```shell
kubectl config set-context --current --namespace=redis
```

### Deploy a Redis cluster using helm

Add bitnami help repo and deploy a Redis cluster:

```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install redis bitnami/redis
```

Write down the Redis host from peevious output:

```text
[...]
redis-master.redis.svc.cluster.local
[...]
```

Get the Redis password:

```shell
kubectl get secret --namespace redis redis -o jsonpath="{.data.redis-password}" | base64 -d
```

### Configure sample redis app

Define connection informations within the k8s/01.secret.yaml as below:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: redis
data:
  host: bG9jYWxob3N0
  port: NjM3OQ==
  username: ZGVmYXVsdA==
  password: 
  database: MA==
  sslenabled: dHJ1ZQ==
```

Now the sample redis app is ready to be deployed using the k8s/02.deployment file :

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-app
spec:
  selector:
    matchLabels:
      app: redis-app
      app.kubernetes.io/name: redis-app
  template:
    metadata:
      labels:
        app: redis-app
        app.kubernetes.io/name: redis-app
    spec:
      containers:
      - name: redis-app
        image: bzhtux/redis-app:v0.0.1
        env:
          - name: REDIS_HOST
            valueFrom:
              secretKeyRef:
                name: redis
                key: host
                optional: false
          - name: REDIS_PORT
            valueFrom:
              secretKeyRef:
                name: redis
                key: port
                optional: false
          - name: REDIS_USERNAME
            valueFrom:
              secretKeyRef:
                name: redis
                key: username
                optional: false
          - name: REDIS_PASSWORD
            valueFrom:
              secretKeyRef:
                name: redis
                key: password
                optional: false
          - name: REDIS_DB
            valueFrom:
              secretKeyRef:
                name: redis
                key: database
                optional: false
          - name: REDIS_SSL
            valueFrom:
              secretKeyRef:
                name: redis
                key: sslenabled
                optional: false

```

The secret file will be used as secretKeyRef to provide redis app env vars.

To access the sample_app-redis, use the k8s/03.service.yaml :

```yaml
apiVersion: v1
kind: Service
metadata:
  name: redis-app
spec:
  ports:
  - name: redis-app
    port: 8080
    targetPort: 8080
  selector:
    app: redis-app
    app.kubernetes.io/name: redis-app
  type: LoadBalancer
```

### Deploy redis sample app

```shell
kubectl create -f k8s/
```

The result should look like:

```shell
configmap/redisconfig created
deployment.apps/redis-app created
```

To test redis sample app can connect to redis, tail logs from the app pod:

```shell
kubectl get pods
NAME                         READY   STATUS             RESTARTS     AGE
redis-app-5558b7999c-5jw8k   0/1     CrashLoopBackOff   1 (4s ago)   6s
redis-master-0               1/1     Running            0            5m7s
redis-replicas-0             1/1     Running            0            5m7s
redis-replicas-1             1/1     Running            0            4m6s
redis-replicas-2             1/1     Running            0            3m41s
```

```shell
kubectl logs -f redis-app-5558b7999c-5jw8k
Testing Golang Redis ...
2022/10/04 11:32:17 Setting a new key Hello with value From Tanzu Application Platform !
2022/10/04 11:32:17 Getting key Hello => From Tanzu Application Platform !
```

It works !