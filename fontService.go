package reportengine

import (
	"github.com/signintech/gopdf"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var fontMap map[string]string

func init() {
	fontMap = make(map[string]string)
	libRegEx, err := regexp.Compile("\\.ttf$")
	if err != nil {
		panic(err.Error())
	}
	err = filepath.Walk(RootDirectoryFonts, func(path string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {
			family := strings.Replace(info.Name(), ".ttf", "", -1)
			fontMap[family] = path
		}
		return nil
	})
	if err != nil {
		panic(err.Error())
	}
}
func LoadFont(pdf *gopdf.GoPdf, family string) {
	err := pdf.AddTTFFont(family, fontMap[family])
	if err != nil {
		panic(err.Error())
	}
}
func Width(pdf *gopdf.GoPdf, fontFamily string, fontSize int, text string) float64 {
	err := pdf.SetFont(fontFamily, "", fontSize)
	if err != nil {
		if err.Error() == "not found font family" {
			LoadFont(pdf, fontFamily)
			err = pdf.SetFont(fontFamily, "", fontSize)
			if err != nil {
				panic(err)
			}
		}
	}
	x, err := pdf.MeasureTextWidth(text)
	if err != nil {
		panic(err)
	}
	return x
}
