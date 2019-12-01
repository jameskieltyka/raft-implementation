package main

import (
	"flag"
	"fmt"
	"net"

	raftserver "github.com/jkieltyka/raft-implementation/internal/raftServer"
	"github.com/jkieltyka/raft-implementation/pkg/discovery"
	raft "github.com/jkieltyka/raft-implementation/raftpb"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8000, "The Raft Server Port")
)

func main() {
	flag.Parse()

	k8sclient := discovery.Kubernetes{}
	fmt.Println(k8sclient.GetNodes("raft", "default"))
	grpcServer := grpc.NewServer()
	raft.RegisterRaftNodeServer(grpcServer, raftserver.NewServer())
	grpcServer.Serve(lis)
}
