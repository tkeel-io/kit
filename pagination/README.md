# Usage example

```go
import "github.com/tkeel.io/kit/pagination"


func (s *SubscribeService) ListSubscribeEntities(ctx context.Context, req *pb.ListSubscribeEntitiesRequest) (*pb.ListSubscribeEntitiesResponse, error) {
    page, err := pagination.Parse(req)
    
	// {Num:10 Size:50 OrderBy:test IsDescending:true KeyWords:key SearchKey:search}
    fmt.Println(page)

    // for example query like an ORM operator 
	query := DB.Query()
	
	if page.Requried {
		cond, fields := page.SearchCondition()
    if cond != nil {
		query.Where(cond)
	}
	if fields != nil {
        query.Select(fields...)
    }
	}
		// Do paginate here
        query.Limit(page.Limit())
        query.Offset(page.Offset())

	} else {
		// Query all date here
    }
	// Query data here
	data := query.Find()
	
	// count all after remove Offset and Limit
	count := query.Offset(-1).Limit(-1).Count()
	
	resp := &pb.ListSubscribeEntitiesResponse{}
    resp.Data = data
	
	// Fill Response Page Fields
    page.FillResponse(resp, count)
}

```

## func Required
Used to determine if the paging request passed meets the paging needs

Judgement conditions：
```go
func (p Page) Required() bool {
	return p.Num > 0 && p.Size > 0
}
```

## func Limit
if no limit set this will return the default value.
```go
func (p Page) Limit() uint32 {
	if p.Size != 0 {
		return uint32(p.Size)
	}

	return uint32(p.defaultSize)
}
```

## func Offset
count the offset of the current page
```go
func (p Page) Offset() uint32 {
	if p.Num <= 0 {
		return 0
	}
	return uint32((p.Num - 1) * p.Size)
}
```

## func SearchCondition
return a `map[string]string` for the search condition and a `[]string` for the search fields.

iIf the search have no condition or fields required, return `nil`
```go
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
```
### Usage
```go

    cond, fields := page.SearchCondition()
    if cond != nil {
        query = query.Where(cond)
    }
    if fields != nil {
        query = query.Select(fields...)
    }
```


## func FillResponse
Automatic padding of paginated data to match paginated responsive design.
```go
func (p Page) FillResponse(resp interface{}, total int) error {
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
				v.Field(i).SetUint(uint64(total))
			case "PageNum":
				v.Field(i).SetUint(uint64(p.Num))
			case "LastPage":
				if p.Size == 0 {
					v.Field(i).SetUint(uint64(0))
					continue
				}
				lastPage := total / int(p.Size)
				if total%int(p.Size) == 0 {
					v.Field(i).SetUint(uint64(lastPage))
					continue
				}
				v.Field(i).SetUint(uint64(lastPage + 1))

			case "PageSize":
				v.Field(i).SetUint(uint64(p.Size))
			}
		}
	}
	return nil
}
```