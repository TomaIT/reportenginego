package reportengine

import (
	"errors"
	"strconv"
)

type Color struct {
	r uint8
	g uint8
	b uint8
}

func NewColor(hex string) (Color, error) {
	var err error
	var r, g, b uint64
	if len(hex) != 7 {
		return Color{}, errors.New("hex color invalid format: #FFFFFF")
	}
	r, err = strconv.ParseUint(hex[1:3], 16, 8)
	if err != nil {
		return Color{}, err
	}
	g, err = strconv.ParseUint(hex[3:5], 16, 8)
	if err != nil {
		return Color{}, err
	}
	b, err = strconv.ParseUint(hex[5:7], 16, 8)
	if err != nil {
		return Color{}, err
	}
	return Color{r: uint8(r), g: uint8(g), b: uint8(b)}, nil
}

func White() Color {
	return Color{255, 255, 255}
}
func Red() Color {
	return Color{255, 0, 0}
}
func Green() Color {
	return Color{0, 255, 0}
}
func Blu() Color {
	return Color{0, 0, 255}
}
func Black() Color {
	return Color{0, 0, 0}
}
func Yellow() Color {
	return Color{255, 255, 0}
}
