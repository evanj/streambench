# Makefile for updating protocol buffer definitions and downloading required tools
BUILD_DIR:=build
PROTOC:=$(BUILD_DIR)/protoc
PROTOC_GEN_GO:=$(BUILD_DIR)/protoc-gen-go

messages/messages.pb.go: messages/messages.proto $(PROTOC) $(PROTOC_GEN_GO)
	$(PROTOC) --plugin=$(PROTOC_GEN_GO) --go_out=paths=source_relative:. $<

# download protoc to a temporary tools directory
$(PROTOC): buildtools/getprotoc.go | $(BUILD_DIR)
	go run $< --output=$@

$(PROTOC_GEN_GO): | $(BUILD_DIR)
	go build --mod=readonly -o $@ github.com/golang/protobuf/protoc-gen-go

$(BUILD_DIR):
	mkdir -p $@
