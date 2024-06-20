all:
	go install ./...


rpc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	mkdir -p pkg/rtreepb
	protoc apis/bobadojo/rtree/v1/rtree.proto \
		--proto_path='apis' \
		--go_opt='module=github.com/bobadojo/stores-server/pkg/rtreepb' \
		--go_out='pkg/rtreepb'
