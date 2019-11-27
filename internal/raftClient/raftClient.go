package raftClient

import (
	raft "github.com/jkieltyka/raft-implementation/raftpb"
	"google.golang.org/grpc"
)

func CreateClient(serverAddr string) (raft.RaftNodeClient, error) {
	conn, err := grpc.Dial(serverAddr)
	if err != nil {
		return nil, err
	}

	return raft.NewRaftNodeClient(conn), err
}
