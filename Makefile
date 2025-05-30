.PHONY: run stop libyaml build clean

# Environment variables file
include .env

build:
	docker build --build-arg PORT=$(PORT) -t schola-theologiae .

run: build
	docker run --name schola-theo \
		-e PORT=$(PORT):$(PORT) -e SECRET_KEY_BASE=$(SECRET_KEY_BASE) \
		-p $(PORT) schola-theologiae

stop: 
	docker stop schola-theo
	docker container rm schola-theo

libyaml:
	docker build -f Dockerfile.libyaml -t libyaml-builder .
	docker create --name extract libyaml-builder
	docker cp extract:/libyaml.tar.gz ./libyaml.tar.gz
	docker rm extract
