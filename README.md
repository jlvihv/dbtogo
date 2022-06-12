# dbtogo

这是一个将数据库中存在的表转换为 go 结构体的命令行工具

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

使用前, 必须在 config.toml 配置文件中配置数据库连接信息:

配置示例: config.toml

```toml
[db]
ip = "127.0.0.1"
port = "3306"
username = "username"
password = "password"
charset = "utf8mb4"
```

使用时, 必须在命令中指定数据库名和表名:

使用示例:

```bash
dbtogo --db db_name --table table_name
```