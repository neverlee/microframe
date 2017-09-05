package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/neverlee/microframe/config"
	"github.com/neverlee/microframe/client"

	hello "greeter/proto"
)

func main() {
	rcl := client.NewClient(&client.ClientOpts{
		RequestTimeout:  config.Duration(5 * time.Second),
		Retries:         1,
		PoolSize:        2,
		PoolTTL:         config.Duration(60 * time.Second),
		BrokerAddress:   nil,
		RegistryAddress: []string{"127.0.0.1:8500"},
		Selector:        "",
	})
	// Use the generated client stub
	cl := hello.NewSayClient("micro.frame.srv.example", rcl)

	ctx := context.Background()
	// ctx := metadata.NewContext(context.Background(), map[string]string{
	// 	"X-User-Id": "john",
	// 	"X-From-Id": "script",
	// })

	rsp, err := cl.Hello(ctx, &hello.Request{
		Name: "Never Lee",
	})
	fmt.Println(rsp.Msg, err)
}
