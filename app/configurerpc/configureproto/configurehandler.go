package configureproto

import (
	"context"

	"github.com/weblazy/core/timex"
)

type (
	ConfigureHandler struct {
	}
)

func NewConfigureHandler() (*ConfigureHandler, error) {
	return &ConfigureHandler{}, nil
}

func (h *ConfigureHandler) Rpc(_ context.Context, req *RpcRequest) (*RpcResponse, error) {
	resp := new(RpcResponse)
	resp.Source = timex.NowTimeStr()
	return resp, nil
}

func (h *ConfigureHandler) Server(_ context.Context, req *ServerRequest) (*ServerResponse, error) {
	resp := new(ServerResponse)
	return resp, nil
}

func (h *ConfigureHandler) Source(_ context.Context, req *SourceRequest) (*SourceResponse, error) {
	resp := new(SourceResponse)
	return resp, nil
}
