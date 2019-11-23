package election

import (
	raft "github.com/jkieltyka/raft-implementation/raftpb"
)

//State - internal state for current term election info
type State struct {
	CurrentLeaderID string
	VotedForID      *string
	CurrentTerm     uint32
	LastLogIndex    uint32
	LastLogTerm     uint32
}

//UpdateTerm updates internal term number and clears any previous vote information
func (e State) UpdateTerm(newTerm uint32) {
	e.CurrentTerm = newTerm
	e.VotedForID = nil
}

//VoteReply response to vote requests with an approval or deny
func (e State) VoteReply(candidateID string, candidateTerm uint32, lastLogTerm uint32, lastLogIndex uint32) error {
	if !e.voteDecision(candidateID, candidateTerm, lastLogIndex, lastLogTerm) {
		denyVote := raft.VoteResponse{
			Term:        e.CurrentTerm,
			VoteGranted: false,
		}

		return nil
	}

	acceptVote := raft.VoteResponse{
		Term:        candidateTerm,
		VoteGranted: true,
	}
	e.UpdateTerm(candidateTerm)
	return nil

}

func (e State) voteDecision(candidateID string, candidateTerm uint32, lastLogTerm uint32, lastLogIndex uint32) bool {
	if e.CurrentTerm > candidateTerm {
		return false
	}

	if (e.VotedForID == nil || *e.VotedForID == candidateID) &&
		lastLogIndex >= e.LastLogIndex &&
		lastLogTerm >= e.LastLogIndex {
		return true
	}

	return false

}

func strPtr(s string) *string {
	return &s
}
