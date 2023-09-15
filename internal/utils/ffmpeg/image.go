package ffmpeg

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"

	"github.com/games4l/internal/errors"
	"github.com/games4l/internal/logger"
)

type ImageInfoStream struct {
	Index         int    `json:"index"`
	CodecName     string `json:"codec_name" validate:"required"`
	CodecLongName string `json:"codec_long_name" validate:"required"`
	CodecType     string `json:"codec_type" validate:"required"`
	Width         int    `json:"width" validate:"required,gt=0"`
	Height        int    `json:"height" validate:"required,gt=0"`
	CodedWidth    int    `json:"coded_width"`
	CodedHeight   int    `json:"coded_height"`
	PixFmt        string `json:"pix_fmt"`
}

type ImageInfo struct {
	Streams []ImageInfoStream
}

func (p *Provider) GetImageStreamInfo(buf io.Reader) (*ImageInfo, error) {
	cmd := exec.Command(p.ffprobeBinPath,
		"-v",
		"error",
		"-select_streams",
		"v",
		"-show_streams",
		"-of",
		"json",
		"pipe:",
	)

	out, err := cmd.Output()
	if err != nil {
		logger.Warn("Failed to execute ffprobe process: " + err.Error())
		return nil, errors.ErrInvalidFormMedia
	}

	info := ImageInfo{}

	if err = json.Unmarshal(out, &info); err != nil {
		return nil, errors.ErrInvalidFormMedia
	}
	if err = validate.Struct(&info); err != nil {
		return nil, errors.ErrInvalidFormMedia
	}

	return &info, nil
}

func (p *Provider) ValidateImage(s io.Reader) (io.Reader, error) {
	buf := bytes.NewBuffer([]byte{})
	tee := io.TeeReader(s, buf)

	info, err := p.GetImageStreamInfo(tee)
	if err != nil {
		return nil, err
	}

	if info.Streams == nil || len(info.Streams) != 1 {
		return nil, errors.ErrInvalidFormMedia
	}

	si := info.Streams[0]

	if si.Height > si.Width*2 || si.Width > si.Height*2 {
		return nil, errors.ErrBadSizedFormMedia
	} else if si.CodedHeight > si.CodedWidth*2 || si.CodedWidth > si.CodedHeight*2 {
		return nil, errors.ErrBadSizedFormMedia
	}

	return buf, nil
}
