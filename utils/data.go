package utils

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jlvihv/dbtogo/defines"
	"log"
	"sync"
)

var (
	configPath = "config.toml"
	config     *defines.Config
	once       sync.Once
)

func GetConfig() *defines.Config {
	once.Do(func() {
		if _, err := toml.DecodeFile(configPath, &config); err != nil {
			fmt.Println("读取配置文件失败, 请检查当前目录下是否存在 config.toml 文件, 以及配置文件格式是否书写正确")
			log.Fatal(err)
		}
	})
	return config
}
