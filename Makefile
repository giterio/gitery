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

production:
	mkdir -p $(GOBIN)/ && cp ./configs/configs.yaml $(GOBIN)/
	cp ./configs/configs.yaml $(GOBIN)/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME)