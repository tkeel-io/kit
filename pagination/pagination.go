package pagination

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrInvalidPageNum      = errors.New("invalid page number")
	ErrInvalidPageSize     = errors.New("invalid page size")
	ErrInvalidOrderBy      = errors.New("invalid order")
	ErrInvalidKeyWords     = errors.New("invalid key words")
	ErrInvalidSearchKey    = errors.New("invalid search key")
	ErrInvalidIsDescending = errors.New("invalid is descending")
	ErrInvalidParseData    = errors.New("invalid data type parsing")
	ErrInvalidResponse     = errors.New("invalid response")
	ErrNoTotal             = errors.New("no total data")
)

type Page struct {
	Num              uint
	Size             uint
	OrderBy          string
	IsDescending     bool
	KeyWords         string
	SearchKey        string
	Total            uint
	defaultSize      uint
	defaultSeparator string
}

func (p Page) Offset() uint32 {
	if p.Num <= 0 {
		return 0
	}
	return uint32((p.Num - 1) * p.Size)
}

func (p Page) Limit() uint32 {
	if p.Size != 0 {
		return uint32(p.Size)
	}

	if p.Num != 0 {
		return uint32(p.defaultSize)
	}

	return 0
}

func (p Page) SearchCondition() (map[string]string, []string) {
	var values []string
	var keys []string
	if len(p.KeyWords) != 0 {
		values = strings.Split(p.KeyWords, p.defaultSeparator)
	}
	if len(p.SearchKey) != 0 {
		keys = strings.Split(p.SearchKey, p.defaultSeparator)
	}

	valLen := len(values)
	cond := make(map[string]string, valLen)
	fields := make([]string, 0, len(keys)-valLen)
	for i := range keys {
		if i < valLen {
			value := strings.TrimSpace(values[i])
			cond[keys[i]] = value
		} else {
			value := strings.TrimSpace(keys[i])
			fields = append(fields, value)
		}
	}

	if len(fields) == 0 {
		fields = nil
	}

	if len(cond) == 0 {
		cond = nil
	}

	return cond, fields
}

func (p Page) Required() bool {
	return p.Num > 0 && p.Size > 0
}

func (p *Page) SetTotal(total uint) {
	p.Total = total
}

func (p Page) FillResponse(resp interface{}) error {
	if p.Total == 0 {
		return ErrNoTotal
	}

	t := reflect.TypeOf(resp)
	v := reflect.ValueOf(resp)
	for t.Kind() != reflect.Struct {
		switch t.Kind() {
		case reflect.Ptr:
			v = v.Elem()
			t = t.Elem()
		default:
			return ErrInvalidResponse
		}
	}

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			switch t.Field(i).Name {
			case "Total":
				v.Field(i).SetUint(uint64(p.Total))
			case "PageNum":
				v.Field(i).SetUint(uint64(p.Num))
			case "LastPage":
				if p.Size == 0 {
					v.Field(i).SetUint(uint64(0))
					continue
				}
				lastPage := p.Total / p.Size
				if p.Total%p.Size == 0 {
					v.Field(i).SetUint(uint64(lastPage))
					continue
				}
				v.Field(i).SetUint(uint64(lastPage + 1))

			case "PageSize":
				if p.Size == 0 {
					v.Field(i).SetUint(uint64(p.Total))
					continue
				}
				v.Field(i).SetUint(uint64(p.Size))
			}
		}
	}
	return nil
}

type Option func(*Page) error

// Parse a struct which have defined Page fields.
func Parse(req interface{}, options ...Option) (Page, error) {
	q := Page{
		Num:              0,
		Size:             0,
		OrderBy:          "",
		IsDescending:     false,
		KeyWords:         "",
		SearchKey:        "",
		defaultSize:      15,
		defaultSeparator: ",",
	}
	v := reflect.ValueOf(req)
	t := reflect.TypeOf(req)
	for t.Kind() != reflect.Struct {
		switch t.Kind() {
		case reflect.Ptr:
			v = v.Elem()
			t = t.Elem()
		default:
			return q, ErrInvalidParseData
		}
	}

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			switch t.Field(i).Name {
			case "PageNum":
				if val, ok := v.Field(i).Interface().(uint); ok {
					q.Num = val
					continue
				}
				if val, ok := v.Field(i).Interface().(uint32); ok {
					q.Num = uint(val)
					continue
				}
				if val, ok := v.Field(i).Interface().(uint64); ok {
					q.Num = uint(val)
					continue
				}
				return q, ErrInvalidPageNum
			case "PageSize":
				if val, ok := v.Field(i).Interface().(uint); ok {
					q.Size = val
					continue
				}
				if val, ok := v.Field(i).Interface().(uint32); ok {
					q.Size = uint(val)
					continue
				}
				if val, ok := v.Field(i).Interface().(uint64); ok {
					q.Size = uint(val)
					continue
				}
				return q, ErrInvalidPageSize
			case "OrderBy":
				if val, ok := v.Field(i).Interface().(string); ok {
					q.OrderBy = val
				} else {
					return q, ErrInvalidOrderBy
				}
			case "IsDescending":
				if val, ok := v.Field(i).Interface().(bool); ok {
					q.IsDescending = val
				} else {
					return q, ErrInvalidIsDescending
				}
			case "KeyWords":
				if val, ok := v.Field(i).Interface().(string); ok {
					q.KeyWords = val
				} else {
					return q, ErrInvalidKeyWords
				}
			case "SearchKey":
				if val, ok := v.Field(i).Interface().(string); ok {
					q.SearchKey = val
				} else {
					return q, ErrInvalidSearchKey
				}
			}
		}
	}

	for i := range options {
		if err := options[i](&q); err != nil {
			return q, err
		}
	}

	return q, nil
}

func WithDefaultSize(size uint) Option {
	return func(p *Page) error {
		p.defaultSize = size
		return nil
	}
}

func WithDefaultSeparator(separator string) Option {
	return func(p *Page) error {
		p.defaultSeparator = separator
		return nil
	}
}
