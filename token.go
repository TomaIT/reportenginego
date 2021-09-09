package reportengine

import (
	"github.com/signintech/gopdf"
	"log"
	"strconv"
)

type Token struct {
	value      string
	fontFamily string
	fontSize   int
	color      Color
}

func (t Token) Render(pdf *gopdf.GoPdf, lowerX float64, upperY float64) {
	var err error
	pdf.SetTextColor(t.color.r, t.color.g, t.color.b)
	err = pdf.SetFont(t.fontFamily, "", t.fontSize)
	if err != nil {
		if err.Error() == "not found font family" {
			LoadFont(pdf, t.fontFamily)
			err = pdf.SetFont(t.fontFamily, "", t.fontSize)
			if err != nil {
				log.Println(err.Error())
				return
			}
		} else {
			log.Println(err.Error())
			return
		}
	}
	pdf.SetX(lowerX)
	pdf.SetY(upperY)
	err = pdf.Text(t.value)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (t Token) Width(pdf *gopdf.GoPdf) float64 {
	return Width(pdf, t.fontFamily, t.fontSize, t.value)
}

func (t Token) Height() float64 {
	return gopdf.ContentObjCalTextHeight(t.fontSize)
}

/**
Return Value
	-true: success
	-false: not possible shorten (remove token)
*/
func (t *Token) Shorten(pdf *gopdf.GoPdf, maxWidth float64) bool {
	width := Width(pdf, t.fontFamily, t.fontSize, t.value)
	if width <= maxWidth {
		return true
	}
	if t.fontFamily == IconFontFamily {
		return false
	}
	temp := ""
	width = Width(pdf, t.fontFamily, t.fontSize, temp+ShortenCharacters)
	if width > maxWidth {
		return false
	}

	for _, r := range t.value {
		t.value = temp + ShortenCharacters
		s, _ := strconv.Unquote(strconv.QuoteRune(r))
		temp += s
		width = Width(pdf, t.fontFamily, t.fontSize, temp+ShortenCharacters)
		if width > maxWidth {
			return true
		}
	}
	return false
}
