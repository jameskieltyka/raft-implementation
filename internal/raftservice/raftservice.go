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
	time.Sleep(5 * time.Second)
	//Start RaftServer
	rs := raftserver.NewServer()
	go rs.Start()

	//Get Raft Nodes
	k8sdiscover := discovery.Kubernetes{}
	// nodes := k8sdiscover.GetNodes("raft", "default")
	nodes := k8sdiscover.GetNodes("raft", "default")
	var raftClients raftclient.ClientList = nodes

	//Start RandomBackoff
	backoffSettings := backoff.NewBackoff(2000, 5000)
	backoffTimer := backoffSettings.SetBackoff()

	heartbeatTimer := time.NewTimer(100 * time.Millisecond)
	heartbeatTimer.Stop()

	for {
		select {
		case <-heartbeatTimer.C:
			if rs.State.CurrentLeaderID == os.Getenv("POD_NAME") {
				raftClients.SendHeartbeat(rs.State)
				rs.Heartbeat <- true
				heartbeatTimer = time.NewTimer(500 * time.Millisecond)
			}
		case <-backoffTimer.C:
			//TODO FIX ISSUE WWHERE BACKOFF IS OCCURING ON LEADER DUE TO NO HEARBEAT
			fmt.Println("starting election")
			//set new election backoff
			backoffTimer = backoffSettings.SetBackoff()
			//Send candidate message to nodes
			accepted := raftClients.RequestVote(rs.State)
			if accepted {
				fmt.Println("leadership established term: ", rs.State.CurrentTerm)
				rs.State.CurrentLeaderID = os.Getenv("POD_NAME")
				raftClients.SendHeartbeat(rs.State)
				heartbeatTimer = time.NewTimer(500 * time.Millisecond)
			}

		case <-rs.Heartbeat:
			backoffSettings.ResetBackoff(backoffTimer)
		}
	}
}
