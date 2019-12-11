package state

import raft "github.com/jkieltyka/raft-implementation/raftpb"

//State - state for current term election info
type State struct {
	//Persistent State
	CurrentTerm uint32
	VoteForID   *string
	Log         []raft.Entry

	//Volatile State
	CommitIndex uint32
	LastApplied uint32
	VotedForID  *string

	//Volatile Leader State
	NextIndex  map[string]uint32
	MatchIndex map[string]uint32
}

func (s *State) ResetLeaderState(nodes []string) {
	for _, node := range nodes {
		s.NextIndex[node] = s.LastApplied + 1
		s.MatchIndex[node] = s.CommitIndex + 1
	}
}

func (s *State) GetLastLogTerm() uint32 {
	var lastLogTerm uint32 = 0
	if len(s.Log) != 0 {
		lastLogTerm = s.Log[s.LastApplied].Term
	}
	return lastLogTerm
}

func (s *State) GetLogTerm(index uint32) uint32 {
	if len(s.Log) == 0 {
		return 0
	}

	return s.Log[index].Term
}
