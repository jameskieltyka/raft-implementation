package election

import (
	"github.com/jkieltyka/raft-implementation/pkg/state"
	raft "github.com/jkieltyka/raft-implementation/raftpb"
)

//VoteReply response to vote requests with an approval or deny
func VoteReply(c *raft.VoteRequest, e *state.State) (raft.VoteResponse, error) {
	if !voteDecision(c, e) {
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
	return acceptVote, nil

}

func voteDecision(c *raft.VoteRequest, e *state.State) bool {
	if e.CurrentTerm > c.Term {
		return false
	}

	if (e.VotedForID == nil || *e.VotedForID == c.CandidateID) &&
		c.LastLogIndex >= e.LastApplied &&
		c.LastLogTerm >= e.GetLastLogTerm() {
		e.VotedForID = &c.CandidateID
		return true
	}

	return false

}
