# spelltext - usage
requirements:
- go `1.25.0` toolchain (possibly compatible with `1.24.x`)
- docker `28.3.2` (tested on build rev. `578ccf6`)
- (opt) minikube `1.35.0` (compatible with `1.33.0`)
- (opt) helm `3.19.0` (goVersion: `1.24.7`)

# contents
- [spelltext - usage](#spelltext---usage)
- [contents](#contents)
- [client](#client)
- [server](#server)
  - [configuration](#configuration)
  - [docker](#docker)
  - [kubernetes](#kubernetes)
  - [example - accessing pods/servers](#example---accessing-podsservers)
  - [example - cleaning up](#example---cleaning-up)

# client
(*as of 25th sept.*): to run a single client, simply run:
```sh
$ set CLIENT_USERNAME=john_doe

$ cd $PROJECT_DIR/client
$ go run client.go --username=$CLIENT_USERNAME
```

# server
## configuration
the spelltext server apps provide a way to be configured via the YAML configuration file, found [here](../config.yml). if you wish to have two separate configurations (e.g. `config.dev.yml`, `config.prod.yml`), you can do so as such:
1) create your configuration file (e.g. `config.new.yml`):
```sh
$ cd $PROJECT_ROOT
$ export $CONFIG_FILE=config.new.yml
$ touch $CONFIG_FILE
```
2) edit the Dockerfiles:
```dockerfile
# change this
ENV CONFIG_FILE=./config.yml

# to this:
ENV CONFIG_FILE=./config.new.yml
```

## docker
using the [composefile](./compose.yml) you can seamlessly deploy a simple service composition:
```sh
# navigate to project root dir
$ cd $PROJECT_ROOT

# start up `docker compose`
$ docker compose up --watch --force-recreate --build 

# ...
$ docker compose logs --follow # or simply -f

# make sure to docker compose down (can cause duplicated stdout if not done)
$ docker compose down -v --remove-orphans
```
  
## kubernetes
deployed using the [provided charts](https://github.com/komadiina/spelltext/tree/main/k8s/) in the repo and [hosted packages](https://github.com/komadiina/spelltext/pkgs/container/spelltext%2Fchatserver).
```sh
# start up minikube (or k3s, microk8s, gke, ...)
$ minikube start

# navigate to project root dir
$ cd $PROJECT_ROOT

$ helm install spelltext k8s/spelltext -f k8s/spelltext/values.yaml -n spelltext --create-namespace 

# (confirm) set current kubectl namespace to 'spelltext'
$ kubectl config set-context --current --namespace=spelltext

# for local development, port-forward nats from inside the cluster, since client requires direct connection
$ kubectl port-forward pods/spelltext-nats 4222:4222 # keep the terminal open

# tunnel to your kubernetes control node (minikube in this example)
$ minikube tunnel # keep the terminal open
```

## example - accessing pods/servers
- via exposing service external-ip and minikube tunnel:
```sh
# start the tunnel
$ minikube tunnel

# select the service you want to dial 
# use chart/templates/TEMPLATE.yaml:.spec.selector.matchLabels.app, e.g. 'chatserver'
$ export SERVICE_NAME=chatserver

$ kubectl get deployments $SERVICE_NAME
> NAME         READY   UP-TO-DATE   AVAILABLE   AGE
> chatserver   2/2     2            2           7m54s

# expose the deployment via a loadbalancer
$ kubectl expose deployment chatserver --type=LoadBalancer --name=chatserver-lb

# external-ip is displayed now, with the help of 'minikube tunnel' command
$ kubectl get svc
> NAME            TYPE           CLUSTER-IP       EXTERNAL-IP     PORT(S)           AGE
> chatserver      ClusterIP      10.106.245.161   <none>          50051/TCP         16m
> chatserver-lb   LoadBalancer   10.108.13.46     10.96.184.178   50051:32671/TCP   4m33s
```

## example - cleaning up
cleanup:
```sh
$ helm uninstall spelltext
$ kubectl delete ns spelltext

# remove persistent volumes via kubectl 
$ kubectl get pv
# WARNING: only use '--all' when it is safe to do so
$ kubectl delete pv <pv-name>

# terminate the nats port-forward terminal
# terminate the minikube tunnel terminal
$ minikube stop
```