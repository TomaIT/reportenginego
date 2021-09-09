package reportengine

import (
	"github.com/signintech/gopdf"
)

type Report struct {
	pdf                      gopdf.GoPdf
	rectangle                Rectangle
	verticalComponentsMargin float64

	headerFP   Component
	headerCP   Component
	headerLP   Component
	footerFP   Component
	footerCP   Component
	footerLP   Component
	contentsFP []Component
	contentsCP []Component
	contentsLP []Component

	pages []Page
}

func NewReport(pageSize gopdf.Rect, marginLeft, marginRight, marginTop, marginBottom float64,
	verticalComponentsMargin float64) *Report {
	report := new(Report)
	report.pdf = gopdf.GoPdf{}
	report.pdf.Start(gopdf.Config{PageSize: pageSize})
	report.rectangle = NewRectangle(0, Solid, 0, White(), White(), true)
	report.rectangle.lowerX = marginLeft
	report.rectangle.lowerY = marginTop
	report.rectangle.width = pageSize.W - marginLeft - marginRight
	report.rectangle.height = pageSize.H - marginTop - marginBottom
	report.verticalComponentsMargin = verticalComponentsMargin
	return report
}

func (t *Report) Build() {
	t.pages = make([]Page, 0)
	if t.headerFP != nil || t.footerFP != nil || len(t.contentsFP) > 0 {
		t.pages = append(t.pages, t.buildSinglePage(t.headerFP, t.footerFP, t.contentsFP))
	}
	if t.headerCP != nil || t.footerCP != nil || len(t.contentsCP) > 0 {
		t.pages = append(t.pages, t.buildMultiplePages(t.headerCP, t.footerCP, t.contentsCP)...)
	}
	if t.headerLP != nil || t.footerLP != nil || len(t.contentsLP) > 0 {
		t.pages = append(t.pages, t.buildSinglePage(t.headerLP, t.footerLP, t.contentsLP))
	}
}

func (t *Report) Render() {
	for i := range t.pages {
		t.pages[i].Render(&t.pdf)
	}
}

func (t *Report) SetHeaderFP(header Component) {
	t.headerFP = header
}
func (t *Report) SetHeaderCP(header Component) {
	t.headerCP = header
}
func (t *Report) SetHeaderLP(header Component) {
	t.headerLP = header
}
func (t *Report) SetFooterFP(footer Component) {
	t.footerFP = footer
}
func (t *Report) SetFooterCP(footer Component) {
	t.footerCP = footer
}
func (t *Report) SetFooterLP(footer Component) {
	t.footerLP = footer
}
func (t *Report) AddContentFP(content Component) {
	t.contentsFP = append(t.contentsFP, content)
}
func (t *Report) AddContentCP(content Component) {
	t.contentsCP = append(t.contentsCP, content)
}
func (t *Report) AddContentLP(content Component) {
	t.contentsLP = append(t.contentsLP, content)
}

func (t Report) GetPdf() gopdf.GoPdf {
	return t.pdf
}

func (t *Report) buildMultiplePages(header, footer Component, contents []Component) []Page {
	pages := make([]Page, 1)
	if header != nil {
		header.Build(&t.pdf, t.rectangle.width)
		header.MoveTo(t.rectangle.lowerX, t.rectangle.lowerY)
	}
	if footer != nil {
		footer.Build(&t.pdf, t.rectangle.width)
		footer.MoveTo(t.rectangle.lowerX, t.rectangle.lowerY+t.rectangle.height-footer.GetRectHeight())
	}
	pages[0] = NewPage(t.rectangle, header, footer)
	if footer != nil {
		lowerX, lowerY, width, height := pages[0].getFirstVoidSpace()
		if height >= t.verticalComponentsMargin/2.0 {
			pages[0].content = append(pages[0].content, getFillComponent(lowerX, lowerY+height-t.verticalComponentsMargin/2.0,
				width, t.verticalComponentsMargin/2.0))
		}
	}

	for i := range contents {
		contents[i].Build(&t.pdf, t.rectangle.width)
	}
	index := 0
	indexPage := 0
	for i := 0; i < MaxIterationInfiniteLoop; i++ {
		lowerX, lowerY, width, height := pages[indexPage].getFirstVoidSpace()
		if index >= len(contents) {
			break
		}
		if height <= 0 { //Page is full
			pages = append(pages, NewPage(t.rectangle, header, footer))
			indexPage++
			continue
		}
		var topMargin = 0.0
		if len(pages[indexPage].content) > 0 {
			topMargin = t.verticalComponentsMargin / 2.0
		}

		if contents[index].GetRectHeight()+topMargin > height { //Too little space for render this component
			if contents[index].IsSplittable() {
				next := contents[index].Split(&t.pdf, height-topMargin, SplitRepeatFirstRow)
				if next == nil { //Impossible split in this space, is too little
					pages[indexPage].content = append(pages[indexPage].content, getFillComponent(lowerX, lowerY, width, height))
					continue
				}
				contents = addSliceElement(contents, next, index+1)
			} else {
				pages[indexPage].content = append(pages[indexPage].content, getFillComponent(lowerX, lowerY, width, height))
				continue
			}
		}
		// Aggiunge spazio vuoto per Vertical Component Margin
		if topMargin > 0.0 {
			pages[indexPage].content = append(pages[indexPage].content, getFillComponent(lowerX, lowerY, width, topMargin))
			lowerY += topMargin
		}
		contents[index].Adjust(&t.pdf, lowerX, lowerY, width, contents[index].GetRectHeight())
		pages[indexPage].content = append(pages[indexPage].content, contents[index])
		index++
	}
	return pages
}

