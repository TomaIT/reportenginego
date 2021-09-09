package reportengine

import "github.com/signintech/gopdf"

type Component interface {
	Build(pdf *gopdf.GoPdf, maxWidth float64)
	Adjust(pdf *gopdf.GoPdf, lowerX, lowerY, width, height float64)
	MoveTo(lowerX, lowerY float64)
	MinWidth(pdf *gopdf.GoPdf) float64
	MinHeight() float64
	FirstVoidSpace() Rectangle
	Render(pdf *gopdf.GoPdf)
	GetRectWidth() float64
	GetRectHeight() float64
	GetRectPosition() (x, y float64)
	SetVisibilityContainer(isVisible bool)
	IsSplittable() bool
	Split(pdf *gopdf.GoPdf, firstHeight float64, splitType int) Component
}
