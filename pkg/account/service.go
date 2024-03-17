package account

import (
	"context"

	"cdr.dev/slog"
	"connectrpc.com/connect"
	"github.com/ppaanngggg/option-bot/pkg/util"
	v1 "github.com/ppaanngggg/option-bot/proto/gen/account/v1"
	"github.com/ppaanngggg/option-bot/proto/gen/account/v1/accountv1connect"
)

var Service accountv1connect.AccountServiceHandler

func init() {
	Service = &service{
		logger: util.DefaultLogger.With(slog.F("account", "service")),
	}
}

type service struct {
	logger slog.Logger
}

func (s *service) Create(
	ctx context.Context, req *connect.Request[v1.CreateRequest],
) (*connect.Response[v1.CreateResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) Get(
	ctx context.Context, req *connect.Request[v1.GetRequest],
) (*connect.Response[v1.GetResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) List(
	ctx context.Context, req *connect.Request[v1.ListRequest],
) (*connect.Response[v1.ListResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) Delete(
	ctx context.Context, req *connect.Request[v1.DeleteRequest],
) (*connect.Response[v1.DeleteResponse], error) {
	//TODO implement me
	panic("implement me")
}
