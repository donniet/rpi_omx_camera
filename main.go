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

	// err = cam.Port(ilclient.CameraCaptureOut).Enable()
	if e := cam.RequestCallback(ilclient.IndexParamCameraDeviceNumber, true); e != nil {
		log.Fatalf("error: requesting callback: %v", e)
	}
	if e := cam.SetCameraDeviceNumber(0); e != nil {
		log.Fatalf("error: setting device number: %v", e)
	}

	if pd, err := cam.Port(ilclient.CameraCaptureOut).GetPortDefinition(); err != nil {
		log.Fatalf("error: getting port definition: %v", err)
	} else {
		log.Printf("direction: %v, domain: %v, video: %v", pd.Direction, pd.Domain, pd.Video)
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
