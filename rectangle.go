package reportengine

import (
	"github.com/signintech/gopdf"
	"log"
)

type Rectangle struct {
	lowerX          float64
	lowerY          float64
	width           float64
	height          float64
	border          int
	backgroundColor Color
	borderColor     Color
	borderType      int
	borderLineWidth float64
	isVisible       bool
}

func NewRectangle(border, borderType int, borderLineWidth float64, backgroundColor, borderColor Color, isVisible bool) Rectangle {
	return Rectangle{border: border, borderType: borderType, borderLineWidth: borderLineWidth,
		backgroundColor: backgroundColor, borderColor: borderColor, isVisible: isVisible}
}

func (t Rectangle) Render(pdf *gopdf.GoPdf) {
	if !t.isVisible {
		return
	}
	var err error
	pdf.SetStrokeColor(t.borderColor.r, t.borderColor.g, t.borderColor.b)
	pdf.SetFillColor(t.backgroundColor.r, t.backgroundColor.g, t.backgroundColor.b)
	pdf.SetLineWidth(t.borderLineWidth)
	switch t.borderType {
	case Dashed:
		pdf.SetLineType("dashed")
	case Dotted:
		pdf.SetLineType("dotted")
	case Solid:
		pdf.SetLineType("")
	}
	margin := 0.0 //t.borderLineWidth / 2.0
	pdf.RectFromUpperLeftWithStyle(t.lowerX+margin, t.lowerY+margin, t.width-margin*2, t.height-margin*2, "F")
	pdf.SetX(t.lowerX + margin)
	pdf.SetY(t.lowerY + margin)
	temp := gopdf.Rect{W: t.width - margin*2, H: t.height - margin*2}
	err = pdf.SetFont(IconFontFamily, "", 14)
	if err != nil {
		if err.Error() == "not found font family" {
			LoadFont(pdf, IconFontFamily)
		}
		err = pdf.SetFont(IconFontFamily, "", 14)
		if err != nil {
			log.Println("4")
			log.Println(err.Error())
			return
		}
	}
	err = pdf.CellWithOption(&temp, "", gopdf.CellOption{Border: t.border})
	if err != nil {
		log.Println("5")
		log.Println(err.Error())
		return
	}
}
