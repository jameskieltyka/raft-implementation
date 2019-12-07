package raftserver

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/jkieltyka/raft-implementation/pkg/election"
	"github.com/jkieltyka/raft-implementation/pkg/state"
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
	State      *state.State
	Role       string
	RoleChange chan string
	Heartbeat  chan bool
	LeaderID   string
}

func NewServer() *RaftServer {
	return &RaftServer{
		State: &state.State{
			Log: make([]state.Log, 0, 1),
		},
		Role:       "follower",
		RoleChange: make(chan string, 5),
		Heartbeat:  make(chan bool, 5),
	}
}

func (r *RaftServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	raft.RegisterRaftNodeServer(grpcServer, r)
	go r.MonitorStateChange()
	grpcServer.Serve(lis)
	return nil
}

func (r *RaftServer) RequestVote(ctx context.Context, req *raft.VoteRequest) (*raft.VoteResponse, error) {
	reply, err := election.VoteReply(req, r.State)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "an error occurred %v", err.Error())
	}
	return &reply, nil
}

func (r *RaftServer) AppendEntries(ctx context.Context, req *raft.EntryData) (*raft.EntryResults, error) {
	if (req.Term < r.State.CurrentTerm) || (r.State.LastApplied != req.PrevLogIndex) {
		return &raft.EntryResults{
			Term:    req.Term,
			Success: false,
		}, nil
	}

	if req.Term >= r.State.CurrentTerm && req.LeaderID != r.LeaderID {
		r.LeaderID = req.LeaderID
		r.State.CurrentTerm = req.Term
		r.Role = "follower"
		r.RoleChange <- "follower"
	}

	if len(req.Entries) == 0 {
		r.Heartbeat <- true
	}

	return &raft.EntryResults{}, nil
}

func (r *RaftServer) MonitorStateChange() {
	for {
		select {
		case s := <-r.RoleChange:
			fmt.Printf("changing role to %s, term %d\n", s, r.State.CurrentTerm)
		}
	}
}
