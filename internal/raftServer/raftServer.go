package raftserver

import (
	"context"

	"github.com/jkieltyka/raft-implementation/pkg/election"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	raft "github.com/jkieltyka/raft-implementation/raftpb"
)

type RaftServer struct {
	raft.UnimplementedRaftNodeServer
	logs  []int
	state election.State
}

func NewServer() RaftServer {
	return RaftServer{}
}

func (r RaftServer) RequestVote(ctx context.Context, req *raft.VoteRequest) (*raft.VoteResponse, error) {
	reply, err := r.state.VoteReply(req)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "an error occurred %v", err.Error())
	}
	return &reply, nil
}

func (r RaftServer) AppendEntries(ctx context.Context, req *raft.EntryData) (*raft.EntryResults, error) {
	return nil, nil
}
