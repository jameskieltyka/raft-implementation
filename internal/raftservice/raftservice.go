package raftservice

import (
	"fmt"
	"os"
	"time"

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
			fmt.Println(err.Error())
			continue
		}
		raftClients = append(raftClients, cl)
	}

	//Start RandomBackoff
	backoffSettings := backoff.NewBackoff(5000, 15000)
	backoffTimer := backoffSettings.SetBackoff()

	heartbeatTimer := time.NewTimer(100 * time.Millisecond)
	heartbeatTimer.Stop()

	fmt.Println("finished init")

	for {
		select {
		case <-heartbeatTimer.C:
			fmt.Println("sending heartbeat")
			if rs.State.CurrentLeaderID == os.Getenv("POD_NAME") {
				raftClients.SendHeartbeat(rs.State)
				heartbeatTimer.Reset(100 * time.Millisecond)
			}
		case <-backoffTimer.C:
			fmt.Println("starting election")
			//set new election backoff
			backoffTimer = backoffSettings.SetBackoff()
			//Send candidate message to nodes
			accepted := raftClients.RequestVote(rs.State)
			if accepted {
				//TODO send append entries RPC to establish leadership
				fmt.Println("leadership established")
				raftClients.SendHeartbeat(rs.State)
				heartbeatTimer = time.NewTimer(100 * time.Millisecond)
			}

		case <-rs.Heartbeat:
			fmt.Println("heartbeat received")
			backoffSettings.ResetBackoff(backoffTimer)
		}
	}
}
