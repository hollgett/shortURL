package models

type RequestJSON struct {
	RequestURL string `json:"url"`
}

type ResponseJSON struct {
	ResponseURL string `json:"result"`
}

type FileStorageData struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}
