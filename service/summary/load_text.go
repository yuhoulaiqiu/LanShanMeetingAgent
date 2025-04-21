package summary

import (
	"context"
	"github.com/cloudwego/eino/compose"
	"time"
)

func LoadText() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input texts) (map[string]any, error) {
		output := make(map[string]any)
		output["unsummarized_text"] = input.unsummarized_text
		output["summarized_text"] = input.summarized_text
		currentTime := time.Now()
		output["current_time"] = currentTime.Format("2006-01-02 15:04:05")
		return output, nil
	})
}
