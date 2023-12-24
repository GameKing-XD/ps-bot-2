.PHONY: docker encoder dev-discord dev-web web
docker: web
	docker buildx build --builder mybuilder --no-cache --push --platform linux/amd64,linux/arm64 --no-cache -t mitaka8/ps-bot-2:latest .
encoder:
	(cd encoder; docker buildx build --builder mybuilder --no-cache --push --platform linux/amd64,linux/arm64 --no-cache -t mitaka8/dca-encoder:latest .)
dev-discord:
	CompileDaemon -build "go build -o /tmp/ps-bot-discord ." -command "/tmp/ps-bot-discord discord"

dev-web:
	CompileDaemon -build "go build -o /tmp/ps-bot-web ." -command "/tmp/ps-bot-web web"
web:
	(cd ./web/src/; npm i; npm run build)
