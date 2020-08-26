package ignition

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	ignitionLogger = ctrl.Log.WithName("ignition")
)

func intToPtr(i int) *int {
	return &i
}

func boolToPtr(b bool) *bool {
	return &b
}

func StringToPtr(s string) *string {
	return &s
}

func GetS3Url(bucket, uri string) string {
	url := &url.URL{
		Scheme: "s3",
		Host:   bucket,
		Path:   uri,
	}
	return url.String()
}

func GetHttpUrl(address, uri string) *string{
	url := &url.URL{
		Scheme: "https",
		Host: address,
		Path: uri,
	}
	res := url.String()
	return &res
}

func UploadFile(address,uploadPath,filename string) error {
	method := "POST"

	playload := &bytes.Buffer{}
	writer := multipart.NewWriter(playload)
	file, err := os.Open(filename)
	if err != nil{
		ignitionLogger.Error(err,"file does not exit")
	}
	defer file.Close()
	part1, err := writer.CreateFormFile("file",filepath.Base(filename))
	_, err = io.Copy(part1,file)
	if err != nil{
		ignitionLogger.Error(err,"copy file error")
		return err
	}
	_ = writer.WriteField("path",uploadPath)
	err = writer.Close()
	if err != nil{
		ignitionLogger.Error(err,"can not close write")
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, "https://"+address+"/upload", playload)
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil{
		ignitionLogger.Error(err,"can not get sercer")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil{
		ignitionLogger.Error(err,"upload file failed")
		return err
	}else{
		ignitionLogger.Info(string(body))
	}
	return nil
}