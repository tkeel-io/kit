package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tkeel-io/kit/encoding/testdata"
)

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type TestModel struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func TestFormCodecMarshal(t *testing.T) {
	req := &LoginRequest{
		Username: "tkeel",
		Password: "tkeel_pwd",
	}
	content, err := NewCodec().Marshal(req)
	require.NoError(t, err)
	require.Equal(t, []byte("password=tkeel_pwd&username=tkeel"), content)

	req = &LoginRequest{
		Username: "tkeel",
		Password: "",
	}
	content, err = NewCodec().Marshal(req)
	require.NoError(t, err)
	require.Equal(t, []byte("username=tkeel"), content)

	m := &TestModel{
		ID:   1,
		Name: "tkeel",
	}
	content, err = NewCodec().Marshal(m)
	t.Log(string(content))
	require.NoError(t, err)
	require.Equal(t, []byte("id=1&name=tkeel"), content)
}

func TestFormCodecUnmarshal(t *testing.T) {
	req := &LoginRequest{
		Username: "tkeel",
		Password: "tkeel_pwd",
	}
	content, err := NewCodec().Marshal(req)
	require.NoError(t, err)

	bindReq := new(LoginRequest)
	err = NewCodec().Unmarshal(content, bindReq)
	require.NoError(t, err)
	require.Equal(t, "tkeel", bindReq.Username)
	require.Equal(t, "tkeel_pwd", bindReq.Password)
}

func TestProtoEncodeDecode(t *testing.T) {
	in := &testdata.TestData{
		A: "A",
		B: 2,
		C: false,
		D: 4.4,
	}
	content, err := NewCodec().Marshal(in)
	require.NoError(t, err)
	require.Equal(t, "a=A&b=2&c=false&d=4.4", string(content))
	in2 := &testdata.TestData{}
	err = NewCodec().Unmarshal(content, in2)
	require.NoError(t, err)
	require.Equal(t, int32(2), in2.B)
	require.Equal(t, "A", in2.A)
	require.Equal(t, false, in2.C)
	require.Equal(t, float32(4.4), in2.D)
}
