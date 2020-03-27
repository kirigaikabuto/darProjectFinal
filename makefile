.PHONY: build
build:
	go build -v
	./DarProject-master -config=config.json
.DEFAULT_GOAL := build