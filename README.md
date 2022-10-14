# eagleserver

go代码模板

common下方业务基础性代码，如配置、命令行、抽取的一些业务公共方法等  
对于代码量不大的情况，业务代码直接放在cmd目录下，不需要再增加其他目录

## vscode调试
直接在cmd/main.go按F5后，程序使用的启动路径是cmd，导致不能正确读取相对路径。  
配置launch.json可以一劳永逸的解决，还可以在launch.json里面设置env进行自测。

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "golang",
            "type": "go",    
            "request": "launch", 
            "mode": "auto", 
            "env": { 
                "common.db_conn": "titan/xianwei@172.16.5.22:3306@titandb"
            },
            "program": "${workspaceFolder}/cmd" ,
            "cwd": "${workspaceFolder}"
        }
    ]
}
```

