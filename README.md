# dbtogo

这是一个将数据库中存在的表转换为 go 结构体的命令行工具, 应当是同类小工具中完成度较高, 实现较优雅的一个, 代码简洁清晰, 易于理解, 欢迎阅读, 哈哈

## 安装

```bash
go install github.com/jlvihv/dbtogo@0.0.1
```

## 使用

使用前, 必须在 config.toml 配置文件中配置数据库连接信息:

可以使用 --config 选项指定配置文件位置

配置示例: config.toml

```toml
[db]
ip = "127.0.0.1"
port = "3306"
username = "root"
password = "your_password"
charset = "utf8mb4"
```

支持三种输出方式:
1. 输出到标准输出
2. 输出到剪贴板
3. 输出到文件

支持生成标签:
1. json 标签
2. toml 标签
3. gorm 标签
4. form 标签
5. yaml 标签
6. comment 标签
7. 以及自定义标签


使用时, 必须在命令中指定数据库名和表名:

使用示例:

输出到命令行
```bash
dbtogo --db 数据库名 --table 表名
```

输出到系统剪贴板
```bash
dbtogo --db 数据库名 --table 表名 --clip
```

输出到文件
```bash
dbtogo --db 数据库名 --table 表名 --file struct.txt
```

生成各种 tag 标签, 可以写在一起, 如下, 将生成 json gorm toml 的 tag
```bash
dbtogo --db 数据库名 --table 表名 -jgt
```

指定配置文件位置
```bash
dbtogo --db 数据库名 --table 表名 --config ~/dbtogo.toml
```

更多信息, 使用 -h 选项查看

## 备注

- 实现过程参考了煎鱼大佬的代码 https://github.com/go-programming-tour-book/tour
- 输出到剪贴板并不总是可以成功的, 在 linux 下, 你可能需要安装 xclip 或 xsel 命令, 具体请看这里 https://github.com/atotto/clipboard

## License

MIT