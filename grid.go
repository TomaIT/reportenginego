package reportengine

import (
	"github.com/signintech/gopdf"
)

type Grid struct {
	rectangle       Rectangle
	matrix          [][]Component
	minMargin       Margin
	horizontalAlign uint
	verticalAlign   uint
}

func NewGrid(matrix [][]Component, rectangle Rectangle, minMargin Margin,
	horizontalAlign, verticalAlign uint) *Grid {
	grid := new(Grid)
	grid.matrix = matrix
	grid.rectangle = rectangle
	grid.minMargin = minMargin
	grid.horizontalAlign = horizontalAlign
	grid.verticalAlign = verticalAlign
	return grid
}

func (t *Grid) Build(pdf *gopdf.GoPdf, maxWidth float64) {
	maxWidthMatrix := maxWidth - t.minMargin.left - t.minMargin.right
	//Built cells
	for i := range t.matrix {
		colWidth := maxWidthMatrix / float64(len(t.matrix[i]))
		for j := range t.matrix[i] {
			t.matrix[i][j].Build(pdf, colWidth)
		}
	}
	//Adjust Alignment Cells
	lowerY := t.minMargin.top
	lowerX := t.minMargin.left
	sumHeight := 0.0
	for i := 0; i < len(t.matrix); i++ {
		colWidth := maxWidthMatrix / float64(len(t.matrix[i]))
		rowHeight := t.minHeightRow(i)
		for j := range t.matrix[i] {
			t.matrix[i][j].Adjust(pdf, lowerX+float64(j)*colWidth, lowerY, colWidth, rowHeight)
		}
		lowerY += rowHeight
		sumHeight += rowHeight
	}
	t.rectangle.lowerY = 0
	t.rectangle.lowerX = 0
	t.rectangle.width = maxWidth
	t.rectangle.height = sumHeight + t.minMargin.top + t.minMargin.bottom
}
func (t *Grid) Adjust(pdf *gopdf.GoPdf, lowerX, lowerY, width, height float64) {
	if t.MinWidth(pdf) > width || t.MinHeight() > height {
		panic("Width/Height are not sufficient")
	}
	t.rectangle.lowerY = lowerY
	t.rectangle.lowerX = lowerX
	t.rectangle.width = width
	t.rectangle.height = height
	lowerX, lowerY = t.getMatrixStartPosition()
	for i := 0; i < len(t.matrix); i++ {
		tempX := lowerX
		max := 0.0
		for j := 0; j < len(t.matrix[i]); j++ {
			t.matrix[i][j].MoveTo(tempX, lowerY)
			tempX += t.matrix[i][j].GetRectWidth()
			if t.matrix[i][j].GetRectHeight() > max {
				max = t.matrix[i][j].GetRectHeight()
			}
		}
		lowerY += max
	}
}
func (t *Grid) MoveTo(lowerX, lowerY float64) {
	offsetX := t.rectangle.lowerX - lowerX
	offsetY := t.rectangle.lowerY - lowerY
	t.rectangle.lowerY = lowerY
	t.rectangle.lowerX = lowerX
	for i := range t.matrix {
		for j := range t.matrix[i] {
			x, y := t.matrix[i][j].GetRectPosition()
			t.matrix[i][j].MoveTo(x-offsetX, y-offsetY)
		}
	}
}
func (t *Grid) SetVisibilityContainer(isVisible bool) {
	t.rectangle.isVisible = isVisible
}
func (t *Grid) Split(pdf *gopdf.GoPdf, firstHeight float64, splitType int) Component {
	var row int
	marginHeight := t.minMargin.top + t.minMargin.bottom
	minHeight := marginHeight + t.matrix[0][0].GetRectHeight()
	if minHeight > firstHeight { //Too little space
		return nil
	}
	height := marginHeight
	for row = 0; row < len(t.matrix); row++ {
		rowHeight := t.matrix[row][0].GetRectHeight()
		if height+rowHeight > firstHeight {
			break
		}
		height += rowHeight
	}
	x := t.matrix[0:row]
	y := t.matrix[row:len(t.matrix)]
	if splitType == SplitRepeatFirstRow {
		y = append([][]Component{t.matrix[0]}, y...)
	}
	t.matrix = x
	rec := t.rectangle
	t.Build(pdf, rec.width)
	next := NewGrid(y, t.rectangle, t.minMargin, t.horizontalAlign, t.verticalAlign)
	next.Build(pdf, t.rectangle.width)
	return next
}
func (t Grid) MinWidth(*gopdf.GoPdf) float64 {
	return t.matrixWidth() + t.minMargin.left + t.minMargin.right
}
func (t Grid) MinHeight() float64 {
	return t.matrixHeight() + t.minMargin.top + t.minMargin.bottom
}
func (t Grid) FirstVoidSpace() Rectangle {
	panic("Not yet implemented")
}
func (t Grid) Render(pdf *gopdf.GoPdf) {
	t.rectangle.Render(pdf)
	for i := range t.matrix {
		for j := range t.matrix[i] {
			t.matrix[i][j].Render(pdf)
		}
	}
}
func (t Grid) GetRectWidth() float64 {
	return t.rectangle.width
}
func (t Grid) GetRectHeight() float64 {
	return t.rectangle.height
}
func (t Grid) GetRectPosition() (x, y float64) {
	return t.rectangle.lowerX, t.rectangle.lowerY
}
func (t Grid) IsSplittable() bool {
	return true
}

func (t Grid) matrixWidth() float64 {
	max := 0.0
	for _, i := range t.matrix {
		sum := 0.0
		for _, j := range i {
			sum += j.GetRectWidth()
		}
		if sum > max {
			max = sum
		}
	}
	return max
}
func (t Grid) matrixHeight() float64 {
	sum := 0.0
	for _, i := range t.matrix {
		max := 0.0
		for _, j := range i {
			temp := j.GetRectHeight()
			if temp > max {
				max = temp
			}
		}
		sum += max
	}
	return sum
}
func (t Grid) minHeightRow(rowIndex int) float64 {
	max := 0.0
	for _, v := range t.matrix[rowIndex] {
		temp := v.MinHeight()
		if temp > max {
			max = temp
		}
	}
	return max
}
func (t Grid) getMatrixStartPosition() (x float64, y float64) {
	matrixWidth := t.matrixWidth()
	matrixHeight := t.matrixHeight()
	switch t.horizontalAlign {
	case gopdf.Left:
		x = t.rectangle.lowerX + t.minMargin.left
	case gopdf.Right:
		x = t.rectangle.lowerX + t.rectangle.width - t.minMargin.right - matrixWidth
	case gopdf.Center:
		x = t.rectangle.lowerX + (t.rectangle.width-matrixWidth)/2.0
	}
	switch t.verticalAlign {
	case gopdf.Top:
		y = t.rectangle.lowerY + t.minMargin.top
	case gopdf.Middle:
		y = t.rectangle.lowerY + (t.rectangle.height-matrixHeight)/2.0
	case gopdf.Bottom:
		y = t.rectangle.lowerY + t.rectangle.height - t.minMargin.bottom - matrixHeight
	}
	return x, y
}
