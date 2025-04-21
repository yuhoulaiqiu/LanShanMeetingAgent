package dao

import (
	"github.com/volcengine/volc-sdk-golang/service/vikingdb"
	"meetingagent/config"
)

func InitVkdb() {
	vkConfig := config.Cfg.VKDB
	vkdb := vikingdb.NewVikingDBService(vkConfig.Host,
		vkConfig.Region,
		vkConfig.Ak,
		vkConfig.Sk,
		vkConfig.Scheme)
	vkdb.SetConnectionTimeout(5)
}
