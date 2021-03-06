package cmd

import (
	"fmt"
	"log"

	"github.com/jlvihv/dbtogo/controller"
	"github.com/jlvihv/dbtogo/utils"
	"github.com/spf13/cobra"
)

var (
	db      string
	table   string
	json    bool
	toml    bool
	gorm    bool
	yaml    bool
	tag     string
	comment bool
	form    bool
	clip    bool
	file    string
)

var rootCmd = &cobra.Command{
	Use:   "dbtogo",
	Short: "将数据库中的表转换为 go 结构体",
	Run: func(_ *cobra.Command, _ []string) {
		run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&db, "db", "", "", "数据库名")
	rootCmd.Flags().StringVarP(&table, "table", "", "", "表名")

	rootCmd.Flags().BoolVarP(&json, "json", "j", false, "生成 json 标签")
	rootCmd.Flags().BoolVarP(&toml, "toml", "t", false, "生成 toml 标签")
	rootCmd.Flags().BoolVarP(&gorm, "gorm", "g", false, "生成 gorm 标签")
	rootCmd.Flags().BoolVarP(&form, "form", "f", false, "生成 form 标签")
	rootCmd.Flags().BoolVarP(&yaml, "yaml", "y", false, "生成 yaml 标签")
	rootCmd.Flags().BoolVarP(&comment, "comment", "c", true, "生成 comment 标签")
	rootCmd.Flags().StringVarP(&tag, "tag", "", "", "生成自定义标签")

	rootCmd.Flags().BoolVarP(&clip, "clip", "", false, "输出到系统剪贴板")
	rootCmd.Flags().StringVarP(&file, "file", "", "", "输出到文件")

	rootCmd.Flags().StringVarP(utils.ConfigPath(), "config", "", "", "指定配置文件所在位置")
}

func run() {
	if len(db) == 0 || len(table) == 0 {
		fmt.Printf("db: %s, table: %s\n", db, table)
		fmt.Println("数据库名与表名不得为空")
		return
	}
	c := controller.NewController()

	c.GetColumns(db, table).ConvertToStructColumns()

	if json {
		c.AddJsonTag()
	}
	if toml {
		c.AddTomlTag()
	}
	if form {
		c.AddFormTag()
	}
	if gorm {
		c.AddGormTag()
	}
	if yaml {
		c.AddYamlTag()
	}
	if comment {
		c.AddCommentTag()
	}
	if len(tag) != 0 {
		c.AddTag(tag)
	}

	c.ToUpperCamelCase().Generate()

	if !clip && len(file) == 0 {
		c.Stdout()
	}
	if clip {
		c.Clipboard()
	}
	if len(file) != 0 {
		c.File(file)
	}
}
