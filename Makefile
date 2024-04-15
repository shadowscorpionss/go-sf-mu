# Makefile to run all services from root dir

# microservice list
SERVICES := censors comments news gateway

.PHONY: all $(SERVICES)

all: $(SERVICES)

$(SERVICES):
	@echo "Launching the service $@"
	@(cd ./cmd/$@ && start cmd /c "go run cmd/main.go")

