package dao

import (
	"github.com/volcengine/volc-sdk-golang/service/vikingdb"
	"log"
	"meetingagent/config"
)

var Vkdb *vikingdb.VikingDBService

func InitVkdb() {
	vkConfig := config.Cfg.VKDB
	Vkdb = vikingdb.NewVikingDBService(vkConfig.Host,
		vkConfig.Region,
		vkConfig.Ak,
		vkConfig.Sk,
		vkConfig.Scheme)
	log.Println("创建vikingdb向量数据库成功")
	Vkdb.SetConnectionTimeout(5)
}
