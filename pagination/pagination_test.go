package pagination

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testRequest struct {
	PageNum      int32
	PageSize     int32
	OrderBy      string
	IsDescending bool
	KeyWords     string
	SearchKey    string
	CustomField  string
}

var (
	pbData = ListRequest{
		PageNum:      10,
		PageSize:     50,
		OrderBy:      "test",
		IsDescending: true,
		KeyWords:     "key",
		SearchKey:    "search",
		PacketId:     0,
	}
	customData = testRequest{
		PageNum:      10,
		PageSize:     50,
		OrderBy:      "test",
		IsDescending: true,
		KeyWords:     "key",
		SearchKey:    "search",
		CustomField:  "my data",
	}
	targetPage = Page{
		Num:              10,
		Size:             50,
		OrderBy:          "test",
		IsDescending:     true,
		KeyWords:         "key",
		SearchKey:        "search",
		defaultSize:      15,
		defaultSeparator: ",",
	}
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		options  []Option
		excepted Page
	}{
		{
			name:     "test pb struct",
			data:     pbData,
			excepted: targetPage,
		},
		{
			name:     "test custom data",
			data:     customData,
			excepted: targetPage,
		},
		{
			name:     "test pb ptr",
			data:     &pbData,
			excepted: targetPage,
		},
		{
			name:     "test custom ptr",
			data:     &customData,
			excepted: targetPage,
		},
		{
			name:     "test custom ptr with default options",
			data:     &customData,
			options:  []Option{WithDefaultSize(20), WithDefaultSeparator(";")},
			excepted: targetPage,
		},
	}
	for _, test := range tests {
		page, err := Parse(test.data, test.options...)
		assert.NoError(t, err)
		if test.name == "test custom ptr with default options" {
			copyOne := targetPage
			copyOne.defaultSize = 20
			copyOne.defaultSeparator = ";"
			assert.Equal(t, copyOne, page)
		} else {
			assert.Equal(t, test.excepted, page)
		}
	}
}

func TestInvalidDataTypeParse(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		excepted error
	}{
		{
			name:     "map type",
			data:     make(map[string]interface{}),
			excepted: ErrInvalidParseData,
		},
		{
			name:     "string",
			data:     "string",
			excepted: ErrInvalidParseData,
		},
		{
			name:     "int",
			data:     32,
			excepted: ErrInvalidParseData,
		},
		{
			name:     "float",
			data:     float32(32.0),
			excepted: ErrInvalidParseData,
		},
		{
			name:     "bool",
			data:     true,
			excepted: ErrInvalidParseData,
		},
		{
			name:     "struct",
			data:     struct{ invalidField string }{"test"},
			excepted: nil,
		},
		{
			name:     "mismatch struct type",
			data:     struct{ PageNum string }{"PageNum"},
			excepted: ErrInvalidPageNum,
		},
	}
	for _, test := range tests {
		_, err := Parse(test.data)
		assert.Equal(t, test.excepted, err)
	}
}

func TestPage_Offset(t *testing.T) {
	tests := []struct {
		name     string
		req      ListRequest
		excepted uint32
	}{
		{
			name:     "test offset",
			req:      pbData,
			excepted: 450,
		},
		{
			name:     "test offset no page",
			req:      ListRequest{},
			excepted: 0,
		},
		{
			name:     "test offset page 1 per page 30",
			req:      ListRequest{PageNum: 1, PageSize: 30},
			excepted: 0,
		},
		{
			name:     "test offset page 2 per page 30",
			req:      ListRequest{PageNum: 2, PageSize: 30},
			excepted: 30,
		},
	}

	for _, test := range tests {
		page, err := Parse(&test.req)
		assert.NoError(t, err)
		assert.Equal(t, test.excepted, page.Offset())
	}
}

func TestPage_Limit(t *testing.T) {
	tests := []struct {
		name     string
		req      ListRequest
		excepted int32
	}{
		{
			name:     "test limit",
			req:      pbData,
			excepted: 50,
		},
		{
			name:     "test limit with default when page query",
			req:      ListRequest{PageNum: 1},
			excepted: 15,
		},
		{
			name:     "test limit with query all data",
			req:      ListRequest{},
			excepted: 0,
		},
	}

	for _, test := range tests {
		page, err := Parse(&test.req)
		assert.NoError(t, err)
		assert.Equal(t, test.excepted, page.Limit())
	}
}

func TestPage_Required(t *testing.T) {
	tests := []struct {
		name     string
		req      ListRequest
		excepted bool
	}{
		{
			name:     "test Required",
			req:      pbData,
			excepted: true,
		},
		{
			name:     "test no required",
			req:      ListRequest{},
			excepted: false,
		},
	}

	for _, test := range tests {
		page, err := Parse(&test.req)
		assert.NoError(t, err)
		assert.Equal(t, test.excepted, page.Required())
	}
}

func TestPage_SearchCondition(t *testing.T) {
	type excepted struct {
		condition map[string]string
		fields    []string
	}
	tests := []struct {
		name     string
		req      ListRequest
		excepted excepted
	}{
		{
			name: "test search condition: only condition",
			req:  pbData,
			excepted: excepted{
				condition: map[string]string{"search": "key"},
				fields:    nil,
			},
		},
		{
			name: "test only fields",
			req:  ListRequest{SearchKey: "test0, test1, test2"},
			excepted: excepted{
				condition: nil,
				fields:    []string{"test0", "test1", "test2"},
			},
		},
		{
			name: "test both required",
			req:  ListRequest{SearchKey: "1,2", KeyWords: "1"},
			excepted: excepted{
				condition: map[string]string{"1": "1"},
				fields:    []string{"2"},
			},
		},
	}

	for _, test := range tests {
		page, err := Parse(&test.req)
		cond, fields := page.SearchCondition()
		assert.NoError(t, err)
		assert.Equal(t, test.excepted.condition, cond)
		assert.Equal(t, test.excepted.fields, fields)
	}
}

func TestPage_FillResponse(t *testing.T) {
	tests := []struct {
		name     string
		req      ListRequest
		count    uint
		excepted ListResponse
	}{
		{
			name:  "test Required",
			req:   pbData,
			count: 500,
			excepted: ListResponse{
				Total:    500,
				PageSize: pbData.PageSize,
				PageNum:  pbData.PageNum,
				LastPage: 10,
			},
		}, {
			name:  "test Required and more data",
			req:   pbData,
			count: 501,
			excepted: ListResponse{
				Total:    501,
				PageSize: pbData.PageSize,
				PageNum:  pbData.PageNum,
				LastPage: 11,
			},
		},
		{
			name:  "test no required",
			req:   ListRequest{},
			count: 500,
			excepted: ListResponse{
				Total:    500,
				PageSize: 500,
				PageNum:  0,
				LastPage: 0,
			},
		},
	}

	for _, test := range tests {
		resp := &ListResponse{}
		page, err := Parse(&test.req)
		page.SetTotal(test.count)
		err = page.FillResponse(resp)
		assert.NoError(t, err)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(test.excepted, *resp))
	}
}
