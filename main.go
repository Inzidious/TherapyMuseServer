package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	server := newAPIServer(":3000")
	server.Init()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	err := server.Start(ctx)
	defer cancel()

	if err != nil {
		fmt.Println("Failed to start app")
	}
}
