package dao

import (
	"fmt"
	"github.com/volcengine/volc-sdk-golang/service/vikingdb"
	"log"
	"meetingagent/config"
)

func InitVkdb() {
	vkConfig := config.Cfg.VKDB
	fmt.Println("vikingdb config", vkConfig)
	vkdb := vikingdb.NewVikingDBService(vkConfig.Host,
		vkConfig.Region,
		vkConfig.Ak,
		vkConfig.Sk,
		vkConfig.Scheme)
	log.Println("创建vikingdb向量数据库成功")
	vkdb.SetConnectionTimeout(5)
}
