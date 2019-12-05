package raftclient

import (
	"context"
	"fmt"
	"math"
	"os"

	"github.com/jkieltyka/raft-implementation/pkg/election"
	raft "github.com/jkieltyka/raft-implementation/raftpb"
	"google.golang.org/grpc"
)

type ClientList []string

func CreateClient(serverAddr string) (raft.RaftNodeClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(serverAddr+":8000", opts...)
	if err != nil {
		return nil, err
	}

	return raft.NewRaftNodeClient(conn), err
}

func (c *ClientList) RequestVote(state *election.State) bool {
	positiveVotes := 1
	vote := &raft.VoteRequest{
		Term:         state.CurrentTerm,
		CandidateID:  os.Getenv("POD_NAME"),
		LastLogIndex: state.LastLogIndex,
		LastLogTerm:  state.LastLogTerm,
	}

	for _, addr := range *c {
		cl, _ := CreateClient(addr)
		response, err := cl.RequestVote(context.Background(), vote)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if response.VoteGranted {
			positiveVotes++
		}
	}

	fmt.Println("votes received: ", positiveVotes)
	if positiveVotes > int(math.Ceil(float64(len(*c))/2.0)) {
		return true
	}
	return false
}

func (c *ClientList) SendHeartbeat(state *election.State) {
	heartbeat := &raft.EntryData{
		Term:         state.CurrentTerm,
		LeaderID:     os.Getenv("POD_NAME"),
		PrevLogIndex: state.LastLogIndex,
		PrevLogTerm:  state.LastLogTerm,
		Entries:      []*raft.Entry{},
		LeaderCommit: state.CommitIndex,
	}

	for _, addr := range *c {
		cl, _ := CreateClient(addr)
		res, err := cl.AppendEntries(context.Background(), heartbeat)
		if err != nil {
			fmt.Println("error :", err.Error())
		}
		fmt.Println("hearbeat response: ", *res)
	}
}
