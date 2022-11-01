package core

import (
	"eagleserver/common"
	"fmt"
	"go.pfgit.cn/letsgo/xdev"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type pkgInfo struct {
	appName string
	pkgName string
	version int
	pkgUrl  string
	mtime   time.Time
	md5     string
	pkgType string
}

var PkgsInfo []pkgInfo

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
	scanInputDir(common.Config.InputDir)
	scanPkgsDir(common.Config.PkgsDir)
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
				processInputDir(dirPath, d.Name())
			}
		}
	}
}

func processInputDir(dirPath string, pkgName string) error {
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
					err := processInputDir(subdirPath, d.Name())
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
		xdev.Log.Infof("move file %s to %s", pkgPath, common.Config.BackupDir)
		err := os.Rename(pkgPath, common.Config.BackupDir)
		if err != nil {
			xdev.Log.Error(err)
			return err
		}
	}
	return nil
}

func scanPkgsDir(dirPath string) {
	dirEntry, err := os.ReadDir(dirPath)
	if err != nil {
		xdev.Log.Error(err)
		return
	}
	processPkgsDir(dirEntry)
}

func processPkgsDir(dirEntry []os.DirEntry) {
	var pkgsInfo []pkgInfo
	for _, d := range dirEntry {
		fs, err := d.Info()
		if err != nil {
			xdev.Log.Error(err)
			continue
		}
		// fs.Name() eg: vendor-kafka-install-20220826.tar.gz
		if strings.HasSuffix(fs.Name(), ".tar.gz") {
			var tpkgInfo pkgInfo
			nameFile := strings.FieldsFunc(fs.Name(), Split)
			tpkgInfo.appName = nameFile[1]
			tpkgInfo.pkgName = fs.Name()
			versionInt, err := strconv.Atoi(nameFile[3])
			if err != nil {
				xdev.Log.Error(err)
				continue
			}
			tpkgInfo.version = versionInt
			tpkgInfo.md5 = CallMd5(fmt.Sprintf("%s/%s", common.Config.PkgsDir, fs.Name()))
			tpkgInfo.pkgUrl = fmt.Sprintf("%s/%s", common.Config.NginxUrl, fs.Name())
			tpkgInfo.mtime = fs.ModTime()
			tpkgInfo.pkgType = getPkgType(fmt.Sprintf("%s/%s", common.Config.PkgsDir, fs.Name()))
			pkgsInfo = append(pkgsInfo, tpkgInfo)
		}
	}
}

func Split(r rune) bool {
	return r == '-' || r == '.'
}

func SplitString(s string) []string {
	return strings.FieldsFunc(s, Split)
}

func CallMd5(filePath string) string {
	return ""
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

func getPkgType(filePath string) string {
	return ""
}

func installDB() {

}

func updatePkgInfo() {

}
