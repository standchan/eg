package core

import (
	"eagleserver/common"
	"fmt"
	"go.pfgit.cn/letsgo/xdev"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	scanInputDir(common.Config.DataDir + "/input")
}

func scanInputDir(dirPath string) {
	dirEntry, err := os.ReadDir(dirPath)
	if err != nil {
		xdev.Log.Error(err)
		return
	}
	for _, d := range dirEntry {
		fs, err := d.Info()
		if err != nil {
			xdev.Log.Error(err)
			continue
		} else {
			now := time.Now().Unix()
			if now-fs.ModTime().Unix() > 10 {
				processPkg(dirPath, d.Name())
			}
		}
	}
}

func processPkg(dirPath string, pkgName string) error {
	pkgPath := filepath.Join(dirPath, pkgName)
	xdev.Log.Info("find file,filepath= ", pkgPath)

	if strings.HasSuffix(pkgPath, ".tar.gz") {
		if nameFile := strings.Split(pkgPath, "-"); len(nameFile) == 4 {
			//模块打包规范：{project}-{module}-install-{date}.tar.gz
			//将符合要求的包移动到pkgs目录下
			newPkgPath := filepath.Join(common.Config.DataDir+"/pkgs", pkgName)
			err := os.Rename(pkgPath, newPkgPath)
			if err != nil {
				xdev.Log.Errorf("move pkg from %s to %s", pkgPath, newPkgPath)
			} else {
				xdev.Log.Errorf("get pkg,pkgpath= %s", newPkgPath)
			}
		} else { //非模块包格式，一般为项目包及其子分类包
			if nameFile := strings.Split(pkgPath, "-"); len(nameFile) == 3 {
				subdirName := fmt.Sprintf("%s-%s", nameFile[0], nameFile[1])
				subdirPath := fmt.Sprintf("%s/%s", dirPath, subdirName)
				//TODO 源代码中有删除一下subdirPath,这里暂时不做。
				//解压项目包
				err := decompress(pkgPath, subdirPath)
				if err != nil {
					xdev.Log.Error(err)
					return err
				}
				dirEntry, err := os.ReadDir(subdirPath)
				if err != nil {
					xdev.Log.Error(err)
					return err
				}
				var errNum int
				for _, d := range dirEntry {
					err := processPkg(subdirPath, d.Name())
					if err != nil {
						errNum++
					}
				}
				if errNum == 0 {
					os.Remove(subdirPath)
					os.Remove(pkgPath)
				}
			}
		}
	} else {
		xdev.Log.Infof("move file %s to %s", pkgPath)
		err := os.Rename(pkgPath, common.Config.BackupDir)
		if err != nil {
			xdev.Log.Error(err)
			return err
		}
	}
	return nil
}

func scanPkgsDir() {

}

func decompress(filePath string, destPath string) error {
	cmd := exec.Command("tar zxf", filePath, "-C", destPath)
	err := cmd.Run()
	if err != nil {
		xdev.Log.Error(err)
		return err
	}
	return nil
}

func installDB() {

}

func updatePkgInfo() {

}
