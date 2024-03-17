package bot

import (
	"context"

	"cdr.dev/slog"
	"connectrpc.com/connect"
	"github.com/ppaanngggg/option-bot/pkg/util"
	botv1 "github.com/ppaanngggg/option-bot/proto/gen/bot/v1"
	"github.com/ppaanngggg/option-bot/proto/gen/bot/v1/botv1connect"
)

var Service botv1connect.BotServiceHandler

func init() {
	Service = &service{
		logger: util.DefaultLogger.With(slog.F("bot", "service")),
	}
}

type service struct {
	logger slog.Logger
}

func (s *service) Create(
	ctx context.Context, c *connect.Request[botv1.CreateRequest],
) (*connect.Response[botv1.CreateResponse], error) {
	panic("implement me")
}

func (s *service) Get(
	ctx context.Context, c *connect.Request[botv1.GetRequest],
) (*connect.Response[botv1.GetResponse], error) {
	panic("implement me")
}
