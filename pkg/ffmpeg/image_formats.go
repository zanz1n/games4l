package ffmpeg

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/games4l/internal/logger"
	"github.com/games4l/pkg/errors"
	"github.com/google/uuid"
)

func (p *Provider) ScaleDownPng(s io.Reader) ([]byte, error) {
	s, err := p.ValidateImage(s)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(p.ffmpegBinPath,
		"-v",
		"error",
		"-i",
		"pipe:",
		"-vf",
		"scale=-1:480",
		"-f",
		"apng",
		"pipe:",
	)

	cmd.Stdin = s

	outBuf, err := cmd.Output()
	if err != nil {
		logger.Error("Failed to get ffmpeg process output buffer: " + err.Error())
		return nil, errors.ErrInvalidFormMedia
	}

	return outBuf, nil
}

func (p *Provider) PngToWebp(b []byte) ([]byte, error) {
	cmd := exec.Command(p.ffmpegBinPath,
		"-v",
		"error",
		"-f",
		"png_pipe",
		"-i",
		"pipe:",
		"-still-picture",
		"1",
		"-f",
		"webp",
		"pipe:",
	)

	cmd.Stdin = bytes.NewReader(b)

	outBuf, err := cmd.Output()
	if err != nil {
		logger.Error("Failed to get ffmpeg process output buffer: " + err.Error())
		return nil, errors.ErrInvalidFormMedia
	}

	return outBuf, nil
}

func (p *Provider) PngToAvif(b []byte) ([]byte, error) {
	fp := path.Join(p.tempFilePath, uuid.NewString())

	cmd := exec.Command(p.ffmpegBinPath,
		"-v",
		"error",
		"-f",
		"png_pipe",
		"-i",
		"pipe:",
		"-map",
		"0",
		"-map",
		"0",
		"-filter:v:1",
		"alphaextract",
		"-frames:v",
		"1",
		"-c:v",
		"libaom-av1",
		"-still-picture",
		"1",
		"-f",
		"avif",
		fp,
	)

	cmd.Stdin = bytes.NewReader(b)

	defer os.Remove(fp)

	if err := cmd.Run(); err != nil {
		logger.Error("Failed to run ffmpeg process: " + err.Error())
		return nil, errors.ErrInvalidFormMedia
	}

	outBuf, err := os.ReadFile(fp)
	if err != nil {
		logger.Error("Failed to open encoded file: " + err.Error())
		return nil, errors.ErrInternalServerError
	}

	return outBuf, nil
}
