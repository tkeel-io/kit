/*
Copyright 2021 The tKeel Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package encoding

import (
	"encoding/json"
	"fmt"
	"github.com/tkeel-io/tdtl"
	"google.golang.org/protobuf/types/known/structpb"
)

func NewStructValue(cc *tdtl.Collect) (*structpb.Value, error) {
	switch cc.Type() {
	case tdtl.Bool:
		ret := cc.To(tdtl.Bool)
		switch ret := ret.(type) {
		case tdtl.BoolNode:
			return structpb.NewBoolValue(bool(ret)), nil
		}
	case tdtl.Int, tdtl.Float, tdtl.Number:
		ret := cc.To(tdtl.Number)
		switch ret := ret.(type) {
		case tdtl.IntNode:
			return structpb.NewNumberValue(float64(ret)), nil
		case tdtl.FloatNode:
			return structpb.NewNumberValue(float64(ret)), nil
		}
	case tdtl.String:
		return structpb.NewStringValue(cc.String()), nil
	case tdtl.JSON, tdtl.Object, tdtl.Array:
		ret := &structpb.Struct{}
		err := json.Unmarshal(cc.Raw(), ret)
		if err != nil {
			fmt.Println("err", err)
			return structpb.NewBoolValue(false), err
		}
		return structpb.NewStructValue(ret), nil
	}
	return structpb.NewBoolValue(false), fmt.Errorf("unknown type")
}

func NewCollectValue(val *structpb.Value) (*tdtl.Collect, error) {
	byt, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	ret := tdtl.New(byt)
	return ret, nil
}
