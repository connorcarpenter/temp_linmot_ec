PROTO_SRC := $(shell find proto -name '*.proto')
GEN_DIR := gen

.PHONY: all gen clean

all: gen

$(GEN_DIR):
	mkdir -p $(GEN_DIR)

# Requires protoc, protoc-gen-go, protoc-gen-go-grpc installed
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

gen: $(GEN_DIR)
	protoc \
		-I proto \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_SRC)

clean:
	rm -rf $(GEN_DIR)
