package summary

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/schema"
	"log"
	"meetingagent/service/utils"
)

func Split(meetingID string) {
	ctx := context.Background()
	// 初始化分割器
	splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":   "h1",
			"##":  "h2",
			"###": "h3",
		},
		TrimHeaders: false,
	})
	if err != nil {
		panic(err)
	}
	doc, _ := utils.ReadMeetingSummaryText(meetingID)
	docs := []*schema.Document{
		{
			ID:      meetingID,
			Content: doc,
		},
	}
	res, err := splitter.Transform(ctx, docs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
