package main

import (
	"os"

	"github.com/ca17/teamsedge/config"
	"gopkg.in/yaml.v2"
)

// 初始化一个本地开发配置文件, 不会提交到git仓库， 可以本地随意修改

func main() {
	bs, err := yaml.Marshal(config.DefaultAppConfig)
	if err != nil {
		panic(err)
	}
	os.WriteFile("teamsedge.yml", bs, 777)
}
