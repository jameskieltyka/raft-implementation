package raftservice

import (
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

	//Start RandomBackoff
	backoffSettings := backoff.NewBackoff(500, 1000)
	backoffTimer := backoffSettings.SetBackoff()

	for {
		select {
		case <-backoffTimer.C:

		case <-rs.Heartbeat:
			backoffSettings.ResetBackoff(backoffTimer)
		}

	}

}
