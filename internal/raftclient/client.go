package raftclient

import (
	"context"
	"fmt"
	"os"

	"github.com/jkieltyka/raft-implementation/pkg/election"
	raft "github.com/jkieltyka/raft-implementation/raftpb"
	"google.golang.org/grpc"
)

type ClientList []raft.RaftNodeClient

func CreateClient(serverAddr string) (raft.RaftNodeClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		return nil, err
	}

	return raft.NewRaftNodeClient(conn), err
}

func (c ClientList) RequestVote(state election.State) bool {
	positiveVotes := 0
	vote := &raft.VoteRequest{
		Term:         state.CurrentTerm + 1,
		CandidateID:  os.Getenv("POD_NAME"),
		LastLogIndex: state.LastLogIndex,
		LastLogTerm:  state.LastLogTerm,
	}

	for _, cl := range c {
		response, err := cl.RequestVote(context.Background(), vote)
		if err != nil {
			continue
		}
		if response.VoteGranted {
			positiveVotes++
		}
	}

	if positiveVotes > len(c)/2 {
		return true
	}
	return false
}

func (c ClientList) SendHeartbeat(state election.State) {
	heartbeat := &raft.EntryData{
		Term:         state.CurrentTerm,
		LeaderID:     os.Getenv("POD_NAME"),
		PrevLogIndex: state.LastLogIndex,
		PrevLogTerm:  state.LastLogTerm,
		Entries:      []*raft.Entry{},
		LeaderCommit: state.CommitIndex,
	}

	for _, cl := range c {
		_, err := cl.AppendEntries(context.Background(), heartbeat)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
