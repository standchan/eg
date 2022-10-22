package core

import (
	"eagleserver/common"
	"time"
)

type pkgInfo struct {
	appName string
	pkgName string
	version int
	pkgUrl  string
	mtime   int
	md5     string
	pkgType string
}

// 每60秒扫描一次包目录
// 防止包仍然在往input传输的时候被移动到pkg目录，导致包不完整情况发生。
func PkgScanThread() {
	timer := time.NewTimer(60 * time.Second)
	for {
		timer.Reset(60 * time.Second)
		select {
		case <-timer.C:
			pkgScan()
		}
	}
}

func pkgScan() {
	scan(common.Config.DataDir + "/input")
}

func scan(dirPath string) {

}

func checkCTime() {

}

func processPkgs() {

}

func installDB() {

}

func updatePkgInfo() {

}
