package common

import "go.pfgit.cn/letsgo/xdev"

type XConfig struct {
	ConnStr  string `key:"common.db_conn" commit:"false"`
	LogLevel string `key:"common.log_level" default:"INFO"`

	//TODO:定义模块自己需要的配置，字段如何定义可查看xdev.ReadConfig说明
}

var Config XConfig

func ReadConfig() error {
	return xdev.ReadConfig(APP_CONFIG_PATH, APP_NAME, &Config)
}
