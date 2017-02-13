DISTFILE := terraform-provider-ucs
.PHONY: dist

default: test clean build

bootstrap: deps
	cp config.tf{.example,}

deps:
	go get github.com/hashicorp/terraform
	go get gopkg.in/xmlpath.v2

build:
	@go build -o $(DISTFILE) .

run:
	@go run main.go provider.go resource_ucs_service_profile.go config.go

test:
	go test -v . ./ipman ./ucsclient ./ucsclient/ucsinternal

coverage:
	@go test -cover

clean:
	@rm -f ./dist/*

dist:
	# Build for darwin-amd64
	GOOS=darwin GOARCH=amd64 go build -o ./dist/terraform-provider-ucs
	cd dist && tar cf terraform-provider-ucs-osx.tar.gz terraform-provider-ucs
	# Build for linux-amd64
	GOOS=linux GOARCH=amd64 go build -o ./dist/terraform-provider-ucs
	cd dist && tar cf terraform-provider-ucs-linux-amd64.tar.gz terraform-provider-ucs
	# Build for linux-386
	GOOS=linux GOARCH=386 go build -o ./dist/terraform-provider-ucs
	cd dist && tar cf terraform-provider-ucs-linux-i386.tar.gz terraform-provider-ucs
