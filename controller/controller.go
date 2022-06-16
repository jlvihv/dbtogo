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
	Error      error
	db         *gorm.DB
	dbName     string
	tableNames []string
	tables     []table
	structText string
}

type table struct {
	name          string
	columns       []*defines.TableColumn
	structColumns []*defines.StructColumn
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

func (self *controller) GetColumns(dbName, tableNames string) *controller {
	if self == nil || self.Error != nil || self.db == nil {
		return self
	}
	fmt.Println("获取表信息...")
	self.dbName = dbName
	self.splitTableNames(tableNames)
	var tables []table
	for _, tableName := range self.tableNames {
		var columns []*defines.TableColumn
		gormResult := self.db.Table("columns").Select([]string{"column_name", "data_type", "column_key", "is_nullable", "column_type", "column_comment"}).
			Where("table_schema = ? and table_name = ?", dbName, tableName).Find(&columns)
		if err := gormResult.Error; err != nil {
			fmt.Println(err)
			self.Error = err
			return self
		}
		if len(columns) == 0 {
			fmt.Printf("db: %s, table: %s 没有任何信息\n", self.dbName, tableName)
			continue
		}
		tables = append(tables, table{columns: columns, name: tableName})
	}
	self.tables = tables
	return self
}

func (self *controller) splitTableNames(tableNames string) *controller {
	if self == nil || self.Error != nil {
		return self
	}
	if strings.Contains(tableNames, ",") {
		tableNames = strings.ReplaceAll(tableNames, " ", "")
		tableNames := strings.Split(tableNames, ",")
		self.tableNames = tableNames
	} else {
		self.tableNames = []string{tableNames}
	}
	return self
}

func (self *controller) ConvertToStructColumns() *controller {
	if self == nil || self.Error != nil || len(self.tables) == 0 {
		return self
	}
	for i, table := range self.tables {
		structColumns := make([]*defines.StructColumn, 0, len(table.columns))
		for _, column := range table.columns {
			structColumns = append(structColumns, &defines.StructColumn{
				Name:    column.ColumnName,
				Type:    getStructType(column.DataType),
				Tag:     "",
				Comment: column.ColumnComment,
			})
		}
		self.tables[i].structColumns = structColumns
	}
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
	if self == nil || self.Error != nil || self.tables == nil || len(self.tables) == 0 {
		return self
	}
	for i := range self.tables {
		self.tables[i].name = utils.UnderscoreToUpperCamelCase(self.tables[i].name)
		for j, column := range self.tables[i].structColumns {
			self.tables[i].structColumns[j].Name = utils.UnderscoreToUpperCamelCase(column.Name)
		}
	}
	return self
}

func (self *controller) Generate() *controller {
	if self == nil || self.Error != nil || len(self.tables) == 0 {
		return self
	}
	result := make([]string, 0, 16)
	for _, table := range self.tables {
		result = append(result, strings.ReplaceAll(defines.StructTemplateText["title"], "STRUCT_NAME", table.name))
		for _, column := range table.structColumns {
			line := defines.StructTemplateText["line"]
			line = strings.ReplaceAll(line, "FIELD_NAME", column.Name)
			line = strings.ReplaceAll(line, "FIELD_TYPE", column.Type)
			if len(column.Tag) == 0 {
				line = strings.ReplaceAll(line, "FIELD_TAG", "")
			} else {
				line = strings.ReplaceAll(line, "FIELD_TAG", fmt.Sprintf("`%s`", column.Tag))
			}
			result = append(result, line)
		}
		result = append(result, defines.StructTemplateText["end"])
	}
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
	if self == nil || self.Error != nil || len(self.tables) == 0 || len(tag) == 0 {
		return self
	}
	for index, table := range self.tables {
		for i, v := range table.structColumns {
			var tagValue string
			switch tag {
			case "comment":
				if len(v.Comment) != 0 {
					tagValue = fmt.Sprintf("%s:\"%s\"", tag, v.Comment)
				}
			case "gorm":
				tagValue = fmt.Sprintf("%s:\"column:%s\"", tag, v.Name)
			default:
				tagValue = fmt.Sprintf("%s:\"%s\"", tag, v.Name)
			}
			if len(v.Tag) != 0 {
				if len(tagValue) == 0 {
					tagValue = v.Tag
				} else {
					tagValue = fmt.Sprintf("%s %s", v.Tag, tagValue)
				}
			}
			self.tables[index].structColumns[i].Tag = tagValue
		}
	}
	return self
}
