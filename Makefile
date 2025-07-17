TEST_PATH := ./db/postgres
COVER_DIR := ./coverage
GEN_PB_DIR := ./api/gen/pb
PROTO_DIR := ./api/proto
GOOGLEAPIS_DIR := $(PROTO_DIR)/third_party

all: run

run: docker
	go run cmd/app/main.go -k=3.0

docker:
	docker-compose --env-file .env up -d

.bin-deps:
	$(info Installing binary dependencies...)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	
	$(info Downloading google well-known types...)
	@if [ ! -d "$(GOOGLEAPIS_DIR)" ]; then \
		curl -LO https://github.com/googleapis/googleapis/archive/refs/heads/master.zip; \
		unzip -q master.zip; \
		mkdir -p $(GOOGLEAPIS_DIR); \
		cp -r googleapis-master/google $(GOOGLEAPIS_DIR)/; \
		rm -rf master.zip googleapis-master; \
	fi

.protoc-generate:
	rm -rf $(GEN_PB_DIR)
	mkdir -p $(GEN_PB_DIR)
	cd $(PROTO_DIR) && protoc \
		--proto_path=. \
		--proto_path=third_party \
		--go_out=../../$(GEN_PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=../../$(GEN_PB_DIR) --go-grpc_opt=paths=source_relative \
		frequencies.proto
	go mod tidy

generate: .bin-deps .protoc-generate

tests:
	mkdir -p $(COVER_DIR)
	go test $(TEST_PATH) -v -coverprofile $(COVER_DIR)/cover.out
	go tool cover -html $(COVER_DIR)/cover.out -o $(COVER_DIR)/cover.html && rm $(COVER_DIR)/cover.out

clean:
	rm -rf $(GEN_PB_DIR)
	rm -rf $(GOOGLEAPIS_DIR)
	rm -rf $(COVER_DIR)