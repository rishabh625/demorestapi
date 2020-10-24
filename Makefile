APP		    := fyndtest
GO_VARS ?=
GO ?= go

help:
	@echo "  clean        clean binary file created"	
	@echo "  build        to build the main binary for current platform"
	@echo "  test         to run unittests"
	@echo "  build-docker to create docker image"
	@echo "  run-docker   to run application in containerized way with postgresql"

build-docker: 
	@echo "Building Docker image"
	docker build -t $(APP) .

build:
	go build -o=$(APP) $(GOPATH)/src/$(APP)/cmd/main/

test:
	$(GO_VARS) $(GO) test -v $(GOPATH)/src/$(APP)/tests

clean:
	rm -f $(APP)

run-docker:
	docker-compose up --build  -d 
