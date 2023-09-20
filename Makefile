.PHONY: docker
docker:
	docker buildx build --builder mybuilder --no-cache --push --platform linux/amd64,linux/arm64 --no-cache -t mitaka8/ps-bot-2:latest .
.PHONY: encoder
encoder:
	(cd encoder; docker buildx build --builder mybuilder --no-cache --push --platform linux/amd64,linux/arm64 --no-cache -t mitaka8/dca-encoder:latest .)
.PHONY: dev-discord
dev-discord:
	CompileDaemon -build "go build -o /tmp/ps-bot-discord ." -command "/tmp/ps-bot-discord discord"

.PHONY: dev-web
dev-web:
	CompileDaemon -build "go build -o /tmp/ps-bot-web ." -command "/tmp/ps-bot-web web"
