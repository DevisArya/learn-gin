PROTO_DIR=internal/proto
OUT_DIR=internal/pb

PROTO_FILES=$(PROTO_DIR)/notification.proto

generate:
	protoc --proto_path=$(PROTO_DIR) \
	       --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	       $(PROTO_FILES)