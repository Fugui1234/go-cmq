package cmq

type MsgInfo struct {
	MsgBody          string `json:"msgBody"`
	MsgId            string `json:"msgId"`
	ReceiptHandle    string `json:"receiptHandle"`
	EnqueueTime      int64  `json:"enqueueTime"`
	FirstDequeueTime int64  `json:"firstDequeueTime"`
	NextVisibleTime  int64  `json:"nextVisibleTime"`
	DequeueCount     int64  `json:"dequeueCount"`
}

type BatchReceiveRes struct {
	Code        int       `json:"code"`
	Message     string    `json:"message"`
	RequestId   string    `json:"requestId"`
	MsgInfoList []MsgInfo `json:"msgInfoList"`
}
