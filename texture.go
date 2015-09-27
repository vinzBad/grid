package grid

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Texture struct {
	Id     uint32
	Width  int32
	Height int32
}

func (t Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.Id)	
}

func CreateTexture(file string) (Texture, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return Texture{}, err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return Texture{}, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return Texture{}, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	width, height := int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y)
	var texture uint32

	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		width,
		height,
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return Texture{Id: texture, Width: width, Height: height}, nil
}
