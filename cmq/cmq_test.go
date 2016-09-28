package cmq

import (
	"fmt"
	"git.gumpcome.com/leonardyp/go-gckit/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func init() {
	logger.Init()
}
func TestCreateQueue(t *testing.T) {
	c := Init("AKIDBSZfcObWmpsreKgDOVTyJdu439JQhkfP", "hNQLCE65uXxDH9qKqJeppiSG5iT3QbaN", "bj", true)

	c.GenSignString("POST", "CreateQueue", map[string]string{
		"queueName":          "test-queue-1",
		"pollingWaitSeconds": "30",
	})

	params := url.Values{}
	params.Add("Action", c.Action)
	params.Add("Region", c.Region)
	params.Add("Timestamp", fmt.Sprint(c.TimeStamp))
	params.Add("Nonce", c.Nonce)
	params.Add("SecretId", c.SecretId)
	params.Add("Signature", c.Sign)
	params.Add("queueName", "yp-test-queue-1")
	params.Add("pollingWaitSeconds", "30")

	logger.DebugStd("%v", params)
	resp, err := http.PostForm("https://cmq-queue-bj.api.qcloud.com/v2/index.php", params)
	if err != nil {
		logger.ErrorStd("error:%v" + err.Error())
		return
	}

	bys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorStd("%v", err)
		return
	}
	logger.DebugStd("%v", string(bys))
}
