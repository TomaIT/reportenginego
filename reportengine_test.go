package reportengine

import (
	"github.com/signintech/gopdf"
	"log"
	"testing"
)

const testOutputDirectory = "testOutput/"

func TestCellTextArea(t *testing.T) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	ct := getCellTextAreaStr("Questo è un  cane blu i{0xF0A43;#0000FF} e questo è un gatto i{0xF011B;#FF0000} rosso. " +
		"FINEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE" +
		"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
	var x Component
	x = ct
	x.Build(pdf, 60)
	x.Adjust(pdf, 20, 20, 60, x.MinHeight())
	x.Render(pdf)
	x.Build(pdf, 80)
	x.Adjust(pdf, 80, 20, 80, x.MinHeight())
	x.Render(pdf)
	x.Build(pdf, 120)
	x.Adjust(pdf, 160, 20, 120, x.MinHeight())
	x.Render(pdf)
	err := pdf.WritePdf(testOutputDirectory + "TestCellTextArea.pdf")
	if err != nil {
		panic(err)
	}
}
func TestCellImage(t *testing.T) {
	var c Component
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	c = getGridWithImages(75)
	c.Build(pdf, gopdf.PageSizeA4.W-10)
	c.Adjust(pdf, 5, 0, c.GetRectWidth(), c.GetRectHeight())
	c.Render(pdf)
	err := pdf.WritePdf(testOutputDirectory + "TestCellImage.pdf")
	if err != nil {
		panic(err)
	}
}
func TestTable(t *testing.T) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	table := getTable(5, 8)
	table.Build(pdf, gopdf.PageSizeA4.W-20)
	table.Adjust(pdf, 10, 10, table.MinWidth(pdf), table.MinHeight())
	table.Render(pdf)
	err := pdf.WritePdf(testOutputDirectory + "TestTable.pdf")
	if err != nil {
		panic(err)
	}
}
func TestGrid(t *testing.T) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	grid := getGrid()
	grid.Build(pdf, 350)
	grid.Adjust(pdf, 50, 100, grid.MinWidth(pdf), grid.MinHeight())
	grid.MoveTo(5, 5)
	grid.Render(pdf)
	err := pdf.WritePdf(testOutputDirectory + "TestGrid.pdf")
	if err != nil {
		panic(err)
	}
}
func TestReport_1(t *testing.T) {
	report := NewReport(*gopdf.PageSizeA4, 10, 10, 10, 10,
		20)
	report.SetHeaderFP(getGrid())
	report.AddContentFP(getTable(26, 5))
	report.AddContentFP(getTable(30, 5))
	report.SetFooterFP(getGrid())
	report.SetHeaderLP(getGrid())
	report.AddContentLP(getTable(5, 7))
	report.SetFooterLP(getGrid())
	report.SetHeaderCP(getGrid())
	report.AddContentCP(getTable(189, 5))
	report.SetFooterCP(getGrid())
	report.Build()
	report.Render()
	err := report.pdf.WritePdf(testOutputDirectory + "TestReport_1.pdf")
	if err != nil {
		log.Print(err.Error())
		return
	}
}
func TestReport_2(t *testing.T) {
	report := NewReport(*gopdf.PageSizeA4, 10, 10, 10, 10,
		20)
	report.SetHeaderCP(getGrid())
	report.AddContentCP(getTable(189, 5))
	report.SetFooterCP(getGrid())
	report.Build()
	report.Render()
	err := report.pdf.WritePdf(testOutputDirectory + "TestReport_2.pdf")
	if err != nil {
		log.Print(err.Error())
		return
	}
}
func TestReport_3(t *testing.T) {
	report := NewReport(*gopdf.PageSizeA4, 10, 10, 10, 10,
		20)
	report.SetHeaderFP(getGrid())
	report.AddContentCP(getTable(189, 5))
	report.SetFooterLP(getGrid())
	report.Build()
	report.Render()
	err := report.pdf.WritePdf(testOutputDirectory + "TestReport_3.pdf")
	if err != nil {
		log.Print(err.Error())
		return
	}
}
func TestReport_4(t *testing.T) {
	report := NewReport(*gopdf.PageSizeA5, 10, 10, 10, 10,
		20)
	report.SetHeaderCP(getGrid())
	report.AddContentCP(getTable(35, 5))
	report.AddContentCP(getGridWithImages(300))
	report.SetFooterCP(getGrid())
	report.Build()
	report.Render()
	err := report.pdf.WritePdf(testOutputDirectory + "TestReport_4.pdf")
	if err != nil {
		log.Print(err.Error())
		return
	}
}
