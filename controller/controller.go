package controller

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/jlvihv/dbtogo/defines"
	"github.com/jlvihv/dbtogo/gorm_db"
	"github.com/jlvihv/dbtogo/utils"
	"gorm.io/gorm"
	"io/ioutil"
	"strings"
)

type controller struct {
	Error         error
	db            *gorm.DB
	dbName        string
	tableName     string
	columns       []*defines.TableColumn
	structColumns []*defines.StructColumn
	structText    string
}

func NewController() *controller {
	ctl := &controller{}
	fmt.Println("连接数据库...")
	db, err := gorm_db.NewGormDB(&utils.GetConfig().DB)
	if err != nil {
		fmt.Println("连接数据库失败")
		ctl.Error = err
		return ctl
	}
	ctl.db = db
	return ctl
}

func (self *controller) GetColumns(dbName, tableName string) *controller {
	if self == nil || self.Error != nil || self.db == nil {
		return self
	}
	fmt.Println("获取表信息...")
	self.dbName, self.tableName = dbName, tableName
	var columns []*defines.TableColumn
	gormResult := self.db.Table("columns").Select([]string{"column_name", "data_type", "column_key", "is_nullable", "column_type", "column_comment"}).Where("table_schema = ? and table_name = ?", dbName, tableName).Find(&columns)
	if err := gormResult.Error; err != nil {
		fmt.Println(err)
		self.Error = err
		return self
	}
	if len(columns) == 0 {
		fmt.Printf("db: %s, table: %s 没有任何信息\n", self.dbName, self.tableName)
		return self
	}
	self.columns = columns
	return self
}

func (self *controller) ConvertToStructColumns() *controller {
	if self == nil || self.Error != nil || len(self.columns) == 0 {
		return self
	}
	structColumns := make([]*defines.StructColumn, 0, len(self.columns))
	for _, column := range self.columns {
		structColumns = append(structColumns, &defines.StructColumn{
			Name:    column.ColumnName,
			Type:    getStructType(column.DataType),
			Tag:     "",
			Comment: column.ColumnComment,
		})
	}
	self.structColumns = structColumns
	return self
}

func getStructType(dbType string) string {
	t, ok := defines.DBTypeToStructType[dbType]
	if !ok {
		t = "unknown"
	}
	return t
}

func (self *controller) ToUpperCamelCase() *controller {
	if self == nil || self.Error != nil || self.structColumns == nil || len(self.structColumns) == 0 {
		return self
	}
	self.tableName = utils.UnderscoreToUpperCamelCase(self.tableName)
	for i, v := range self.structColumns {
		self.structColumns[i].Name = utils.UnderscoreToUpperCamelCase(v.Name)
	}
	return self
}

func (self *controller) Generate() *controller {
	if self == nil || self.Error != nil || len(self.structColumns) == 0 {
		return self
	}
	var result []string
	result = append(result, strings.ReplaceAll(defines.StructTemplateText["title"], "STRUCT_NAME", self.tableName))
	for _, v := range self.structColumns {
		line := defines.StructTemplateText["line"]
		line = strings.ReplaceAll(line, "FIELD_NAME", v.Name)
		line = strings.ReplaceAll(line, "FIELD_TYPE", v.Type)
		if len(v.Tag) == 0 {
			line = strings.ReplaceAll(line, "FIELD_TAG", "")
		} else {
			line = strings.ReplaceAll(line, "FIELD_TAG", fmt.Sprintf("`%s`", v.Tag))
		}
		result = append(result, line)
	}
	result = append(result, defines.StructTemplateText["end"])
	self.structText = strings.Join(result, "\n")
	return self
}

func (self *controller) Stdout() *controller {
	if self == nil || self.Error != nil || len(self.structText) == 0 {
		return self
	}
	fmt.Println(self.structText)
	return self
}

func (self *controller) File(filename string) *controller {
	if self == nil || self.Error != nil || len(self.structText) == 0 {
		return self
	}
	fmt.Printf("输出到文件 %s ...", filename)
	err := ioutil.WriteFile(filename, []byte(self.structText), 0644)
	if err != nil {
		fmt.Println("失败")
		self.Error = err
		return self
	}
	fmt.Println("成功")
	return self
}

func (self *controller) String() string {
	if self == nil || self.Error != nil || len(self.structText) == 0 {
		return ""
	}
	fmt.Println("输出到标准输出")
	return self.structText
}

func (self *controller) Clipboard() {
	if self == nil || self.Error != nil || len(self.structText) == 0 {
		return
	}
	fmt.Print("输出到系统剪贴板...")
	err := clipboard.WriteAll(self.structText)
	if err != nil {
		fmt.Printf("\n输出到剪贴板失败 error: %s\n", err)
		fmt.Println("请手动复制")
		self.Stdout()
		return
	}
	fmt.Println("成功")
}

func (self *controller) AddJsonTag() *controller {
	return self.AddTag("json")
}

func (self *controller) AddTomlTag() *controller {
	return self.AddTag("toml")
}

func (self *controller) AddGormTag() *controller {
	return self.AddTag("gorm")
}

func (self *controller) AddCommentTag() *controller {
	return self.AddTag("comment")
}

func (self *controller) AddFormTag() *controller {
	return self.AddTag("form")
}

func (self *controller) AddYamlTag() *controller {
	return self.AddTag("yaml")
}

func (self *controller) AddTag(tag string) *controller {
	if self == nil || self.Error != nil || len(self.structColumns) == 0 || len(tag) == 0 {
		return self
	}
	for i, v := range self.structColumns {
		var tagValue string
		if tag == "comment" {
			if len(v.Comment) != 0 {
				tagValue = fmt.Sprintf("%s:\"%s\"", tag, v.Comment)
			}
		} else if tag == "gorm" {
			tagValue = fmt.Sprintf("%s:\"column:%s\"", tag, v.Name)
		} else {
			tagValue = fmt.Sprintf("%s:\"%s\"", tag, v.Name)
		}
		if len(v.Tag) != 0 {
			if len(tagValue) == 0 {
				tagValue = v.Tag
			} else {
				tagValue = fmt.Sprintf("%s %s", v.Tag, tagValue)
			}
		}
		self.structColumns[i].Tag = tagValue
	}
	return self
}
