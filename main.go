package main

import (
	"flag"

	"github.com/jkieltyka/raft-implementation/internal/raftservice"
)

func main() {
	flag.Parse()
	raftservice.StartService()
}
