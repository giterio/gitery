PROJECTNAME=$(shell basename "$(PWD)")

GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

install:
	go mod download

development:
	go build -o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME) || exit
	cp ./configs/configs.yaml $(GOBIN)/ || exit

start:
	go build -o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME) || exit
	cp ./configs/configs.yaml $(GOBIN)/ || exit
	./bin/gitery -env=development

production:
	GOOS=linux go build -ldflags="-s -w" -o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME) || exit
	cp ./configs/configs.yaml $(GOBIN)/ || exit