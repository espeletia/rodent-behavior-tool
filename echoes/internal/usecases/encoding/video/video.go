package video

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

type VideoMediaEncoder struct {
	ffmpegPath   string
	ffmprobePath string
}

func NewVideoMediaEncoder(ffmpegPath string, ffmprobePath string) *VideoMediaEncoder {
	return &VideoMediaEncoder{
		ffmpegPath:   ffmpegPath,
		ffmprobePath: ffmprobePath,
	}
}

func (v *VideoMediaEncoder) EncodeVideoWith256(ctx context.Context, videoPath string, dir string) (string, error) {
	outputFile := filepath.Join(dir, fmt.Sprintf("clip%d.mp4", time.Now().UnixNano()))

	cmd := exec.Command(
		v.ffmpegPath,
		"-i", videoPath,
		"-c:v", "libx264",
		"-c:a", "aac",
		outputFile,
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run FFMPEG: %v", err)
	}

	return outputFile, nil
}
