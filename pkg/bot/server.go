package bot

import (
	"context"

	"connectrpc.com/connect"
	botv1 "github.com/ppaanngggg/option-bot/proto/gen/bot/v1"
)

type Server struct{}

func (s *Server) Create(
	ctx context.Context, c *connect.Request[botv1.CreateRequest],
) (*connect.Response[botv1.CreateResponse], error) {
	panic("implement me")
}

func (s *Server) Get(
	ctx context.Context, c *connect.Request[botv1.GetRequest],
) (*connect.Response[botv1.GetResponse], error) {
	panic("implement me")
}
