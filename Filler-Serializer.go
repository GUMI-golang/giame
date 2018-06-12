package giame

import (
	"image"
	"image/png"
	"bytes"
	"image/draw"
)

type RawFiller struct {
	Type   FillerType `json:"type"`
	Width  int        `json:"width"`
	Height int        `json:"height"`
	Data   []uint8    `json:"data"`
}
type CompressedRawFiller struct {
	Type   FillerType `json:"type"`
	Width  int        `json:"width"`
	Height int        `json:"height"`
	Data   []uint8    `json:"data"`
}

func (s *RawFiller) Restore() Filler {
	return NewFiller(s.Type, &image.RGBA{
		Rect:image.Rect(0,0,s.Width, s.Height),
		Pix:s.Data,
		Stride: 4 * s.Width,
	})
}
func (s *RawFiller) Compress() (*CompressedRawFiller, error) {
	buf := bytes.NewBuffer(nil)
	err := png.Encode(buf, &image.RGBA{
		Rect:image.Rect(0,0, s.Width, s.Height),
		Pix:s.Data,
		Stride:s.Width * 4,
	})
	if err != nil {
		return nil, err
	}
	return &CompressedRawFiller{
		Type:s.Type,
		Width:s.Width,
		Height:s.Height,
		Data:buf.Bytes(),
	}, nil
}
func (s *CompressedRawFiller) Uncompress() (*RawFiller, error) {
	v, err := png.Decode(bytes.NewReader(s.Data))
	if err != nil {
		return nil, err
	}
	rgba := image.NewRGBA(v.Bounds())
	draw.Draw(rgba, rgba.Rect, v, v.Bounds().Min, draw.Src)
	return &RawFiller{
		Type:s.Type,
		Width:s.Width,
		Height:s.Height,
		Data:rgba.Pix,
	}, nil
}