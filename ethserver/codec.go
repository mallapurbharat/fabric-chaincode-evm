package ethserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

type rpcCodec struct {
	codec *json.Codec
}

type codecRequest struct {
	rpc.CodecRequest
}

func NewRPCCodec() rpc.Codec {
	codec := json.NewCodec()
	return &rpcCodec{codec: codec}
}

func (c *rpcCodec) NewRequest(r *http.Request) rpc.CodecRequest {
	req := c.codec.NewRequest(r)
	return &codecRequest{req}
}

func (r *codecRequest) Method() (string, error) {
	m, err := r.CodecRequest.Method()
	if err != nil {
		return "", err
	}
	method := strings.Split(m, "_")
	return fmt.Sprintf("%s.%s", method[0], strings.Title(method[1])), nil
}