func (t *Report) buildSinglePage(header, footer Component, contents []Component) Page {
	if header != nil {
		header.Build(&t.pdf, t.rectangle.width)
		header.MoveTo(t.rectangle.lowerX, t.rectangle.lowerY)
	}
	if footer != nil {
		footer.Build(&t.pdf, t.rectangle.width)
		footer.MoveTo(t.rectangle.lowerX, t.rectangle.lowerY+t.rectangle.height-footer.GetRectHeight())
	}
	page := NewPage(t.rectangle, header, footer)
	if footer != nil {
		lowerX, lowerY, width, height := page.getFirstVoidSpace()
		if height >= t.verticalComponentsMargin/2.0 {
			page.content = append(page.content, getFillComponent(lowerX, lowerY+height-t.verticalComponentsMargin/2.0,
				width, t.verticalComponentsMargin/2.0))
		}
	}

	for i := range contents {
		contents[i].Build(&t.pdf, t.rectangle.width)
	}
	index := 0
	for i := 0; i < MaxIterationInfiniteLoop; i++ {
		lowerX, lowerY, width, height := page.getFirstVoidSpace()
		if index >= len(contents) || height <= 0 { //Page is full
			break
		}
		var topMargin = 0.0
		if len(page.content) > 0 {
			topMargin = t.verticalComponentsMargin / 2.0
		}

		if contents[index].GetRectHeight()+topMargin > height { //Too little space for render this component
			if contents[index].IsSplittable() {
				next := contents[index].Split(&t.pdf, height-topMargin, SplitNormal)
				if next == nil { //Impossible split in this space, is too little
					page.content = append(page.content, getFillComponent(lowerX, lowerY, width, height))
					continue
				}
				contents = addSliceElement(contents, next, index+1)
			} else {
				page.content = append(page.content, getFillComponent(lowerX, lowerY, width, height))
				continue
			}
		}
		// Aggiunge spazio vuoto per Vertical Component Margin
		if topMargin > 0.0 {
			page.content = append(page.content, getFillComponent(lowerX, lowerY, width, topMargin))
			lowerY += topMargin
		}
		contents[index].Adjust(&t.pdf, lowerX, lowerY, width, contents[index].GetRectHeight())
		page.content = append(page.content, contents[index])
		index++
	}
	return page
}

func getFillComponent(lowerX, lowerY, width, height float64) Component {
	rect := NewRectangle(0, Solid, 0.0, White(), White(), false)
	rect.lowerY = lowerY
	rect.lowerX = lowerX
	rect.width = width
	rect.height = height
	return NewCellText(gopdf.Center, gopdf.Middle, "", false, IconFontFamily, 1, White(), NewMargin(0.0), rect)
}

func addSliceElement(contents []Component, element Component, index int) []Component {
	temp := make([]Component, len(contents)+1)
	for i := 0; i < index; i++ {
		temp[i] = contents[i]
	}
	temp[index] = element
	for i := index + 1; i < len(contents)+1; i++ {
		temp[i] = contents[i-1]
	}
	return temp
}
