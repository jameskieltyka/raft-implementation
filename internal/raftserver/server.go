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
			Log:        make([]raft.Entry, 0, 1),
			NextIndex:  make(map[string]uint32),
			MatchIndex: make(map[string]uint32),
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
	if len(req.Entries) == 0 {
		r.HandleHeartbeat(req)
		return &raft.EntryResults{Term: r.State.CurrentTerm, Success: true}, nil
	}

	if (req.Term < r.State.CurrentTerm) || (r.State.LastApplied < req.PrevLogIndex) ||
		r.State.GetLogTerm(req.GetPrevLogTerm()) != req.GetPrevLogTerm() {
		fmt.Println(req.Term, r.State.CurrentTerm, r.State.LastApplied, req.PrevLogIndex, r.State.GetLogTerm(req.GetPrevLogTerm()), req.GetPrevLogTerm())
		fmt.Println("returning failure", req.LeaderID)
		return EntryFailure(r.State.CurrentTerm), nil
	}

	for i, entry := range req.Entries {
		if r.State.LastApplied == req.GetPrevLogIndex()+uint32(i) {
			fmt.Println("adding new entry")
			r.State.Log = append(r.State.Log, *entry)
			r.State.LastApplied = r.State.LastApplied + uint32(i) + 1
		} else {
			currentEntry := r.State.Log[r.State.LastApplied+uint32(i)]
			if currentEntry.Term != entry.Term {
				fmt.Println("replacing old entry")
				r.State.Log[r.State.LastApplied+uint32(i)] = *entry
				r.State.Log = r.State.Log[0 : r.State.LastApplied+uint32(i)]
				r.State.LastApplied = r.State.LastApplied + uint32(i)
			}
		}

	}

	r.State.CommitIndex = MinCommitIndex(req.LeaderCommit, r.State.CommitIndex)

	return &raft.EntryResults{Term: r.State.CurrentTerm, Success: true}, nil
}

func (r *RaftServer) MonitorStateChange() {
	for {
		select {
		case s := <-r.RoleChange:
			fmt.Printf("changing role to %s, term %d\n", s, r.State.CurrentTerm)
		}
	}
}

func MinCommitIndex(leaderCommit, reqCommit uint32) uint32 {
	if leaderCommit > reqCommit {
		return reqCommit
	}
	return leaderCommit
}

func EntryFailure(term uint32) *raft.EntryResults {
	return &raft.EntryResults{
		Term:    term,
		Success: false,
	}
}

func (r *RaftServer) HandleHeartbeat(req *raft.EntryData) {
	if req.Term >= r.State.CurrentTerm && req.LeaderID != r.LeaderID {
		r.LeaderID = req.LeaderID
		r.State.CurrentTerm = req.Term
		r.Role = "follower"
		r.RoleChange <- "follower"
	}

	r.Heartbeat <- true
}
