package defines

var StructTemplateText = map[string]string{
	"title": "type STRUCT_NAME struct {",
	"line": "	FIELD_NAME	FIELD_TYPE	FIELD_TAG",
	"end": "}",
}

// DBTypeToStructType 数据库数据类型到 go 结构体数据类型的转换规则
var DBTypeToStructType = map[string]string{
	"int":                "int",
	"integer":            "int",
	"tinyint":            "int8",
	"smallint":           "int16",
	"mediumint":          "int32",
	"bigint":             "int64",
	"int unsigned":       "uint",
	"integer unsigned":   "uint",
	"tinyint unsigned":   "uint8",
	"smallint unsigned":  "uint16",
	"mediumint unsigned": "uint32",
	"bigint unsigned":    "uint64",
	"bit":                "byte",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "time.Time",
	"datetime":           "time.Time",
	"timestamp":          "time.Time",
	"time":               "time.Time",
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
}

// TableColumn 数据库中字段信息
type TableColumn struct {
	ColumnName    string
	DataType      string
	ColumnKey     string
	IsNullable    string
	ColumnType    string
	ColumnComment string
}

// StructColumn go 结构体字段信息
type StructColumn struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}
