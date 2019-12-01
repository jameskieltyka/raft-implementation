package raftservice

import (
	"fmt"

	"github.com/jkieltyka/raft-implementation/internal/raftclient"
	"github.com/jkieltyka/raft-implementation/internal/raftserver"
	"github.com/jkieltyka/raft-implementation/pkg/backoff"
	"github.com/jkieltyka/raft-implementation/pkg/discovery"
)

func StartService() error {

	//Start RaftServer
	rs := raftserver.NewServer()
	go rs.Start()

	//Get Raft Nodes
	k8sdiscover := discovery.Kubernetes{}
	nodes := k8sdiscover.GetNodes("raft", "default")

	//Create RaftClients
	raftClients := raftclient.ClientList{}
	for _, node := range nodes {
		cl, err := raftclient.CreateClient(node)
		if err != nil {
			return err
		}
		raftClients = append(raftClients, cl)
	}

	//Start heartbeat client
	//TODO this should send out an AppendEntries heartbeat if election state has the current leader id == to this node

	//Start RandomBackoff
	backoffSettings := backoff.NewBackoff(500, 1500)
	backoffTimer := backoffSettings.SetBackoff()

	for {
		select {
		case <-backoffTimer.C:
			//set new election backoff
			backoffTimer = backoffSettings.SetBackoff()
			//Send candidate message to nodes
			accepted := raftClients.RequestVote(rs.State)
			if accepted {
				//TODO send append entries RPC to establish leadership

				fmt.Println("leadership established")
			}

		case <-rs.Heartbeat:
			backoffSettings.ResetBackoff(backoffTimer)
		}

	}

}