package main

import (
	"fmt"
	"net/http"

	"github.com/ppaanngggg/option-bot/pkg/account"
	"github.com/ppaanngggg/option-bot/pkg/bot"
	"github.com/ppaanngggg/option-bot/pkg/datasource"
	"github.com/ppaanngggg/option-bot/pkg/util"
	"github.com/ppaanngggg/option-bot/proto/gen/account/v1/accountv1connect"
	"github.com/ppaanngggg/option-bot/proto/gen/bot/v1/botv1connect"
	"github.com/ppaanngggg/option-bot/proto/gen/datasource/v1/datasourcev1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	mux := http.NewServeMux()
	{
		path, handler := accountv1connect.NewAccountServiceHandler(account.Service)
		mux.Handle(path, handler)
	}
	{
		path, handler := botv1connect.NewBotServiceHandler(bot.Service)
		mux.Handle(path, handler)
	}
	{
		path, handler := datasourcev1connect.NewDataSourceServiceHandler(datasource.Service)
		mux.Handle(path, handler)
	}
	http.ListenAndServe(
		fmt.Sprintf("%s:%d", util.Conf.Server.Host, util.Conf.Server.Port),
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
