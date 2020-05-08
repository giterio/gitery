PROJECTNAME=$(shell basename "$(PWD)")

GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

install:
	go mod download

development:
	mkdir -p $(GOBIN)/ && cp ./configs/configs.yaml $(GOBIN)/
	go build -o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME) || exit

start:
	mkdir -p $(GOBIN)/ && cp ./configs/configs.yaml $(GOBIN)/
	go build -o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME) || exit
	./bin/gitery -env=development
