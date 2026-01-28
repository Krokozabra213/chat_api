package request

import (
	"net/http"

	"github.com/Krokozabra213/test_api/pkg/resp"
)

func HandleBody[T any](w http.ResponseWriter, req *http.Request) (*T, error) {
	body, err := Decode[T](req.Body)
	if err != nil {
		resp.JsonResp(w, 402, err.Error())

		return nil, err
	}
	err = IsValid[T](body)
	if err != nil {
		resp.JsonResp(w, 402, err.Error())

		return nil, err
	}
	return &body, nil
}
