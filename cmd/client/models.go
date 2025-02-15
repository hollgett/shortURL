package main

type RequestJSON struct {
	RequestURL string `json:"url"`
}

type ResponseJSON struct {
	ResponseURL string `json:"result"`
}

type RequestBatch struct {
	Id       string `json:"correlation_id"`
	Original string `json:"original_url"`
}

type ResponseBatch struct {
	Id    string `json:"correlation_id"`
	Short string `json:"original_url"`
}
