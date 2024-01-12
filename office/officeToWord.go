package office

import (
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"log"
	"strings"
)

// excel文件转文字
func ExcelToWord(filePath string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(filePath)
	if err != nil {
		return "", "", 0, err
	}
	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, err
	}
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return "", "", 0, err
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
	return text, suffix, size, nil
}

// excel地址文件转文字
func ExcelUrlToWord(url string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(url)
	if err != nil {
		return "", "", 0, err
	}

	filePath, err := saveFile(url, suffix)
	if err != nil {
		return
	}

	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, err
	}

	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return "", "", 0, err
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

	return text, suffix, size, nil
}

// txt文件转文字
func TxtToWord(filePath string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(filePath)
	if err != nil {
		return "", "", 0, err
	}
	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, err
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	text := string(content)
	text = strings.Replace(text, "\n", " ", -1)

	return text, suffix, size, nil
}

// txt地址文件转文字
func TxtUrlToWord(url string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(url)
	if err != nil {
		return "", "", 0, err
	}

	filePath, err := saveFile(url, suffix)
	if err != nil {
		return
	}

	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, err
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	text := string(content)
	text = strings.Replace(text, "\n", " ", -1)

	return text, suffix, size, nil
}
