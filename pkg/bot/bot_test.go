package bot

import (
	"context"
	"github.com/Azure/go-asynctask"
	"testing"
	"time"
)

func Test_AsyncTask(t *testing.T) {
	ctx := context.Background()
	task := asynctask.Start(
		ctx, asynctask.ActionToFunc(
			func(ctx context.Context) error {
				for {
					select {
					case <-ctx.Done():
						println("done")
						return nil
					default:
						time.Sleep(1 * time.Second)
						println(time.Now().UTC().String())
					}
				}
			},
		),
	)
	time.Sleep(5 * time.Second)
	task.Cancel()
	err := task.Wait(ctx)
	println(err.Error())
	time.Sleep(5 * time.Second)
}
