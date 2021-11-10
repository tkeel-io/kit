package encoding

import (
	"net/url"
	"reflect"
	"sync"

	"github.com/go-playground/form/v4"
	"google.golang.org/protobuf/proto"
)

const (
	nullStr = ""
)

var (
	once        *sync.Once
	globalCodec *codec
)

type codec struct {
	encoder *form.Encoder
	decoder *form.Decoder
}

func GetCodec() *codec {
	once.Do(func() {
		globalCodec = NewCodec()
	})
	return globalCodec
}

func NewCodec() *codec {
	return &codec{
		encoder: form.NewEncoder(),
		decoder: form.NewDecoder(),
	}
}

func (c codec) Marshal(v interface{}) ([]byte, error) {
	var vs url.Values
	var err error
	if m, ok := v.(proto.Message); ok {
		vs, err = EncodeMap(m)
		if err != nil {
			return nil, err
		}
	} else {
		vs, err = c.encoder.Encode(v)
		if err != nil {
			return nil, err
		}
	}
	for k, v := range vs {
		if len(v) == 0 {
			delete(vs, k)
		}
	}
	return []byte(vs.Encode()), nil
}

func (c codec) Unmarshal(data []byte, v interface{}) error {
	vs, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	if m, ok := v.(proto.Message); ok {
		return MapProto(m, vs)
	} else if m, ok := reflect.Indirect(reflect.ValueOf(v)).Interface().(proto.Message); ok {
		return MapProto(m, vs)
	}

	return c.decoder.Decode(v, vs)
}
