##
## Code gen
##

VERSION ?= `go run github.com/mdomke/git-semver/v5`

all: generate

gen_src := $(wildcard ../api/*.proto) $(wildcard ../api/*.yaml) Makefile
gen_sum := gen.sum

generate: gen.sum
$(gen_sum): $(gen_src)
	@shasum $(gen_src) | sort -k 2 > $(gen_sum).tmp
	@diff -q $(gen_sum).tmp $(gen_sum) || ( \
	  set -xe; \
	  (set -e; GO111MODULE=on go mod download); \
	  docker run \
	    --user=`id -u` \
	    --volume="$(PWD)/..:/go/src/berty.tech/berty" \
	    --volume="`go env GOPATH`/pkg/mod:/go/pkg/mod" \
	    --workdir="/go/src/berty.tech/berty/docs" \
	    --entrypoint="sh" \
	    --rm \
	    bertytech/buf:1 \
	    -xec 'make generate_local' \
	)
.PHONY: generate

protoc_opts := -I ../api:`go list -m -mod=mod -f {{.Dir}} github.com/grpc-ecosystem/grpc-gateway`/third_party/googleapis:`go list -m -mod=mod -f {{.Dir}} github.com/gogo/protobuf`:/protobuf

generate_local:
	mkdir -p protocol
	buf generate --template ./buf-doc.gen.yaml ../api/protocol/protocoltypes.proto -o protocol
	@# replace multiple empty lines with one
	cat protocol/api.md.tmp | sed '/^$$/N;/^\n$$/D' > protocol/api.md
	rm -f */*.md.tmp
	shasum $(gen_src) | sort -k 2 > $(gen_sum).tmp
	mv $(gen_sum).tmp $(gen_sum)
	mv protocol/api.md apis/protocoltypes.md
	mv protocol/protocoltypes.swagger.json apis/
.PHONY: generate_local

regenerate: gen.clean generate
.PHONY: regenerate

gen.clean:
	rm -f gen.sum $(wildcard */*.md.tmp) $(wildcard */*.swagger.json)
.PHONY: gen.clean

openapi.prepare: gen.sum
	mkdir -p .tmp/openapi
	cat ./apis/protocoltypes.swagger.json   | jq '.info.version="'$(VERSION)'"' > .tmp/openapi/protocoltypes.swagger.json
	cat .tmp/openapi/*.json | jq .info.version
.PHONY: openapi.prepare

BUMP_TOKEN ?=
bump.validate: openapi.prepare
	@# gem install bump-cli
	bump validate --token=$(BUMP_TOKEN) --doc=6eb1bb1e-c65d-4b73-a8c4-0e545742f6db .tmp/openapi/protocoltypes.swagger.json
.PHONY: bump.validate

bump.deploy: bump.validate
	@# gem install bump-cli
	bump deploy --token=$(BUMP_TOKEN) --doc=6eb1bb1e-c65d-4b73-a8c4-0e545742f6db .tmp/openapi/protocoltypes.swagger.json
.PHONY: bump.deploy

openapi.clean:
	rm -rf .tmp/openapi
.PHONY: openapi.clean

clean: gen.clean openapi.clean
.PHONY: clean

tidy:
.PHONY: tidy
