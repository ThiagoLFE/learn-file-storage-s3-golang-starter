package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os/exec"
)

type GetAspectRatio struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

const (
	Landscape string = "landscape" // 16:9
	Portrait  string = "portrait"  // 9:16
	Other     string = "other"     // any else
)

const aspectRatioTolerance = 0.3

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-select_streams", "v:0", "-show_streams", filePath)

	data, err := cmd.Output()
	if err != nil {
		log.Printf("fail to run ffprobe command : %v", err)
		return "", fmt.Errorf("internal server error")
	}

	var result GetAspectRatio
	if err := json.Unmarshal(data, &result); err != nil {
		log.Printf("fail to read data from ffprobe command : %v", err)
		return "", fmt.Errorf("internal server error")
	}

	if len(result.Streams) < 1 {
		return "", fmt.Errorf("no metadata founded on this file")
	}

	stream := result.Streams[0]
	return classifyAspectRatio(stream.Width, stream.Height), nil
}

func classifyAspectRatio(width, height int) string {
	if width <= 0 || height <= 0 {
		return Other
	}

	ratio := float64(width) / float64(height)

	if isCloseAspectRatio(ratio, 16.0/9.0) {
		return Landscape
	}
	if isCloseAspectRatio(ratio, 9.0/16.0) {
		return Portrait
	}
	return Other
}

func isCloseAspectRatio(actual, expected float64) bool {
	return math.Abs(actual-expected)/expected <= aspectRatioTolerance
}
