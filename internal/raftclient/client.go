package raftclient

import (
	"context"
	"fmt"
	"math"
	"os"

	"github.com/jkieltyka/raft-implementation/pkg/state"
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

func (c *ClientList) RequestVote(state *state.State) bool {
	positiveVotes := 1
	var lastLogTerm uint32 = 0
	if len(state.Log) != 0 {
		lastLogTerm = state.Log[state.LastApplied].Term
	}

	vote := &raft.VoteRequest{
		Term:         state.CurrentTerm,
		CandidateID:  os.Getenv("POD_NAME"),
		LastLogIndex: state.LastApplied,
		LastLogTerm:  lastLogTerm,
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

	if positiveVotes > int(math.Ceil(float64(len(*c))/2.0)) {
		return true
	}
	return false
}

func (c *ClientList) SendHeartbeat(state *state.State) {

	heartbeat := &raft.EntryData{
		Term:         state.CurrentTerm,
		LeaderID:     os.Getenv("POD_NAME"),
		PrevLogIndex: state.LastApplied,
		PrevLogTerm:  state.GetLastLogTerm(),
		Entries:      []*raft.Entry{},
		LeaderCommit: state.CommitIndex,
	}

	for _, addr := range *c {
		cl, _ := CreateClient(addr)
		_, err := cl.AppendEntries(context.Background(), heartbeat)
		if err != nil {
			fmt.Println("error :", err.Error())
		}
	}
}

// func (c *ClientList) SendLog(state *state.State) {
// 	heartbeat := &raft.EntryData{
// 		Term:         state.CurrentTerm,
// 		LeaderID:     os.Getenv("POD_NAME"),
// 		PrevLogIndex: state.LastLogIndex,
// 		PrevLogTerm:  state.LastLogTerm,
// 		Entries:      []*raft.Entry{},
// 		LeaderCommit: state.CommitIndex,
// 	}

// 	for _, addr := range *c {
// 		cl, _ := CreateClient(addr)
// 		_, err := cl.AppendEntries(context.Background(), heartbeat)
// 		if err != nil {
// 			fmt.Println("error :", err.Error())
// 		}
// 	}
// }
