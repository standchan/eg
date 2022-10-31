package main

import (
	"eagleserver/common"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"go.pfgit.cn/letsgo/xdev"
)

func main() {
	Log := xdev.XNewLoggerDefault("./log/" + common.APP_NAME + ".log")
	Log.Info(common.APP_NAME + "/" + common.APP_VERSION)

	if err := common.ReadConfig(); err != nil {
		os.Exit(1)
	}

	Log.Info("==========================config init ok=============================")
	Log.Infof("%+v", common.Config)
	xdev.SetLogLevel(common.Config.LogLevel)
	//TODO:
}
