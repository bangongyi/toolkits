package office

import (
	"errors"
	"github.com/tealeg/xlsx"
	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
	unipdflicense "github.com/unidoc/unipdf/v3/common/license"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// word文件转文字
func WordToContent(key string, filePath string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(filePath)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}
	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}

	err = unipdflicense.SetMeteredKey(key)
	if err != nil {
		return "", "", 0, errors.New("初始化key失败1！")
	}
	// This example requires both for unioffice and unipdf.
	err = license.SetMeteredKey(key)
	if err != nil {
		return "", "", 0, errors.New("初始化key失败2！")
	}

	doc, err := document.Open(filePath)
	if err != nil {
		return "", "", 0, errors.New("打开文件失败！")
	}

	// 遍历文档中的每个段落
	text := ""
	for _, p := range doc.Paragraphs() {
		// 遍历段落中的每个 Run
		for _, r := range p.Runs() {
			// 打印 Run 的文本
			text += r.Text()
		}
	}

	return text, suffix, size, nil
}

// word地址文件转文字
func WordUrlToContent(key string, url string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(url)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}

	filePath, err := saveFile(url, suffix)
	if err != nil {
		return "", "", 0, errors.New("文件保存在本地失败！")
	}

	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}

	err = unipdflicense.SetMeteredKey(key)
	if err != nil {
		return "", "", 0, errors.New("初始化key失败1！")
	}
	// This example requires both for unioffice and unipdf.
	err = license.SetMeteredKey(key)
	if err != nil {
		return "", "", 0, errors.New("初始化key失败2！")
	}

	defer os.Remove(filePath)

	doc, err := document.Open(filePath)
	if err != nil {
		return "", "", 0, errors.New("打开文件失败！")
	}

	// 遍历文档中的每个段落
	text := ""
	for _, p := range doc.Paragraphs() {
		// 遍历段落中的每个 Run
		for _, r := range p.Runs() {
			// 打印 Run 的文本
			text += r.Text()
		}
	}
	return text, suffix, size, nil
}

// excel文件转文字
func ExcelToContent(filePath string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(filePath)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}
	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return "", "", 0, errors.New("打开文件失败！")
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
func ExcelUrlToContent(url string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(url)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}

	filePath, err := saveFile(url, suffix)
	if err != nil {
		return "", "", 0, errors.New("文件保存在本地失败！")
	}

	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}

	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return "", "", 0, errors.New("打开文件失败！")
	}

	defer os.Remove(filePath)

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
func TxtToContent(filePath string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(filePath)
	if err != nil {
		return "", "", 0, errors.New("获取前缀失败！")
	}
	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
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
func TxtUrlToContent(url string) (word string, fileSuffix string, FileSize int, err error) {
	suffix, err := getSuffix(url)
	if err != nil {
		return "", "", 0, err
	}

	filePath, err := saveFile(url, suffix)
	if err != nil {
		return "", "", 0, errors.New("文件保存在本地失败！")
	}

	size, err := countSize(filePath)
	if err != nil {
		return "", "", 0, errors.New("计算文件大小失败！")
	}

	defer os.Remove(filePath)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	text := string(content)
	text = strings.Replace(text, "\n", " ", -1)

	return text, suffix, size, nil
}
