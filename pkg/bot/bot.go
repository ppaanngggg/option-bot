package bot

import botv1 "github.com/ppaanngggg/option-bot/proto/gen/bot/v1"

type Bot struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	EnableAutoOpen  bool           `json:"enable_auto_open"`
	EnableAutoClose bool           `json:"enable_auto_close"`
	Setting         *botv1.Setting `json:"setting"`
}
