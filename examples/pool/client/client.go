package main

import (
	"log"
	"net"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
)

func main() {
	pool, err := easyRpc.NewClientPool(func() (net.Conn, error) {
		return net.DialTimeout("tcp", "localhost:8888", time.Second*3)
	}, 5)
	if err != nil {
		panic(err)
	}
	for i := 0; i < pool.Size(); i++ {
		pool.Get(i).Set("user", i)
	}
	defer pool.Stop()

	for i := 0; i < 10; i++ {
		req := "hello"
		rsp := ""
		client := pool.Next()
		err = client.Call("/echo", &req, &rsp, time.Second*5)
		user, _ := pool.Get(i).Get("user")
		if err != nil {
			log.Fatalf("client[%v] Call failed: %v", user, err)
		} else {
			log.Printf("client[%v] Call Response: \"%v\"", user, rsp)
		}
	}
}
