package reportengine

import (
	"github.com/signintech/gopdf"
	"regexp"
	"strings"
)

type CellTextArea struct {
	horizontalAlign uint
	verticalAlign   uint
	rectangle       Rectangle
	underline       bool
	fontSize        int
	fontFamily      string
	color           Color
	originalValue   string
	minMarginText   Margin
	cellsText       []CellText
	cellsTextMerged []CellText
}

func NewCellTextArea(horizontalAlign uint, verticalAlign uint, value string, underline bool,
	fontFamily string, fontSize int, color Color, minMarginText Margin, rectangle Rectangle) *CellTextArea {
	cta := new(CellTextArea)
	cta.horizontalAlign = horizontalAlign
	cta.verticalAlign = verticalAlign
	cta.underline = underline
	cta.fontSize = fontSize
	cta.color = color
	cta.originalValue = value
	cta.minMarginText = minMarginText
	cta.fontFamily = fontFamily
	regexSpaces := regexp.MustCompile(`\s+`)
	str := strings.TrimSpace(regexSpaces.ReplaceAllString(value, " "))
	cta.cellsText = make([]CellText, 0)
	for _, val := range strings.Split(str, " ") {
		margin := gopdf.ContentObjCalTextHeight(fontSize) * 0.58
		if underline {
			margin = margin - float64(fontSize)*UnderlineWidthFactor - UnderlineMargin
		}
		if margin < 0 {
			margin = 0
		}
		margin /= 2.0
		ct := *NewCellText(horizontalAlign, verticalAlign, val, underline, fontFamily,
			fontSize, color, NewVerticalMargin(margin), rectangle)
		cta.cellsText = append(cta.cellsText, ct)
	}
	cta.rectangle = rectangle
	return cta
}

func (t *CellTextArea) Build(pdf *gopdf.GoPdf, maxWidth float64) {
	var i, j int
	//To reset Shorten execution for cellText after that Build is called
	for i = 0; i < len(t.cellsText); i++ {
		t.cellsText[i].toOriginal()
	}
	t.cellsTextMerged = make([]CellText, 0)
	for i = 0; i < len(t.cellsText); i++ {
		nct := t.cellsText[i]
		for j = i + 1; j < len(t.cellsText); j++ {
			mergedWidth := nct.MinWidth(pdf) + t.cellsText[j].MinWidth(pdf) +
				t.minMarginText.left + t.minMarginText.right + Width(pdf, t.fontFamily, t.fontSize, " ")
			if mergedWidth <= maxWidth { //Merge
				nct = merge(nct, t.cellsText[j], " ")
			} else {
				break
			}
		}
		t.cellsTextMerged = append(t.cellsTextMerged, nct)

		i = j - 1
	}
	startY := t.minMarginText.top
	for i := 0; i < len(t.cellsTextMerged); i++ {
		t.cellsTextMerged[i].Build(pdf, maxWidth-t.minMarginText.left-t.minMarginText.right)
		if i > 0 {
			startY += t.cellsTextMerged[i-1].rectangle.height
		}
		t.cellsTextMerged[i].rectangle.lowerY = startY
		t.cellsTextMerged[i].rectangle.lowerX = t.minMarginText.left
	}
	t.rectangle.width = maxWidth
	t.rectangle.height = t.MinHeight()
	t.rectangle.lowerX = 0
	t.rectangle.lowerY = 0
}
func (t *CellTextArea) Adjust(pdf *gopdf.GoPdf, lowerX, lowerY, width, height float64) {
	if t.MinWidth(pdf) > width || t.MinHeight() > height {
		panic("Width/Height are not sufficient")
	}
	t.rectangle.lowerX = lowerX
	t.rectangle.lowerY = lowerY
	t.rectangle.width = width
	t.rectangle.height = height
	x, y := t.getCellTextStartPosition(pdf)
	w := t.cellWidth(pdf)
	for i := 0; i < len(t.cellsTextMerged); i++ {
		h := t.cellsTextMerged[i].MinHeight()
		t.cellsTextMerged[i].Adjust(pdf, x, y, w, h)
		y += h
	}
}
func (t *CellTextArea) MoveTo(lowerX, lowerY float64) {
	offsetX := t.rectangle.lowerX - lowerX
	offsetY := t.rectangle.lowerY - lowerY
	t.rectangle.lowerY = lowerY
	t.rectangle.lowerX = lowerX
	for i := range t.cellsTextMerged {
		x := t.cellsTextMerged[i].rectangle.lowerX - offsetX
		y := t.cellsTextMerged[i].rectangle.lowerY - offsetY
		t.cellsTextMerged[i].MoveTo(x, y)
	}
}
func (t *CellTextArea) SetVisibilityContainer(isVisible bool) {
	t.rectangle.isVisible = isVisible
}
func (t *CellTextArea) Split(*gopdf.GoPdf, float64, int) Component {
	return nil
}
func (t CellTextArea) MinHeight() float64 {
	tot := t.minMarginText.top + t.minMarginText.bottom
	for _, v := range t.cellsTextMerged {
		tot += v.MinHeight()
	}
	return tot
}
func (t CellTextArea) MinWidth(pdf *gopdf.GoPdf) float64 {
	max := 0.0
	for _, v := range t.cellsTextMerged {
		temp := v.MinWidth(pdf)
		if temp > max {
			max = temp
		}
	}
	return max + t.minMarginText.left + t.minMarginText.right
}
func (t CellTextArea) Render(pdf *gopdf.GoPdf) {
	t.rectangle.Render(pdf)
	for _, v := range t.cellsTextMerged {
		//Not rendering cellText rectangle
		v.renderTokens(pdf)
	}
}
func (t CellTextArea) FirstVoidSpace() Rectangle {
	panic("Not implemented")
}
func (t CellTextArea) GetRectWidth() float64 {
	return t.rectangle.width
}
func (t CellTextArea) GetRectHeight() float64 {
	return t.rectangle.height
}
func (t CellTextArea) GetRectPosition() (x, y float64) {
	return t.rectangle.lowerX, t.rectangle.lowerY
}
func (t CellTextArea) IsSplittable() bool {
	return false
}

func (t CellTextArea) getCellTextStartPosition(pdf *gopdf.GoPdf) (x float64, y float64) {
	cellWidth := t.cellWidth(pdf)
	cellHeight := t.cellHeight()
	switch t.horizontalAlign {
	case gopdf.Left:
		x = t.rectangle.lowerX + t.minMarginText.left
	case gopdf.Right:
		x = t.rectangle.lowerX + t.rectangle.width - t.minMarginText.right - cellWidth
	case gopdf.Center:
		x = t.rectangle.lowerX + (t.rectangle.width-cellWidth)/2.0
	}
	switch t.verticalAlign {
	case gopdf.Top:
		y = t.rectangle.lowerY + t.minMarginText.top
	case gopdf.Middle:
		y = t.rectangle.lowerY + (t.rectangle.height-cellHeight)/2.0
	case gopdf.Bottom:
		y = t.rectangle.lowerY + t.rectangle.height - t.minMarginText.bottom - cellHeight
	}
	return x, y
}
func (t CellTextArea) cellHeight() float64 {
	sum := 0.0
	for _, v := range t.cellsTextMerged {
		sum += v.MinHeight()
	}
	return sum
}
func (t CellTextArea) cellWidth(pdf *gopdf.GoPdf) float64 {
	max := 0.0
	for _, v := range t.cellsTextMerged {
		temp := v.MinWidth(pdf)
		if temp > max {
			max = temp
		}
	}
	return max
}
