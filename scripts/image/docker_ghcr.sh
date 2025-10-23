

docker build -t "ghcr.io/komadiina/spelltext/flyway:latest" -f scripts/db/Dockerfile scripts/db
docker build -t "ghcr.io/komadiina/spelltext/storeserver:latest" -f server/store/Dockerfile .
docker build -t "ghcr.io/komadiina/spelltext/charserver:latest" -f server/character/Dockerfile .
docker build -t "ghcr.io/komadiina/spelltext/inventoryserver:latest" -f server/inventory/Dockerfile .
docker build -t "ghcr.io/komadiina/spelltext/chatserver:latest" -f server/chat/Dockerfile .
docker build -t "ghcr.io/komadiina/spelltext/gambaserver:latest" -f server/gamba/Dockerfile .
docker build -t "ghcr.io/komadiina/spelltext/authserver:latest" -f server/auth/Dockerfile .
docker build -t "ghcr.io/komadiina/spelltext/buildserver:latest" -f server/build/Dockerfile .
docker build -t "ghcr.io/komadiina/spelltext/combatserver:latest" -f server/combat/Dockerfile .

docker push "ghcr.io/komadiina/spelltext/flyway:latest"
docker push "ghcr.io/komadiina/spelltext/storeserver:latest"
docker push "ghcr.io/komadiina/spelltext/charserver:latest"
docker push "ghcr.io/komadiina/spelltext/inventoryserver:latest"
docker push "ghcr.io/komadiina/spelltext/chatserver:latest"
docker push "ghcr.io/komadiina/spelltext/gambaserver:latest"
docker push "ghcr.io/komadiina/spelltext/authserver:latest"
docker push "ghcr.io/komadiina/spelltext/buildserver:latest"
docker push "ghcr.io/komadiina/spelltext/combatserver:latest"