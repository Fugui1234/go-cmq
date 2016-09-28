package cmq

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"time"
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
	l, _ := url.Parse("http://www.baidu.com?" + base64.StdEncoding.EncodeToString([]byte(Hmac_Sha1(orignal, c.SecretKey))))
	c.Sign = l.Query().Encode()
}

// @Description
func Hmac_Sha1(data, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return string(mac.Sum(nil))
}
