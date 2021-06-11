package stryd

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HttpUpload(urlStr string, content string) ([]byte, error) {

	var extName string

	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(data)

	if !strings.Contains(contentType, "image") {
		return nil, errors.New("File mime type invalid, only accept image.")
	}

	if strings.Contains(contentType, "jpeg") {
		extName = "jpg"
	} else if strings.Contains(contentType, "png") {
		extName = "png"
	} else {
		return nil, errors.New("File mime type invalid, only accept jpg or png.")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileName := fmt.Sprintf("%s.%s", Md5(fmt.Sprintf("%d%s", time.Now().UnixNano(), urlStr)), extName)

	fileData, _ := writer.CreateFormFile("media", fileName)

	_, err = fileData.Write(data)
	if err != nil {
		return nil, err
	}

	err = writer.Close()

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept-Encoding", "gzip, deflate")
	request.Header.Set("Content-Type", writer.FormDataContentType())

	tr := &http.Transport{
		//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	httpClient.Timeout = 60 * time.Second

	resp, err := httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") == "gzip" {
		resp.Body, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("wechat api status is :non-200")
	}

	bodys, err := ioutil.ReadAll(resp.Body)
	return bodys, err

}

func HttpPost(api string, headers map[string]string, param map[string]interface{}) ([]byte, error) {

	buf := new(bytes.Buffer)

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	err := enc.Encode(param)

	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("POST", api, strings.NewReader(buf.String()))

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func HttpGet(api string, headers map[string]string, param map[string]interface{}) ([]byte, error) {

	queryStr, err := build(param)

	if err != nil {
		return nil, err
	}

	apiInfo, err := url.Parse(api)

	if err != nil {
		return nil, err
	}

	if apiInfo.RawQuery == "" {
		api = fmt.Sprintf("%s?%s", api, queryStr)
	} else {
		api = fmt.Sprintf("%s&%s", api, queryStr)
	}

	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, _ := http.NewRequest("GET", api, nil)

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := (&http.Client{}).Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func build(raw map[string]interface{}) (string, error) {

	p := make(map[string]string)

	for k, v := range raw {

		switch vv := v.(type) {
		case []interface{}:

			parseNormal(p, vv, []string{k})
		case map[string]interface{}:

			parseKeyValue(p, vv, []string{k})
		default:

			p[k] = fmt.Sprintf("%s", vv)
		}
	}

	data := url.Values{}

	for k, v := range p {
		data.Add(k, v)
	}

	return data.Encode(), nil
}

func parseKeyValue(p map[string]string, raw map[string]interface{}, keys []string) {

	for k, v := range raw {
		switch vv := v.(type) {
		case []interface{}:

			tmpKeys := append(keys, k)

			parseNormal(p, vv, tmpKeys)

		case map[string]interface{}:

			tmpKeys := append(keys, k)

			parseKeyValue(p, vv, tmpKeys)

		default:

			//keys = append(keys, k)

			var tmp []string

			for m, n := range keys {
				if m > 0 {
					n = fmt.Sprintf("[%s]", n)
				}

				tmp = append(tmp, n)
			}

			kStr := strings.Join(tmp, "")

			p[fmt.Sprintf("%s[%s]", kStr, k)] = fmt.Sprintf("%s", vv)
		}
	}
}

func parseNormal(p map[string]string, raw []interface{}, keys []string) {

	for k, v := range raw {
		switch vv := v.(type) {
		case []interface{}:

			tmpKeys := append(keys, fmt.Sprintf("%d", k))

			parseNormal(p, vv, tmpKeys)

		case map[string]interface{}:

			tmpKeys := append(keys, fmt.Sprintf("%d", k))

			parseKeyValue(p, vv, tmpKeys)

		default:

			//keys = append(keys, fmt.Sprintf("%d", k))

			var tmp []string

			for m, n := range keys {
				if m > 0 {
					n = fmt.Sprintf("[%s]", n)
				}

				tmp = append(tmp, n)
			}

			kStr := strings.Join(tmp, "")

			p[fmt.Sprintf("%s[%d]", kStr, k)] = fmt.Sprintf("%s", vv)
		}
	}
}
