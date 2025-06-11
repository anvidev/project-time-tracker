package main

import (
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	api, err := NewApiContext(ctx)
	if err != nil {
		log.Fatal("server initialization failed", "error", err)
	}

	if err := api.Run(); err != nil {
		log.Fatal("server error:", "error", err)
	}
}
