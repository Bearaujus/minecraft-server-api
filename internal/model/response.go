package model

type ResponseHeader struct {
	ProcessTime interface{} `json:"process_time"`
	IsSuccess   bool        `json:"is_success"`
	Messages    interface{} `json:"messages"`
}

type Response struct {
	Header ResponseHeader `json:"header"`
	Data   interface{}    `json:"data"`
}
