package office

import (
	"github.com/tealeg/xlsx"
)

func ExcelToWord(location string) (word string, err error) {
	xlFile, err := xlsx.OpenFile(location)
	if err != nil {
		return "", err
	}

	text := ""
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				text += cell.String() + " "
			}
			text += ""
		}
	}
	return text, nil
}
