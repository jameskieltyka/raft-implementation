package main

import (
	"flag"
	"fmt"
	"net"

	raftserver "github.com/jkieltyka/raft-implementation/internal/raftServer"
	raft "github.com/jkieltyka/raft-implementation/raftpb"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8000, "The Raft Server Port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		return
	}
	grpcServer := grpc.NewServer()
	raft.RegisterRaftNodeServer(grpcServer, &raftserver.NewServer())
	grpcServer.Serve(lis)
}
