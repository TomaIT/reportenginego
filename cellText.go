package reportengine

import (
	"github.com/signintech/gopdf"
	"regexp"
	"strconv"
	"strings"
)

type CellText struct {
	horizontalAlign uint
	verticalAlign   uint
	rectangle       Rectangle
	underline       bool
	fontSize        int
	color           Color
	fontFamily      string
	originalValue   string
	minMarginText   Margin
	tokens          []Token
}

func NewCellText(horizontalAlign uint, verticalAlign uint, value string, underline bool,
	fontFamily string, fontSize int, color Color, minMarginText Margin, rectangle Rectangle) *CellText {
	ct := new(CellText)
	ct.horizontalAlign = horizontalAlign
	ct.verticalAlign = verticalAlign
	ct.underline = underline
	ct.fontSize = fontSize
	ct.color = color
	ct.originalValue = value
	ct.minMarginText = minMarginText
	ct.fontFamily = fontFamily
	ct.setTokens(value, fontFamily, fontSize, color)
	ct.rectangle = rectangle
	return ct
}

func (t *CellText) Build(pdf *gopdf.GoPdf, maxWidth float64) {
	//To reset previous Shorten called
	t.toOriginal()
	for range t.tokens {
		minWidth := t.MinWidth(pdf)
		if minWidth <= maxWidth {
			break
		}
		index := len(t.tokens) - 1
		tokenWidth := t.tokens[index].Width(pdf)
		if tokenWidth > minWidth-maxWidth && t.tokens[index].fontFamily != IconFontFamily {
			if t.tokens[index].Shorten(pdf, tokenWidth-minWidth+maxWidth) {
				break
			}
		}
		t.tokens = t.tokens[:len(t.tokens)-1]
	}
	t.rectangle.width = maxWidth
	t.rectangle.height = t.MinHeight()
	t.rectangle.lowerX = 0
	t.rectangle.lowerY = 0
}
func (t *CellText) Adjust(pdf *gopdf.GoPdf, lowerX, lowerY, width, height float64) {
	if t.MinWidth(pdf) > width || t.MinHeight() > height {
		panic("Width/Height are not sufficient")
	}
	t.rectangle.lowerX = lowerX
	t.rectangle.lowerY = lowerY
	t.rectangle.width = width
	t.rectangle.height = height
}
func (t *CellText) MoveTo(lowerX, lowerY float64) {
	t.rectangle.lowerX = lowerX
	t.rectangle.lowerY = lowerY
}
func (t *CellText) SetVisibilityContainer(isVisible bool) {
	t.rectangle.isVisible = isVisible
}
func (t *CellText) Split(*gopdf.GoPdf, float64, int) Component {
	return nil
}
func (t CellText) MinWidth(pdf *gopdf.GoPdf) float64 {
	return t.textWidth(pdf) + t.minMarginText.left + t.minMarginText.right
}
func (t CellText) MinHeight() float64 {
	temp := t.textHeight() + t.minMarginText.top + t.minMarginText.bottom
	if t.underline {
		return temp + (float64(t.fontSize) * UnderlineWidthFactor) + UnderlineMargin
	}
	return temp
}
func (t CellText) Render(pdf *gopdf.GoPdf) {
	t.rectangle.Render(pdf)
	t.renderTokens(pdf)
}
func (t CellText) FirstVoidSpace() Rectangle {
	panic("Not implemented")
}
func (t CellText) GetRectWidth() float64 {
	return t.rectangle.width
}
func (t CellText) GetRectHeight() float64 {
	return t.rectangle.height
}
func (t CellText) GetRectPosition() (x, y float64) {
	return t.rectangle.lowerX, t.rectangle.lowerY
}
func (t CellText) IsSplittable() bool {
	return false
}

