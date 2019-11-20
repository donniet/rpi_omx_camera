package main

import (
	"log"

	"github.com/donniet/ilclient"
)

func main() {
	client := ilclient.Get()
	defer client.Close()

	log.Printf("creating component")

	cam, err := client.NewComponent("camera",
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
	f, err := cam.Port(ilclient.CameraCaptureOut).GetVideoPortFormat()

	log.Printf("format error: %v", err)
	log.Printf("formats: %v", f)

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
