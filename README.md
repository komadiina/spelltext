# spelltext

![logo](./docs/spelltext_logo.png)

command-line interface based mmorpg (derives *spell* from mmorpg **spells**, *text* because it's, well, *text-based*). uses [tview](https://github.com/rivo/tview) as the graphical interface. started as a [bachelor's thesis/project](https://efee.etf.unibl.org/vector/zavrsni-radovi/2520) of mine, relying on the Kubernetes orchestration & containerization concepts, as well as high-availability and consistency-related prerequisites.


## usage
requirements:
- go `1.25.0` toolchain
- docker `28.3.2` (tested on build rev. `578ccf6`)
- minikube `1.35.0` (compatible with `1.33.0`)

### client
(*as of 25th sept.*): to run a single client, simply run:
```sh
$ go run client/client.go --username=$YOUR_USERNAME
```

### server
1. `docker compose`:
using the [composefile](./compose.yml) you can seamlessly deploy a simple service composition
  
2. `minikube`:

## components

![screenshot](./docs/spelltext_diagram.jpg)

- `chatserver`:
  - **global**: uses fanout MQ (built upon NATS JetStream durable streams)
  - **whisper**: sends a forward-proxy request to `chatserver` (destination: `username`), which queries a shared Redis database to identify the `chatserver` containing the `username`-identified client, to route the request on
- `inventoryserver`: contains a **pgsql** database (will be sharded, to support high-availability and balanced load), detailing player inventory status, currency (coins, boss tokens, etc.)
- `progserver`: progre**sss**erver is a bit too much, no? keeps track of player story progression
- `toonserver`: servers as a primary character build service - specializations, talents
- `combatserver`: instances a isolated '1v1' environment for two entities (players, NPCs)
- `gambaserver`: lets players open chests of various tiers by spending currency - contacts `inventoryserver` upon ChestOpenEvent 

## todo:
`9th sep, 2025`:
- charts for k8s deployment
- implement FO MQ for `chatserver::global`