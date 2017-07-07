package main

import (
	"fmt"
	"pluralsight.com/coordinator"
)

func main() {
	fmt.Println("Starting coordinator")
	ea := coordinator.NewEventAggregator()
	//dc := coordinator.NewDatabaseConsumer(ea)
	ql := coordinator.NewQueueListener(ea)
	go ql.ListenForNewSource()

	fmt.Println("Listening for new sources...")

	<-make(chan struct{})
}
