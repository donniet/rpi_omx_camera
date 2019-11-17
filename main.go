package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"

	"github.com/donniet/ilclient"
)

func main() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			os.Exit(1)
		}
	}()

	go func() {
		<-sig

		debug.PrintStack()
		panic(fmt.Errorf("interrupted"))
	}()

	client := ilclient.New()

	cam, err := client.NewComponent("image_encode",
		ilclient.CreateFlagDisableAllPorts,
		ilclient.CreateFlagEnableOutputBuffers,
		ilclient.CreateFlagEnableInputBuffers)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	state, err := cam.State()
	if err != nil {
		log.Printf("error getting state: %v", err)
		return
	}
	log.Printf("camera state: %v", state)

	cam.SetState(ilclient.StateIdle)

	state, err = cam.State()
	if err != nil {
		log.Printf("error getting state: %v", err)
		return
	}
	log.Printf("camera state: %v", state)

	enc, err := client.NewComponent("video_encode",
		ilclient.CreateFlagEnableInputBuffers,
		ilclient.CreateFlagEnableOutputBuffers,
		ilclient.CreateFlagDisableAllPorts)
	if err != nil {
		log.Printf("error: %v", err)
	}

	log.Printf("cam: %v", cam)
	log.Printf("enc: %v", enc)
}
