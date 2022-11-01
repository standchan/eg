package common

import (
	"fmt"
	"go.pfgit.cn/letsgo/xdev"
)

type XConfig struct {
	ConnStr  string `key:"common.db_conn" commit:"false"`
	LogLevel string `key:"common.log_level" default:"INFO"`

	//TODO:定义模块自己需要的配置，字段如何定义可查看xdev.ReadConfig说明
	DataDir   string `key:"common.data_dir" default:"/opt/eagle_data"`
	NginxUrl  string `key:"common.nginx_url" default:"127.0.0.1:18420/eagle/download"`
	InputDir  string
	PkgsDir   string
	BackupDir string
	TmpDir    string
}

var Config XConfig

func ReadConfig() error {
	err := xdev.ReadConfig(APP_CONFIG_PATH, APP_NAME, &Config)
	if err != nil {
		return err
	} else {
		Config.InputDir = fmt.Sprintf("%s/%s", Config.DataDir, "input")
		Config.PkgsDir = fmt.Sprintf("%s/%s", Config.DataDir, "pkgs")
		Config.BackupDir = fmt.Sprintf("%s/%s", Config.DataDir, "backup")
		Config.TmpDir = fmt.Sprintf("%s/%s", Config.DataDir, "tmp")
	}
	return nil
}
