package main

import (
	"fmt"
	"time"

	"github.com/djordjev/go-session-im/imsession"
)

func main() {

	fmt.Println("Starting up example")
	manager := imsession.Initialize(time.Second)

	token := manager.Create("test 1")

	payload, _ := manager.Get(token)
	fmt.Println(payload)

	manager.Update(token, "test 2")

	latest, _ := manager.Get(token)
	fmt.Println(latest)

	time.Sleep(time.Second * 8)

	afterRemoval, err := manager.Get(token)
	fmt.Println(afterRemoval, err)

	manager.Stop()
}
