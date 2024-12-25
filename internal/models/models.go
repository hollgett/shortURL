package models

type RequestJson struct {
	RequestURL string `json:"url"`
}

type ResponseJson struct {
	ResponseURL string `json:"result"`
}
