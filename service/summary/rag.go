package summary

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/philippgille/chromem-go"
	"log"
	"meetingagent/config"
	"meetingagent/dao"
	"meetingagent/service/utils"
)

func Split(meetingID string) error {
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
		log.Println(err)
		return err
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
		log.Println(err)
		return err
	}
	collection, err := dao.ChromemDB.GetOrCreateCollection("rag", nil, chromem.NewEmbeddingFuncOpenAICompat(
		"https://ark.cn-beijing.volces.com/api/v3",
		config.Cfg.ModelInfo.ApiKey,
		"doubao-embedding-large-text-240915",
		nil,
	))
	if err != nil {
		log.Println("获取或创建集合失败")
		return err
	}

	// 删除与当前meetingID相关的所有旧文档
	where := map[string]string{"source_meeting_id": meetingID}
	err = collection.Delete(ctx, where, nil)
	if err != nil {
		log.Printf("删除旧文档时发生错误：%v\n", err)
	}

	for i, v := range res {
		fmt.Println("片段", i+1, ":", v.Content)
		doc1 := chromem.Document{
			ID:      "rag" + uuid.New().String(),
			Content: v.Content,
			Metadata: map[string]string{
				"source":            "文档",
				"source_meeting_id": meetingID, // 添加meetingID到元数据中以便后续删除
			},
		}
		err = collection.AddDocument(ctx, doc1)
		if err != nil {
			log.Println("添加文档失败")
			return err
		}
	}
	log.Println("文档加载完毕")
	return nil
}
