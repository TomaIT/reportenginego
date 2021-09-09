package reportengine

import (
	"bytes"
	"encoding/base64"
	"github.com/signintech/gopdf"
	"image"
	"log"
	"strings"
)

type CellImage struct {
	rectangle       Rectangle
	horizontalAlign uint
	verticalAlign   uint
	valueBase64     string
	dpi             float64
	minMarginImg    Margin

	imgBytes       []byte
	imgPixelWidth  int
	imgPixelHeight int
	imgType        string
}

func NewCellImage(horizontalAlign, verticalAlign uint, minMarginImg Margin, rectangle Rectangle, valueBase64 string, dpi float64) *CellImage {
	var err error
	ci := new(CellImage)
	ci.rectangle = rectangle
	ci.horizontalAlign = horizontalAlign
	ci.verticalAlign = verticalAlign
	ci.dpi = dpi
	ci.minMarginImg = minMarginImg

	err = ci.setValue(valueBase64)
	if err != nil {
		err = ci.setValue(ImgUnsupportedFormat)
		if err != nil {
			panic(err.Error())
		}
		ci.dpi = 450
	}

	return ci
}

func (t *CellImage) Build(pdf *gopdf.GoPdf, maxWidth float64) {
	if maxWidth < t.MinWidth(pdf) {
		panic("Width is not sufficient")
	}
	t.rectangle.width = maxWidth
	t.rectangle.height = t.MinHeight()
	t.rectangle.lowerX = 0
	t.rectangle.lowerY = 0
}
func (t *CellImage) Adjust(pdf *gopdf.GoPdf, lowerX, lowerY, width, height float64) {
	if t.MinWidth(pdf) > width || t.MinHeight() > height {
		panic("Width/Height are not sufficient")
	}
	t.rectangle.lowerX = lowerX
	t.rectangle.lowerY = lowerY
	t.rectangle.width = width
	t.rectangle.height = height
}
func (t *CellImage) MoveTo(lowerX, lowerY float64) {
	t.rectangle.lowerX = lowerX
	t.rectangle.lowerY = lowerY
}
func (t *CellImage) SetVisibilityContainer(isVisible bool) {
	t.rectangle.isVisible = isVisible
}
func (t *CellImage) Split(*gopdf.GoPdf, float64, int) Component {
	return nil
}
func (t CellImage) MinWidth(*gopdf.GoPdf) float64 {
	return t.imgWidth() + t.minMarginImg.left + t.minMarginImg.right
}
func (t CellImage) MinHeight() float64 {
	return t.imgHeight() + t.minMarginImg.top + t.minMarginImg.bottom
}
func (t CellImage) Render(pdf *gopdf.GoPdf) {
	t.rectangle.Render(pdf)
	lowerX, upperY := t.getImgStartPosition()
	imgH1, err := gopdf.ImageHolderByBytes(t.imgBytes)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	h := t.imgHeight()
	w := t.imgWidth()
	err = pdf.ImageByHolder(imgH1, lowerX, upperY-h, &gopdf.Rect{H: h, W: w})
	if err != nil {
		log.Fatal(err.Error())
	}
}
func (t CellImage) FirstVoidSpace() Rectangle {
	panic("Not implemented")
}
func (t CellImage) GetRectWidth() float64 {
	return t.rectangle.width
}
func (t CellImage) GetRectHeight() float64 {
	return t.rectangle.height
}
func (t CellImage) GetRectPosition() (x, y float64) {
	return t.rectangle.lowerX, t.rectangle.lowerY
}
func (t CellImage) IsSplittable() bool {
	return false
}

func (t *CellImage) setValue(valueBase64 string) error {
	var err error
	t.valueBase64 = valueBase64
	i := strings.Index(valueBase64, ",")
	if i < 0 {
		t.imgBytes, err = base64.StdEncoding.DecodeString(valueBase64)
	} else {
		t.imgBytes, err = base64.StdEncoding.DecodeString(valueBase64[i+1:])
	}
	if err != nil {
		log.Println("Error Image: ", err)
		return err
	}
	var img image.Image
	img, t.imgType, err = image.Decode(bytes.NewReader(t.imgBytes))
	if err != nil {
		log.Println("Error Image: ", err)
		return err
	}
	t.imgPixelWidth = img.Bounds().Max.X
	t.imgPixelHeight = img.Bounds().Max.Y
	return nil
}
func (t CellImage) imgWidth() float64 {
	return 2.54 * float64(t.imgPixelWidth) / t.dpi * 28.3
}
func (t CellImage) imgHeight() float64 {
	return 2.54 * float64(t.imgPixelHeight) / t.dpi * 28.3
}
func (t CellImage) getImgStartPosition() (x float64, y float64) {
	imgWidth := t.imgWidth()
	imgHeight := t.imgHeight()
	switch t.horizontalAlign {
	case gopdf.Left:
		x = t.rectangle.lowerX + t.minMarginImg.left
	case gopdf.Right:
		x = t.rectangle.lowerX + t.rectangle.width - t.minMarginImg.right - imgWidth
	case gopdf.Center:
		x = t.rectangle.lowerX + (t.rectangle.width-imgWidth)/2.0
	}
	switch t.verticalAlign {
	case gopdf.Top:
		y = t.rectangle.lowerY + imgHeight + t.minMarginImg.top
	case gopdf.Middle:
		y = t.rectangle.lowerY + t.rectangle.height - (t.rectangle.height-imgHeight)/2.0
	case gopdf.Bottom:
		y = t.rectangle.lowerY + t.rectangle.height - t.minMarginImg.bottom
	}
	return x, y
}
