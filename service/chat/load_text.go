package chat

import (
	"context"
	"github.com/cloudwego/eino/compose"
	"time"
)

func LoadText() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		output := input
		output["context"] = ""
		currentTime := time.Now()
		output["current_time"] = currentTime.Format("2006-01-02 15:04:05")
		return output, nil
	})
}
