package main

import (
	"fmt"

	"github.com/donniet/ilclient"
)

func main() {
	client := ilclient.New()

	cam, err := client.NewComponent("camera")

	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Printf("client: %v", cam)
}
