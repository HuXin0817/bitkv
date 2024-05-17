package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/HuXin0817/bitkv"
	"github.com/HuXin0817/bitkv/interval/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var serverAddr = flag.String("h", "0.0.0.0:7070", "server address")

func main() {
	flag.Parse()
	logx.Disable()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		if err := bitkv.Merge(); err != nil {
			fmt.Printf("[ERROR] %v\n", err)
		}
		os.Exit(0)
	}()

	c := Config{
		RpcServerConf: zrpc.RpcServerConf{
			ListenOn: *serverAddr,
		},
	}

	ctx := NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterServeServer(grpcServer, NewServeServer(ctx))
	})
	defer s.Stop()

	fmt.Printf("[INFO] start listen on %s...\n", *serverAddr)
	s.Start()
}
