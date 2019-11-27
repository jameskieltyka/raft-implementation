.PHONY: generate
generate:
	protoc -I raftpb/ raftpb/raft.proto --go_out=plugins=grpc:raftpb