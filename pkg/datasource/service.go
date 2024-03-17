package datasource

import (
	"context"

	"cdr.dev/slog"
	"connectrpc.com/connect"
	"github.com/ppaanngggg/option-bot/pkg/util"
	v1 "github.com/ppaanngggg/option-bot/proto/gen/datasource/v1"
	"github.com/ppaanngggg/option-bot/proto/gen/datasource/v1/datasourcev1connect"
	"golang.org/x/xerrors"
)

var Service datasourcev1connect.DataSourceServiceHandler

func init() {
	Service = &service{
		logger: util.DefaultLogger.With(slog.F("datasource", "service")),
	}
}

type service struct {
	globalDataSource *DataSource
	logger           slog.Logger
}

func (s *service) SetGlobal(
	ctx context.Context, c *connect.Request[v1.SetGlobalRequest],
) (*connect.Response[v1.SetGlobalResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) checkGlobalDataSource() error {
	if s.globalDataSource == nil {
		return xerrors.New("global data source is not initialized")
	}
	return nil
}

func (s *service) SearchSymbols(
	ctx context.Context, req *connect.Request[v1.SearchSymbolsRequest],
) (*connect.Response[v1.SearchSymbolsResponse], error) {
	if err := s.checkGlobalDataSource(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	symbols, err := s.globalDataSource.Search(ctx, req.Msg.Query)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &connect.Response[v1.SearchSymbolsResponse]{
		Msg: &v1.SearchSymbolsResponse{
			Symbols: symbols,
		},
	}, nil
}

func (s *service) GetOptionExpirations(
	ctx context.Context, req *connect.Request[v1.GetOptionExpirationsRequest],
) (*connect.Response[v1.GetOptionExpirationsResponse], error) {
	if err := s.checkGlobalDataSource(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	expirations, err := s.globalDataSource.GetOptionExpirations(ctx, req.Msg.Underlying)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &connect.Response[v1.GetOptionExpirationsResponse]{
		Msg: &v1.GetOptionExpirationsResponse{
			Expirations: expirations,
		},
	}, nil
}

func (s *service) GetOptionChains(
	ctx context.Context, req *connect.Request[v1.GetOptionChainsRequest],
) (*connect.Response[v1.GetOptionChainsResponse], error) {
	if err := s.checkGlobalDataSource(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	chains, err := s.globalDataSource.GetOptionChains(
		ctx, req.Msg.Underlying, req.Msg.Expiration,
	)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &connect.Response[v1.GetOptionChainsResponse]{
		Msg: &v1.GetOptionChainsResponse{
			Chains: chains,
		},
	}, nil
}
