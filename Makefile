.PHONY: generate
generate:
	protoc -I raftpb/ raftpb/raft.proto --go_out=plugins=grpc:raftpb

.PHONY: build
build:
	docker build . -t raft
	kind load docker-image raft --name raft