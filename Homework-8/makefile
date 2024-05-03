
ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=test password=test dbname=test host=localhost port=5432 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal/pkg
MOCKGEN_TAG=1.2.0
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: test-migration-up
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down


.PHONY: .generate-mockgen-deps
.generate-mockgen-deps:
ifeq ($(wildcard $(MOCKGEN_BIN)),)
	@GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@$(MOCKGEN_BIN)
endif

.PHONY: .generate-mockgen
.generate-mockgen:
ifeq ($(wildcard $(MOCKGEN_BIN)),)
	PATH="$(LOCAL_BIN):$$PATH" go generate -x -run=mockgen ./...
endif

.PHONY: gofmt
gofmt:
	goimports -l -w $(CURDIR)

.test:
	$(info Running tests...)
	go test ./...

test: .test

generate:
	protoc \
		--proto_path=api/ \
		--go_out=internal/pkg/server \
		--go-grpc_out=internal/pkg/server \
		api/*.proto

