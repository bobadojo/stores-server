all:
	go install ./...

rpc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	mkdir -p pkg/rtreepb
	protoc apis/bobadojo/rtree/v1/rtree.proto \
		--proto_path='apis' \
		--go_opt='module=github.com/bobadojo/stores-server/pkg/rtreepb' \
		--go_out='pkg/rtreepb'

artifact-registry:
	gcloud auth configure-docker us-west1-docker.pkg.dev
	docker build . --tag us-west1-docker.pkg.dev/bobadojo/stores/stores:latest
	docker push us-west1-docker.pkg.dev/bobadojo/stores/stores:latest

