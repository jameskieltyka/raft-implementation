package raftclient

import (
	"context"

	"github.com/jkieltyka/raft-implementation/pkg/election"
	raft "github.com/jkieltyka/raft-implementation/raftpb"
	"google.golang.org/grpc"
)

type ClientList []raft.RaftNodeClient

func CreateClient(serverAddr string) (raft.RaftNodeClient, error) {
	conn, err := grpc.Dial(serverAddr)
	if err != nil {
		return nil, err
	}

	return raft.NewRaftNodeClient(conn), err
}

func (c ClientList) RequestVote(state election.State) bool {
	positiveVotes := 0
	vote := &raft.VoteRequest{
		Term:         election.State.CurrentTerm + 1,
		CandidateID:  "asd",
		LastLogIndex: election.State.LastLogIndex,
		LastLogTerm:  election.State.LastLogTerm,
	}

	for _, cl := range c {
		response, err = cl.RequestVote(context.Background(), vote)
		if err != nil {
			continue
		}
		if raft.VoteResponse.VoteGranted {
			positiveVotes++
		}
	}

	if positiveVotes > len(c)/2 {
		return true
	}
	return false
}
