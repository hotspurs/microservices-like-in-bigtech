# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN := $(CURDIR)/bin

# Добавляем bin в текущей директории в PATH при запуске protoc
PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

# Путь до protobuf файлов
PROTO_PATH := $(CURDIR)/api

# Путь до сгенеренных .pb.go файлов
PKG_PROTO_PATH := $(CURDIR)/pkg

# Путь до завендореных protobuf файлов
VENDOR_PROTO_PATH := $(CURDIR)/vendor.protobuf

# устанавливаем необходимые плагины
.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	$(info Installing binary dependencies...)

	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/bufbuild/buf/cmd/buf@v1.32.2
	go install github.com/yoheimuta/protolint/cmd/protolint@latest


# vendor
vendor:	.vendor-reset .vendor-googleapis .vendor-google-protobuf .vendor-protovalidate .vendor-protoc-gen-openapiv2 .vendor-tidy

.vendor-reset:
	rm -rf $(VENDOR_PROTO_PATH)
	mkdir -p $(VENDOR_PROTO_PATH)

.vendor-tidy:
	find $(VENDOR_PROTO_PATH) -type f ! -name "*.proto" -delete
	find $(VENDOR_PROTO_PATH) -empty -type d -delete

# Устанавливаем proto описания google/protobuf
.vendor-google-protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf $(VENDOR_PROTO_PATH)/protobuf &&\
	cd $(VENDOR_PROTO_PATH)/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p $(VENDOR_PROTO_PATH)/google
	mv $(VENDOR_PROTO_PATH)/protobuf/src/google/protobuf $(VENDOR_PROTO_PATH)/google
	rm -rf $(VENDOR_PROTO_PATH)/protobuf

# Устанавливаем proto описания validate
.vendor-protovalidate:
	git clone -b main --single-branch --depth=1 --filter=tree:0 \
		https://github.com/bufbuild/protovalidate $(VENDOR_PROTO_PATH)/protovalidate && \
	cd $(VENDOR_PROTO_PATH)/protovalidate
	git checkout
	mv $(VENDOR_PROTO_PATH)/protovalidate/proto/protovalidate/buf $(VENDOR_PROTO_PATH)
	rm -rf $(VENDOR_PROTO_PATH)/protovalidate

# Устанавливаем proto описания google/api
.vendor-googleapis:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/googleapis/googleapis $(VENDOR_PROTO_PATH)/googleapis &&\
	cd $(VENDOR_PROTO_PATH)/googleapis &&\
	git checkout
	mv $(VENDOR_PROTO_PATH)/googleapis/google $(VENDOR_PROTO_PATH)
	rm -rf $(VENDOR_PROTO_PATH)/googleapis

# Устанавливаем proto описания protoc-gen-openapiv2/options
.vendor-protoc-gen-openapiv2:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway $(VENDOR_PROTO_PATH)/grpc-gateway && \
 	cd $(VENDOR_PROTO_PATH)/grpc-gateway && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p $(VENDOR_PROTO_PATH)/protoc-gen-openapiv2
	mv $(VENDOR_PROTO_PATH)/grpc-gateway/protoc-gen-openapiv2/options $(VENDOR_PROTO_PATH)/protoc-gen-openapiv2
	rm -rf $(VENDOR_PROTO_PATH)/grpc-gateway

# генерация .go файлов с помощью protoc
.protoc-generate:
	$(info run protoc...)

	mkdir -p $(PKG_PROTO_PATH)
	$(PROTOC) -I $(VENDOR_PROTO_PATH) --proto_path=$(CURDIR) \
	--go_out=$(PKG_PROTO_PATH) --go_opt paths=source_relative \
	--go-grpc_out=$(PKG_PROTO_PATH) --go-grpc_opt paths=source_relative \
	--grpc-gateway_out=$(PKG_PROTO_PATH) --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	$(PROTO_PATH)/notes/messages.proto $(PROTO_PATH)/notes/service.proto

	$(PROTOC) -I $(VENDOR_PROTO_PATH) --proto_path=$(CURDIR) \
	--openapiv2_out=. --openapiv2_opt logtostderr=true \
	$(PROTO_PATH)/notes/service.proto

# Генерация .pb файлов с помощью buf
.buf-generate:
	$(info run buf generate...)
	PATH="$(LOCAL_BIN):$(PATH)" $(LOCAL_BIN)/buf generate

# Форматирование protobuf файлов
proto-format:
	$(info run buf format...)
	$(LOCAL_BIN)/buf format -w

# go mod tidy
.tidy:
	GOBIN=$(LOCAL_BIN) go mod tidy

# Генерация кода из protobuf
generate: .bin-deps .buf-generate proto-format .tidy

# Билд приложения
build:
	go build -o $(LOCAL_BIN) ./cmd/notes/client
	go build -o $(LOCAL_BIN) ./cmd/notes/server

# Линтер
lint: .proto-lint

# Линтер proto файлов
.proto-lint:
	$(LOCAL_BIN)/protolint -config_path ./.protolint.yaml ./proto/


# Объявляем, что текущие команды не являются файлами и
# интсрументируем Makefile не искать изменения в файловой системе
.PHONY: \
	.bin-deps \
	.protoc-generate \
	.buf-generate \
	.tidy \
	.vendor-protovalidate \
	.proto-lint \
	proto-format \
	vendor \
	generate \
	build \
	lint
