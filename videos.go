package main

import (
	"bytes"
	"encoding/json"
	"math"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {
	const (
		sixteenByNine = 1.77777777778
		nineBySixteen = 0.5625
		epsilon       = 1e-3
	)

	args := []string{"-v", "error", "-print_format", "json", "-show_streams", filePath}

	cmd := exec.Command("ffprobe", args...)
	buf := bytes.Buffer{}
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	metadata := metadata{}
	err = json.Unmarshal(buf.Bytes(), &metadata)
	if err != nil {
		return "", err
	}

	width := metadata.Streams[0].Width
	height := metadata.Streams[0].Height

	ratio := float64(width) / float64(height)
	switch {
	case math.Abs(sixteenByNine-ratio) <= epsilon:
		return "16:9", nil
	case math.Abs(nineBySixteen-ratio) <= epsilon:
		return "9:16", nil
	default:
		return "other", nil
	}
}
