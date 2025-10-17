docker build -t "spelltext-flyway:latest" -f scripts/db/Dockerfile scripts/db
docker build -t "spelltext-storeserver:latest" -f server/store/Dockerfile .
docker build -t "spelltext-charserver:latest" -f server/character/Dockerfile .
docker build -t "spelltext-inventoryserver:latest" -f server/inventory/Dockerfile .
docker build -t "spelltext-chatserver:latest" -f server/chat/Dockerfile .