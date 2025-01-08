package jsonutil

import (
	"encoding/json"
	"io"

	"github.com/hollgett/shortURL.git/internal/models"
)

func EncodeJSON(w io.Writer, respData models.ResponseJSON) error {
	encode := json.NewEncoder(w)
	return encode.Encode(respData)
}

func DecodeJSON(r io.Reader, reqData *models.RequestJSON) error {
	decode := json.NewDecoder(r)
	return decode.Decode(reqData)
}
