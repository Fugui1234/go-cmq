package cmq

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/golang/glog"
)

type Cmq struct {
	Action     string
	SecretId   string
	SecretKey  string
	Region     string
	RandSize   int
	Inited     bool
	IsInner    bool
	InnerAddr  string
	OutterAddr string
	Uri        string
	TimeStamp  int64
	Nonce      string
	Sign       string
}

const CLOUD_API_URI = "/v2/index.php"

func Init(secretid, secretkey, region string, isInner bool) *Cmq {
	if secretid == "" || secretkey == "" || region == "" {
		panic("auth appid,secretid or secretkey or region 不能为空")
	}
	return &Cmq{
		SecretId:   secretid,
		SecretKey:  secretkey,
		Region:     region,
		Inited:     true,
		IsInner:    isInner,
		InnerAddr:  "cmq-queue-" + region + ".api.tencentyun.com",
		OutterAddr: "cmq-queue-" + region + ".api.qcloud.com",
	}
}

func (c *Cmq) GenSignString(method, action string, params map[string]string) {
	now := time.Now()
	timestamp := now.Unix()
	if c.RandSize == 0 {
		c.RandSize = 10000
	}

	nonce := strconv.Itoa(rand.New(rand.NewSource(now.UnixNano())).Intn(c.RandSize))
	c.Nonce = nonce
	c.TimeStamp = timestamp
	c.Action = action

	params["Action"] = action
	params["Nonce"] = strconv.Itoa(rand.New(rand.NewSource(now.UnixNano())).Intn(c.RandSize))
	params["Region"] = c.Region
	params["SecretId"] = c.SecretId
	params["Timestamp"] = fmt.Sprint(timestamp)

	var keys = make([]string, 0, len(params))
	for key, _ := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var orignal string
	for _, key := range keys {
		orignal += key + "=" + params[key] + "&"
	}
	if c.IsInner {
		orignal = method + c.InnerAddr + CLOUD_API_URI + "?" + orignal[:len(orignal)-1]
	} else {
		orignal = method + c.OutterAddr + CLOUD_API_URI + "?" + orignal[:len(orignal)-1]
	}
	c.Sign = base64.StdEncoding.EncodeToString([]byte(Hmac_Sha1(orignal, c.SecretKey)))
}

func (c *Cmq) commonParams() url.Values {
	params := url.Values{}
	params.Add("Action", c.Action)
	params.Add("Region", c.Region)
	params.Add("Timestamp", fmt.Sprint(c.TimeStamp))
	params.Add("Nonce", c.Nonce)
	params.Add("SecretId", c.SecretId)
	params.Add("Signature", c.Sign)
	return params
}

func (c *Cmq) getUrl() string {
	addr := ""
	if c.IsInner {
		addr = "http://" + c.InnerAddr
	} else {
		addr = "https://" + c.OutterAddr
	}
	return addr + CLOUD_API_URI
}

func (c *Cmq) sendRequest(action string, specificParams map[string]string) ([]byte, error) {
	c.GenSignString("POST", action, specificParams)
	params := c.commonParams()
	for key, value := range specificParams {
		params.Add(key, value)
	}
	resp, err := http.PostForm(c.getUrl(), params)
	if err != nil {
		glog.Infoln(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	bys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Infoln("%v", err)
		return nil, err
	}
	return bys, nil
}

func (c *Cmq) CreateQueue(queueName string) {
	if _, err := c.sendRequest("CreateQueue", map[string]string{
		"queueName": queueName,
	}); err != nil {
		glog.Infoln(err)
	}
}

func (c *Cmq) DeleteQueue(queueName string) {
	if _, err := c.sendRequest("DeleteQueue", map[string]string{
		"queueName": queueName,
	}); err != nil {
		glog.Infoln(err)
	}
}

func (c *Cmq) SendMessage(queueName string, msgBody string) error {
	var i int
	for i = 3; i > 0; i-- {
		bys, err := c.sendRequest("SendMessage", map[string]string{
			"queueName": queueName,
			"msgBody":   msgBody,
		})
		glog.Infoln(string(bys))
		if result, jsonErr := simplejson.NewJson(bys); jsonErr == nil {
			if code, jsonErr := result.Get("code").Int(); jsonErr == nil {
				if code == 0 {
					return err
				} else {
					glog.Infoln("Error to send message!!!")
				}
			}
		}
	}
	return errors.New("send error")
}

func (c *Cmq) BatchSendMessage(queueName string, msgBodys []string) {
	params := map[string]string{}
	for index, msgBody := range msgBodys {
		params["msgBody."+strconv.Itoa(index+1)] = msgBody
	}
	params["queueName"] = queueName
	if _, err := c.sendRequest("BatchSendMessage", params); err != nil {
		glog.Infoln(err)
	}
}

func (c *Cmq) BatchReceiveMessage(queueName string, msgNum int) ([]byte, error) {
	return c.sendRequest("BatchReceiveMessage", map[string]string{
		"queueName": queueName,
		"numOfMsg":  strconv.Itoa(msgNum),
	})
}

func (c *Cmq) BatchDeleteMessage(queueName string, receiptHandles []string) {
	params := map[string]string{}
	for index, receiptHandle := range receiptHandles {
		params["receiptHandle."+strconv.Itoa(index+1)] = receiptHandle
	}
	params["queueName"] = queueName
	if _, err := c.sendRequest("BatchDeleteMessage", params); err != nil {
		glog.Infoln(err)
	}
}

func Hmac_Sha1(data, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return string(mac.Sum(nil))
}
