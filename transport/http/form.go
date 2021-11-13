package http

import (
	"fmt"

	"github.com/emicklei/go-restful"
	"github.com/tkeel-io/kit/encoding"
)

type FormEntityReadWriter struct{}

func (ferw FormEntityReadWriter) Read(req *restful.Request, v interface{}) error {
	if err := req.Request.ParseForm(); err != nil {
		return fmt.Errorf("error request parse form: %w", err)
	}
	if err := encoding.GetCodec().Unmarshal([]byte(req.Request.Form.Encode()), v); err != nil {
		return fmt.Errorf("error encoding unmarshal: %w", err)
	}

	return nil
}

func (ferw FormEntityReadWriter) Write(resp *restful.Response,
	status int, v interface{}) error {
	b, err := encoding.GetCodec().Marshal(v)
	if err != nil {
		return fmt.Errorf("error encoding marshal: %w", err)
	}
	resp.ResponseWriter.WriteHeader(status)
	_, err = resp.ResponseWriter.Write(b)
	if err != nil {
		return fmt.Errorf("error write response writer:%w", err)
	}
	return nil
}
