# spelltext

![logo](./docs/spelltext_logo.png)

**(WIP!!)** command-line interface based mmorpg (derives *spell* from mmorpg **spells**, *text* because it's, well, *text-based*). uses [tview](https://github.com/rivo/tview) as the graphical interface. started as a [bachelor's thesis/project](https://efee.etf.unibl.org/vector/zavrsni-radovi/2520) of mine, relying on the Kubernetes orchestration & containerization concepts, as well as high-availability and consistency-related prerequisites.

https://github.com/user-attachments/assets/af9055f2-41c9-4150-abc8-a8ec40877163

## contents
- [spelltext](#spelltext)
  - [contents](#contents)
  - [usage](#usage)
    - [quickstart](#quickstart)
  - [components](#components)

## usage
### quickstart
to get up and running with **spelltext** with *kubernetes*, simply run:
```sh
# rebuild proto (if necessary):
$ cd proto 
$ protoc -I. -I$(dirname "$(which protoc)")/../include --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. $(cat ./files) 
$ cd ..

# deploy cluster:
$ minikube start
$ helm install spelltext k8s/ -f k8s/values.yaml -n spelltext --create-namespace
$ kubectl config set-context --current --namespace=spelltext
$ minikube tunnel # occupies terminal
$ kubectl port-forward pods/spelltext-nats 4222:4222 # also occupies terminal

# run client:
$ cd client 
$ go run client.go
```

or with *docker* via *docker-compose*:
```sh
$ docker compose up --build --force-recreate --no-attach pgadmin
```

for further usage, refer to the [docs/USAGE.md](./docs/USAGE.md) doc, where everything from setting up, configuring services, maintaining and cleaning up is explained.

## components
- `chatserver`: uses fanout MQ (built upon NATS JetStream durable streams)
- `charserver`: servers as a central characters' hub
- `buildserver`: allows players to upgrade their characters via a build system (TODO)
- `storeserver`: a vendor-based marketplace service
- `combatserver`: instances a isolated '1v1' environment for two entities (players, NPCs)
- `gambaserver`: lets players open chests of various tiers by spending currency
- `authserver`: responsible for authentication
- `inventoryserver`: sole purpose of organising (CRUD) items into the character inventories (backpack)
- `progserver`: progre**sss**erver is a bit too much, no? keeps track of player story progression (TODO)
