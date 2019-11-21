package main

import (
	"log"

	"github.com/donniet/ilclient"
)

func createCamera(client *ilclient.Client, width uint, height uint, framerate float64) *ilclient.Component {
	cam, err := client.NewComponent("camera",
		ilclient.CreateFlagDisableAllPorts,
		ilclient.CreateFlagEnableOutputBuffers,
		ilclient.CreateFlagEnableInputBuffers)
	if err != nil {
		log.Printf("error: %v", err)
		return cam
	}
	state, err := cam.State()
	if err != nil {
		log.Printf("error getting state: %v", err)
		return cam
	}

	// err = cam.Port(ilclient.CameraCaptureOut).Enable()
	if e := cam.RequestCallback(ilclient.IndexParamCameraDeviceNumber, true); e != nil {
		log.Fatalf("error: requesting callback: %v", e)
	}
	if e := cam.SetCameraDeviceNumber(0); e != nil {
		log.Fatalf("error: setting device number: %v", e)
	}

	camCapture := cam.Port(ilclient.CameraCaptureOut)

	if pd, err := camCapture.GetPortDefinition(); err != nil {
		log.Fatalf("error: getting port definition: %v", err)
	} else {
		log.Printf("direction: %v, domain: %v, video: %v", pd.Direction, pd.Domain, pd.Video)
		pd.Video.Width = width
		pd.Video.Height = height
		pd.Video.Framerate = framerate
		pd.Video.Stride = ilclient.CalculateStride(pd.Video.Width, pd.BufferAlignment)
		pd.Video.Color = ilclient.ColorFormatYUV420PackedPlanar

		err = camCapture.SetPortDefinition(pd)
		if err != nil {
			log.Fatalf("error: setting port definition: %v", err)
		}
	}

	//check port def
	if pd, err := camCapture.GetPortDefinition(); err != nil {
		log.Fatalf("error: getting port definition: %v", err)
	} else {
		log.Printf("direction: %v, domain: %v, video: %v", pd.Direction, pd.Domain, pd.Video)

	}

	f, err := camCapture.GetVideoPortFormat()

	log.Printf("format error: %v", err)
	log.Printf("formats: %v", f)

	if err := camCapture.SetFramerate(framerate); err != nil {
		log.Fatalf("error: setting framerate: %v", err)
	}

	log.Printf("camera state: %v", state)

	// cam.SetState(ilclient.StateIdle)

	return cam
}

func createEncoder(client *ilclient.Client, width uint, height uint, framerate float64, bitrate uint) *ilclient.Component {
	enc, err := client.NewComponent("video_encode",
		ilclient.CreateFlagEnableInputBuffers,
		ilclient.CreateFlagEnableOutputBuffers,
		ilclient.CreateFlagDisableAllPorts)
	if err != nil {
		log.Printf("error: %v", err)
	}

	encOut := enc.Port(ilclient.VideoEncodeCompressedOut)

	if pd, err := encOut.GetPortDefinition(); err != nil {
		log.Fatalf("error: getting port definition: %v", err)
	} else {
		log.Printf("encoder direction: %v, domain: %v, video: %v", pd.Direction, pd.Domain, pd.Video)
		pd.Video.Width = width
		pd.Video.Height = height
		pd.Video.Framerate = framerate
		pd.Video.Stride = ilclient.CalculateStride(pd.Video.Width, pd.BufferAlignment)
		pd.Video.Bitrate = bitrate

		err = encOut.SetPortDefinition(pd)
		if err != nil {
			log.Fatalf("error: setting port definition: %v", err)
		}
	}

	if pd, err := encOut.GetPortDefinition(); err != nil {
		log.Fatalf("error: getting port definition: %v", err)
	} else {
		log.Printf("encoder direction: %v, domain: %v, video: %v", pd.Direction, pd.Domain, pd.Video)
	}

	// cam.SetState(ilclient.StateIdle)

	return enc
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client := ilclient.Get()
	defer client.Close()

	cam := createCamera(client, 1440, 1080, 15.)
	enc := createEncoder(client, 1440, 1080, 15., 17000000)

	tun, err := client.NewTunnel(cam.Port(ilclient.CameraCaptureOut), enc.Port(ilclient.VideoEncodeRawVideoIn))
	if err != nil {
		log.Fatalf("error: create tunnel: %v", err)
	}

	log.Printf("cam: %v", cam)
	log.Printf("enc: %v", enc)
	log.Printf("tun: %v", tun)
}
