package summary

import (
	"context"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func Store() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, msg *schema.Message) (string, error) {
		println(msg.Content)
		return msg.Content, nil
	})
}
