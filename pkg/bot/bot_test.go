package bot

import (
	"testing"
)

func Test_Bot(t *testing.T) {
	bot := Bot{
		Name: "unit_test",
		Setting: Setting{
			Legs: []Leg{
				{
					Underlying: "SPX",
					Action:     BUY,
					Type:       CALL,
					Quantity:   1,
					Strike:     Strike{},
				},
			},
		},
		EnableAutoOpen:  true,
		EnableAutoClose: true,
	}
	println(bot)
}
