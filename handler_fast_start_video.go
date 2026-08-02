package main

import (
	"fmt"
	"os/exec"
)

func processVideoForFastStart(filepath string) (string, error) {
	outputFilepath := filepath + ".processing"

	cmd := exec.Command(
		"ffmpeg",
		"-i", filepath,
		"-c", "copy",
		"-movflags", "faststart",
		"-f", "mp4",
		outputFilepath,
	)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg fast-start processing failed: $w", err)
	}

	return outputFilepath, nil
}
