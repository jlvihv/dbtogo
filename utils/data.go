package utils

import (
	"fmt"
	"log"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/jlvihv/dbtogo/defines"
	"github.com/mitchellh/go-homedir"
)

var (
	configPath = ""
	config     *defines.Config
	once       sync.Once
)

func GetConfig() *defines.Config {
	once.Do(func() {
		if configPath == "" {
			defaultConfigPath := "~/.config/dbtogo/config.toml"
			absConfigPath, err := homedir.Expand(defaultConfigPath)
			if err != nil {
				fmt.Println("获取家目录失败，无法读取配置文件")
				log.Fatal(err)
			}
			configPath = absConfigPath
		}
		if _, err := toml.DecodeFile(configPath, &config); err != nil {
			fmt.Println("读取配置文件失败，请检查配置文件路径是否正确，以及配置文件格式是否书写正确")
			log.Fatal(err)
		}
	})
	return config
}

func ConfigPath() *string {
	return &configPath
}
