GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOFMT = $(GOCMD) fmt
GOVET = $(GOCMD) vet
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

BASENAME=virgo4-collections-ws

build: darwin web

all: darwin linux web

linux-full: linux web

darwin-full: darwin web

web:
	mkdir -p bin/
	cd frontend && yarn install && yarn build
	rm -rf bin/public
	mv frontend/dist bin/public

darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -a -o bin/$(BASENAME).darwin backend/*.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -o bin/$(BASENAME).linux backend/*.go

clean:
	$(GOCLEAN) backend/
	rm -rf bin

fmt:
	cd backend; $(GOFMT)

vet:
	cd backend; $(GOVET)

dep:
	cd frontend && yarn upgrade
	$(GOGET) -u ./backend/...
	$(GOMOD) tidy
	$(GOMOD) verify

check:
	go install honnef.co/go/tools/backend/staticcheck
	$(HOME)/go/bin/staticcheck -checks all,-S1002,-ST1003 backend/*.go
	go install golang.org/x/tools/go/analysis/passes/shadow/backend/shadow
	$(GOVET) -vettool=$(HOME)/go/bin/shadow ./backend/...
