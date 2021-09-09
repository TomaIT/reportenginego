package reportengine

import (
	"github.com/signintech/gopdf"
	"sort"
)

type Page struct {
	rectangle Rectangle
	//header    Component
	//footer    Component
	content []Component
}

func NewPage(rectangle Rectangle, header, footer Component) Page {
	page := new(Page)
	page.rectangle = rectangle
	if header != nil {
		page.content = append(page.content, header)
	}
	if footer != nil {
		page.content = append(page.content, footer)
	}
	return *page
}

func (t *Page) getFirstVoidSpace() (lowerX, lowerY, width, height float64) {
	var tempLowerY, tempUpperY float64
	width = t.rectangle.width   // Assumption is immutable
	lowerX = t.rectangle.lowerX // Assumption is immutable

	if len(t.content) <= 0 {
		lowerY = t.rectangle.lowerY
		height = t.rectangle.height
		return lowerX, lowerY, width, height
	}

	sort.Slice(t.content, func(i, j int) bool {
		_, iy := t.content[i].GetRectPosition()
		_, jy := t.content[j].GetRectPosition()
		return iy < jy
	})

	//First Rectangle
	_, tempLowerY = t.content[0].GetRectPosition()
	if t.rectangle.lowerY < tempLowerY {
		lowerY = t.rectangle.lowerY
		height = tempLowerY - t.rectangle.lowerY
		return lowerX, lowerY, width, height
	}

	//Center Rectangle
	for i := 1; i < len(t.content); i++ {
		_, tempUpperY = t.content[i-1].GetRectPosition()
		tempUpperY += t.content[i-1].GetRectHeight()

		_, tempLowerY = t.content[i].GetRectPosition()

		if tempLowerY > tempUpperY {
			lowerY = tempUpperY
			height = tempLowerY - tempUpperY
			return lowerX, lowerY, width, height
		}
	}
	//Last Rectangle
	_, tempUpperY = t.content[len(t.content)-1].GetRectPosition()
	tempUpperY += t.content[len(t.content)-1].GetRectHeight()

	tempLowerY = t.rectangle.lowerY + t.rectangle.height

	if tempUpperY < tempLowerY {
		lowerY = tempUpperY
		height = tempLowerY - tempUpperY
	}
	return lowerX, lowerY, width, height
}

func (t *Page) Render(pdf *gopdf.GoPdf) {
	pdf.AddPage()
	//t.header.Render(pdf)
	for i := range t.content {
		t.content[i].Render(pdf)
	}
	//t.footer.Render(pdf)
}
