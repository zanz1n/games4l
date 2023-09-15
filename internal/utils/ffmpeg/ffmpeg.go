package ffmpeg

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Provider struct {
	ffprobeBinPath string
	ffmpegBinPath  string
	tempFilePath   string
}

func NewProvider(tempPath, ffmpegBin, ffprobeBin string) *Provider {
	return &Provider{
		ffprobeBinPath: ffprobeBin,
		ffmpegBinPath:  ffmpegBin,
		tempFilePath:   tempPath,
	}
}
