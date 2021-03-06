package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"

	"proxy/common"
)

var config *viper.Viper

func otcPrice(w http.ResponseWriter, req *http.Request) {
	url := "https://otc-api.huobi.co/v1/data/market/detail"

	httpClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
        log.Printf("otcprice req create failed, msg error: %s", err)
        w.WriteHeader(http.StatusBadRequest)
        return
	}

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Printf("otcprice request failed, msg error: %s", getErr)
        w.WriteHeader(http.StatusBadRequest)
        return
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Printf("otcprice data read failed, msg error: %s", readErr)
        w.WriteHeader(http.StatusBadRequest)
        return
	}
	fmt.Fprintf(w, string(body))
}

func ossUpload(w http.ResponseWriter, req *http.Request) {
	endpoint := config.GetString("oss.endpoint")
	accesskeyId := config.GetString("oss.accesskey.id")
	accesskeySecret := config.GetString("oss.accesskey.secret")
	bucketName := config.GetString("oss.bucket.name")
	objectNameBaseDir := config.GetString("oss.objectname.basedir")

    err := req.ParseMultipartForm(32 << 20)
    if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
    
	file, handler, formFileErr := req.FormFile("files")
	if formFileErr != nil {
        log.Printf("parse request failed, msg error: %s", formFileErr)
        w.WriteHeader(http.StatusBadRequest)
        return
	}
    if file != nil {
        defer file.Close()
    }

	client, ossClientCreateErr := oss.New(endpoint, accesskeyId, accesskeySecret)
	if ossClientCreateErr != nil {
        log.Printf("oss client create failed, msg error: %s", ossClientCreateErr)
        w.WriteHeader(http.StatusBadRequest)
        return
	}

	bucket, bucketCreateErr := client.Bucket(bucketName)
	if bucketCreateErr != nil {
        log.Printf("oss bucket create failed, msg error: %s", bucketCreateErr)
        w.WriteHeader(http.StatusBadRequest)
        return
	}

	objectName := objectNameBaseDir + handler.Filename
	putObjectErr := bucket.PutObject(objectName, file)
	if putObjectErr != nil {
        log.Printf("oss upload file failed, msg error: %s", putObjectErr)
        w.WriteHeader(http.StatusBadRequest)
        return
	}
    log.Printf("upload file %s to %s", handler.Filename, objectName)
    return
}

func main() {
	common.Init()
	config = common.GetConfig()

	http.HandleFunc("/otcprice", otcPrice)
	http.HandleFunc("/aliyun/oss/upload", ossUpload)

	http.ListenAndServe(":8080", nil)
}
