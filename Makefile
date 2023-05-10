GO:=/usr/bin/go
CONT:=docker
TAG:=generator:latest

run-local:
	$(GO) run ./cmd/generator -devel

ui/node_modules:
	cd ui && npm install

build: generator ui/dist

ui/dist: ui/node_modules build-ui


run-devel: build-ui generator
	./generator -devel

test: build-ui generator
	$(GO) test ./pkg/...
	bash scripts/run_tests.sh


build: build-ui generator

build-ui: ui/node_modules
	cd ui/ && npm run build

generator:
	$(GO) build ./cmd/...

dist: clean
	$(CONT) build . --tag $(TAG)

release: dist
	$(CONT) push $(TAG)

full-clean: clean
	rm -rf ui/node_modules ui/.angular tmp

clean:
	rm -rf generator ui/dist

.PHONY: test run-local run-devel clean full-clean