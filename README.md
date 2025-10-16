# spelltext

![logo](./docs/spelltext_logo.png)

**(WIP!!)** command-line interface based mmorpg (derives *spell* from mmorpg **spells**, *text* because it's, well, *text-based*). uses [tview](https://github.com/rivo/tview) as the graphical interface. started as a [bachelor's thesis/project](https://efee.etf.unibl.org/vector/zavrsni-radovi/2520) of mine, relying on the Kubernetes orchestration & containerization concepts, as well as high-availability and consistency-related prerequisites.

https://github.com/komadiina/spelltext/blob/dev-docker/docs/demo_s.mp4

## contents
- [spelltext](#spelltext)
  - [contents](#contents)
  - [usage](#usage)
  - [components](#components)
  - [todo:](#todo)

## usage
for further usage, refer to the [docs/USAGE.md](./docs/USAGE.md) doc, where everything from setting up, configuring services, maintaining and cleaning up is explained.

## components
- `chatserver`: uses fanout MQ (built upon NATS JetStream durable streams)
- `inventoryserver`: contains a **pgsql** database (will be sharded, to support high-availability and balanced load), detailing player inventory status, currency (coins, boss tokens, etc.)
- `progserver`: progre**sss**erver is a bit too much, no? keeps track of player story progression
- `charserver`: servers as a primary character build service - specializations, talents
- `storeserver`: a vendor-based marketplace service
- `combatserver`: instances a isolated '1v1' environment for two entities (players, NPCs)
- `gambaserver`: lets players open chests of various tiers by spending currency - contacts `inventoryserver` upon ChestOpenEvent 

## todo:
`1st oct, 2025`:
- fix pgpool `pool_passwd` not using `values.yaml:.pgpool.customUsers[.usernames, .passwords]` field 
