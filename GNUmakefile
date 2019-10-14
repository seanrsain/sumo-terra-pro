TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=sumologic
DIR=~/.terraform.d/plugins

default: build

build: fmtcheck
	go install

install: fmtcheck
	mkdir -vp $(DIR)
	go build -o $(DIR)/terraform-provider-sumologic

uninstall:
	@rm -vf $(DIR)/terraform-provider-sumologic

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: build install uninstall test testacc fmt fmtcheck