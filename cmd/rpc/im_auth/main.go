package main

import (
	"flag"

	"github.com/qingw1230/studyim/internal/rpc/auth"
)

func main() {
	rpcPort := flag.Int("port", 10600, "RpcToken default listen port 10600")
	flag.Parse()
	rpcServer := auth.NewRpcAuthServer(*rpcPort)
	rpcServer.Run()
}
