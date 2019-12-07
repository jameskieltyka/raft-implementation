.PHONY: generate
generate:
	protoc -I raftpb/ raftpb/raft.proto --go_out=plugins=grpc:raftpb

.PHONY: build
build:
	docker build . -t raft
	kind load docker-image raft --name raft

.PHONY: redeploy
redeploy: build
	kubectl scale deployment raft --replicas 0
	sleep 10
	kubectl scale deployment raft --replicas 3