func (t *CellText) toOriginal() {
	t.setTokens(t.originalValue, t.fontFamily, t.fontSize, t.color)
}
func (t *CellText) setTokens(value string, fontFamily string, fontSize int, color Color) {
	regexIcon := regexp.MustCompile(`i{(0x[a-zA-Z0-9]{4,5};#[a-zA-Z0-9]{6})}`)
	icons := regexIcon.FindAllStringSubmatch(value, -1)
	temp := value
	t.tokens = make([]Token, 0)
	for _, val := range icons {
		if temp != "" {
			x := strings.Split(temp, val[0])
			temp = x[1]
			if x[0] != "" {
				t.tokens = append(t.tokens, Token{fontFamily: fontFamily, fontSize: fontSize, value: x[0], color: color})
			}
		}
		x, y := getIcon(val[1])
		t.tokens = append(t.tokens, Token{fontFamily: IconFontFamily, fontSize: fontSize, value: x, color: y})
	}
	if temp != "" {
		t.tokens = append(t.tokens, Token{fontFamily: fontFamily, fontSize: fontSize, value: temp, color: color})
	}

}
func (t CellText) textWidth(pdf *gopdf.GoPdf) float64 {
	tot := 0.0
	for i := range t.tokens {
		tot += t.tokens[i].Width(pdf)
	}
	return tot
}
func (t CellText) textHeight() float64 {
	max := 0.0
	for i := range t.tokens {
		temp := t.tokens[i].Height()
		if temp > max {
			max = temp
		}
	}
	return max
}
func (t CellText) getTextStartPosition(pdf *gopdf.GoPdf) (x float64, y float64) {
	textWidth := t.textWidth(pdf)
	textHeight := t.textHeight()
	switch t.horizontalAlign {
	case gopdf.Left:
		x = t.rectangle.lowerX + t.minMarginText.left
	case gopdf.Right:
		x = t.rectangle.lowerX + t.rectangle.width - t.minMarginText.right - textWidth
	case gopdf.Center:
		x = t.rectangle.lowerX + (t.rectangle.width-textWidth)/2.0
	}
	underlineMargin := 0.0
	if t.underline {
		underlineMargin = float64(t.fontSize)*UnderlineWidthFactor + UnderlineMargin
	}
	switch t.verticalAlign {
	case gopdf.Top:
		y = t.rectangle.lowerY + textHeight + t.minMarginText.top
	case gopdf.Middle:
		y = t.rectangle.lowerY + t.rectangle.height - underlineMargin - (t.rectangle.height-textHeight-underlineMargin)/2.0
	case gopdf.Bottom:
		y = t.rectangle.lowerY + t.rectangle.height - t.minMarginText.bottom - underlineMargin
	}
	return x, y
}
func (t CellText) renderTokens(pdf *gopdf.GoPdf) {
	lowerX, upperY := t.getTextStartPosition(pdf)
	startX := lowerX
	for i := range t.tokens {
		t.tokens[i].Render(pdf, lowerX, upperY)
		lowerX += t.tokens[i].Width(pdf)
	}
	if t.underline {
		pdf.SetLineWidth(float64(t.fontSize) * UnderlineWidthFactor)
		pdf.SetStrokeColor(t.color.r, t.color.g, t.color.b)
		pdf.SetLineType("")
		pdf.Line(startX, upperY+UnderlineMargin+(float64(t.fontSize)*UnderlineWidthFactor), lowerX, upperY+UnderlineMargin+(float64(t.fontSize)*UnderlineWidthFactor))
	}
}
func merge(a CellText, b CellText, delimiter string) (res CellText) {
	res.horizontalAlign = a.horizontalAlign
	res.verticalAlign = a.verticalAlign
	res.rectangle = a.rectangle
	res.underline = a.underline
	res.fontSize = a.fontSize
	res.color = a.color
	res.fontFamily = a.fontFamily
	res.originalValue = a.originalValue + delimiter + b.originalValue
	res.minMarginText = a.minMarginText
	res.tokens = make([]Token, 0)
	res.tokens = append(res.tokens, a.tokens...)
	res.tokens = append(res.tokens, Token{fontFamily: res.fontFamily, fontSize: res.fontSize, value: delimiter, color: res.color})
	res.tokens = append(res.tokens, b.tokens...)
	return res
}
func getIcon(value string) (icon string, color Color) {
	//value = i{0xF0A43;#0000FF}
	fields := strings.Split(value, ";")
	icon = fields[0]
	dec := hexToInt32(icon)
	icon, _ = strconv.Unquote(strconv.QuoteRuneToGraphic(dec))
	color, _ = NewColor(fields[1])
	return icon, color
}
func hexToInt32(hexString string) int32 {
	numberStr := strings.Replace(hexString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	val, _ := strconv.ParseInt(numberStr, 16, 32)
	return int32(val)
}
