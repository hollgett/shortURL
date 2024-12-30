package models

type RequestJSON struct {
	RequestURL string `json:"url"`
}

type ResponseJSON struct {
	ResponseURL string `json:"result"`
}
