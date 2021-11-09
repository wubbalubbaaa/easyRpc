package main

import (
	"log"
	"net"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
)

func main() {
	client, err := easyRpc.NewClient(func() (net.Conn, error) {
		addr, err := net.ResolveUnixAddr("unix", "bench.unixsock")
		if err != nil {
			return nil, err
		}
		return net.DialUnix("unix", nil, addr)
	})
	if err != nil {
		panic(err)
	}

	client.Run()
	defer client.Stop()

	req := "hello"
	rsp := ""
	err = client.Call("/echo", &req, &rsp, time.Second*5)
	if err != nil {
		log.Fatalf("Call failed: %v", err)
	} else {
		log.Printf("Call Response: \"%v\"", rsp)
	}
}
