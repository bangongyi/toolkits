package baidu

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	tokenUrlBaiDu     = "https://aip.baidubce.com/oauth/2.0/token"
	transformUrlBaidu = "https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic?access_token=%s"
)

type BodyResultResponse struct {
	LogId          int         `json:"log_id"`
	WordsResultNum int         `json:"words_result_num"`
	WordsResult    []WordsList `json:"words_result"`
}

type WordsList struct {
	Words string `json:"words"`
}

type BaiDuTokenResponse struct {
	RefreshToken     string `json:"refresh_token,omitempty"`
	ExpiresIn        int64  `json:"expires_in,omitempty"`
	SessionKey       string `json:"session_key,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	Scope            string `json:"scope,omitempty"`
	SessionSecret    string `json:"session_secret,omitempty"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type Cache interface {
	Set(key string, value string, expires int) error
	Get(key string) (string, error)
}

type BaiduOcr struct {
	cache     Cache
	apiKey    string
	apiSecret string
}

func NewBaiduOcr(apiKey string, apiSecret string, cache Cache) (*BaiduOcr, error) {
	c := &BaiduOcr{apiKey: apiKey, apiSecret: apiSecret, cache: cache}
	_, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// 图片转文字
func (b *BaiduOcr) ImageToWord(image string) (word string, err error) {

	//image := "/Users/zhangsan/go/test.jpeg"
	encode := b.getFileContentAsBase64(image)
	contextLen := len(encode)
	if contextLen/1024/1024 > 8 {
		return "", errors.New("文件大小不能大于8M")
	}
	payload := strings.NewReader("image=" + url.QueryEscape(encode) + "&detect_direction=false&detect_language=false&paragraph=false&probability=false")
	str, err := b.commonFun(payload)
	if err != nil {
		return "", err
	}

	return str, nil
}

// 图片地址转文字
func (b *BaiduOcr) ImageUrlToWord(imageUrl string) (word string, err error) {
	if len(imageUrl) > 1024 {
		return "", errors.New("图片地址不能超过 1024 个字节")
	}

	encode := b.getFileContentAsBase64(imageUrl)
	contextLen := len(encode)
	if contextLen/1024/1024 > 8 {
		return "", errors.New("文件大小不能大于8M")
	}
	payload := strings.NewReader("url=" + url.QueryEscape(encode) + "&detect_direction=false&detect_language=false&paragraph=false&probability=false")
	str, err := b.commonFun(payload)
	if err != nil {
		return "", err
	}

	return str, nil
}

// pdf转文字
func (b *BaiduOcr) PdfToWord(pdf string) (word string, err error) {

	//pdf := "/Users/zhangsan/go/test.pdf"
	encode := b.getFileContentAsBase64(pdf)
	contextLen := len(encode)
	if contextLen/1024/1024 > 8 {
		return "", errors.New("文件大小不能大于8M")
	}
	payload := strings.NewReader("pdf_file=" + url.QueryEscape(encode) + "&detect_direction=false&detect_language=false&paragraph=false&probability=false")

	str, err := b.commonFun(payload)
	if err != nil {
		return "", err
	}

	return str, nil
}

func (b *BaiduOcr) commonFun(payload *strings.Reader) (word string, err error) {
	token, err := b.getAccessToken()
	if err != nil {
		return "", err
	}

	requestUrl := fmt.Sprintf(transformUrlBaidu, token)

	client := &http.Client{}
	req, err := http.NewRequest("POST", requestUrl, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	resBody1 := BodyResultResponse{}
	err = json.Unmarshal(body, &resBody1)
	if err != nil {
		return
	}

	var str string
	for _, val := range resBody1.WordsResult {
		str += val.Words + ","
	}
	str = strings.TrimRight(str, ",")

	return str, nil
}

// base64编码后进行urlEncode
func (b *BaiduOcr) getFileContentAsBase64(path string) string {
	srcByte, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(srcByte)
}

// 获取token
func (b *BaiduOcr) getAccessToken() (token string, err error) {

	md5String, _ := md5ByString(b.apiKey)
	tokenKey := "kpai:baiduocr:" + md5String
	token, err = b.cache.Get(tokenKey)
	if err != nil {
		fmt.Printf("baidu gettoken redis token, err = %v \n", err)
	}
	if len(token) > 1 {
		fmt.Printf("baidu gettoken redis token is not empty, token 1%v1 \n", token)
		return token, nil
	}

	url := tokenUrlBaiDu + "?client_id=%s&client_secret=%s&grant_type=client_credentials"
	url = fmt.Sprintf(url, b.apiKey, b.apiSecret)
	payload := strings.NewReader(``)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Printf("baidu gettoken http.NewRequest, err %v\n", err)
		return token, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("baidu gettoken http.NewRequest Do, err %v\n", err)
		return token, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("baidu gettoken http.NewRequest ioutil.ReadAll, err %v\n", err)
		return token, err
	}

	var baiDuTokenResponse BaiDuTokenResponse
	err = json.Unmarshal(body, &baiDuTokenResponse)
	if err != nil {
		fmt.Printf("baidu gettoken http.NewRequest json.Unmarshal, body %v ;err %v\n", string(body), err)
		return token, err
	}
	if len(baiDuTokenResponse.Error) > 1 {
		fmt.Printf("baidu gettoken http.NewRequest err, ErrorMsg %v, ErrorCode %v \n", baiDuTokenResponse.Error, baiDuTokenResponse.ErrorDescription)
		return token, err
	}

	token = baiDuTokenResponse.AccessToken
	if len(token) > 0 {
		fmt.Printf("baidu translate save token, %v  \n", token)
		if err != nil {
			fmt.Printf("baidu gettoken save token, %v  \n", err.Error())
			return token, err
		}

		err = b.cache.Set(tokenKey, token, int(baiDuTokenResponse.ExpiresIn))
		if err != nil {
			fmt.Printf("baidu gettoken save token Expire,token = %v  err = %v \n", token, err)
			return token, err
		}
	}
	return token, nil
}

func md5ByString(str string) (string, error) {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		return "", err
	}
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr), nil
}
