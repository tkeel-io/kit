package http

import (
	"net/url"

	"github.com/emicklei/go-restful"
	"github.com/tkeel-io/kit/encoding"
)

func GetQuery(req *restful.Request, in interface{}) error {
	if err := encoding.GetCodec().Unmarshal([]byte(req.Request.URL.Query().Encode()), in); err != nil {
		return err
	}
	return nil
}

func GetPathValue(req *restful.Request, in interface{}) error {
	pathValue := req.PathParameters()
	vars := make(url.Values, len(pathValue))
	for k, v := range pathValue {
		vars[k] = []string{v}
	}
	if err := encoding.GetCodec().Unmarshal([]byte(vars.Encode()), in); err != nil {
		return err
	}
	return nil
}

func GetBody(req *restful.Request, in interface{}) error {
	return req.ReadEntity(in)
}
