package raftserver

import (
	"context"
	"flag"
	"net"
	"os"

	"github.com/jkieltyka/raft-implementation/pkg/election"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	raft "github.com/jkieltyka/raft-implementation/raftpb"
)

var (
	port = flag.Int("port", 8000, "The Raft Server Port")
)

type RaftServer struct {
	raft.UnimplementedRaftNodeServer
	Logs      []int
	State     *election.State
	Heartbeat chan bool
}

func NewServer() RaftServer {
	return RaftServer{
		Logs:      make([]int, 0, 1),
		State:     &election.State{},
		Heartbeat: make(chan bool),
	}
}

func (r RaftServer) Start() error {
	lis, err := net.Listen("tcp", ":8000") //fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	raft.RegisterRaftNodeServer(grpcServer, r)
	grpcServer.Serve(lis)
	return nil
}

func (r RaftServer) RequestVote(ctx context.Context, req *raft.VoteRequest) (*raft.VoteResponse, error) {
	reply, err := r.State.VoteReply(req)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "an error occurred %v", err.Error())
	}
	return &reply, nil
}

func (r RaftServer) AppendEntries(ctx context.Context, req *raft.EntryData) (*raft.EntryResults, error) {
	if req.Term > r.State.CurrentTerm && r.State.CurrentLeaderID == os.Getenv("POD_NAME") {
		r.State.CurrentLeaderID = req.LeaderID
		r.State.CurrentTerm = req.Term
	}
	r.Heartbeat <- true
	return &raft.EntryResults{}, nil
}
