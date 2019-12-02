package election

import (
	raft "github.com/jkieltyka/raft-implementation/raftpb"
)

//State - state for current term election info
type State struct {
	CurrentLeaderID string
	VotedForID      *string
	CurrentTerm     uint32
	LastLogIndex    uint32
	LastLogTerm     uint32
	CommitIndex     uint32
}

//UpdateTerm updates internal term number and clears any previous vote information
func (e State) UpdateTerm(newTerm uint32) {
	e.CurrentTerm = newTerm
	e.VotedForID = nil
}

//VoteReply response to vote requests with an approval or deny
func (e State) VoteReply(c *raft.VoteRequest) (raft.VoteResponse, error) {
	if !e.voteDecision(c) {
		denyVote := raft.VoteResponse{
			Term:        e.CurrentTerm,
			VoteGranted: false,
		}
		return denyVote, nil
	}

	acceptVote := raft.VoteResponse{
		Term:        c.Term,
		VoteGranted: true,
	}
	e.UpdateTerm(c.Term)
	return acceptVote, nil

}

func (e State) voteDecision(c *raft.VoteRequest) bool {
	if e.CurrentTerm > c.Term {
		return false
	}

	if (e.VotedForID == nil || *e.VotedForID == c.CandidateID) &&
		c.LastLogIndex >= e.LastLogIndex &&
		c.LastLogTerm >= e.LastLogIndex {
		return true
	}

	return false

}

func strPtr(s string) *string {
	return &s
}
