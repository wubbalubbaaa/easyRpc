package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wubbalubbaaa/easyRpc"
)

var mux = sync.RWMutex{}
var server = easyRpc.NewServer()
var clientMap = make(map[*easyRpc.Client]struct{})

func main() {
	server.Handler.EnablePool(true)

	server.Handler.Handle("/enter", func(ctx *easyRpc.Context) {
		passwd := ""
		ctx.Bind(&passwd)
		if passwd == "123qwe" {
			// keep client
			mux.Lock()
			clientMap[ctx.Client] = struct{}{}
			mux.Unlock()

			ctx.Write(nil)

			log.Printf("enter success")
		} else {
			log.Printf("enter failed invalid passwd: %v", passwd)
			ctx.Client.Stop()
		}
	})
	// release client
	server.Handler.HandleDisconnected(func(c *easyRpc.Client) {
		mux.Lock()
		delete(clientMap, c)
		mux.Unlock()
	})

	go func() {
		ticker := time.NewTicker(time.Second)
		for i := 0; true; i++ {
			<-ticker.C
			broadcast(i)
		}
	}()

	server.Run("localhost:8888")
}

func broadcast(i int) {
	msg := server.NewMessage(easyRpc.CmdNotify, "/broadcast", fmt.Sprintf("broadcast msg %d", i))
	mux.RLock()
	defer func() {
		mux.RUnlock()
		msg.ReleaseAndPayback(server.Handler)
	}()

	for client := range clientMap {
		msg.Retain()
		client.PushMsg(msg, easyRpc.TimeZero)
	}
}
