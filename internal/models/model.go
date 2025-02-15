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

type RequestBatch struct {
	CorrId   string `json:"correlation_id"`
	Original string `json:"original_url"`
}

type ResponseBatch struct {
	CorrId string `json:"correlation_id"`
	Short  string `json:"short_url"`
}

type DBBatch struct {
	Short    string `db:"short"`
	CorrId   string
	Original string `db:"origin"`
}
