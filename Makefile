.PHONY: docker
docker:
	docker buildx build --builder mybuilder --no-cache --push --platform linux/amd64,linux/arm64 --no-cache -t mitaka8/ps-bot-2:latest .
.PHONY: encoder
encoder:
	(cd encoder; docker buildx build --builder mybuilder --no-cache --push --platform linux/amd64,linux/arm64 --no-cache -t mitaka8/dca-encoder:latest .)
