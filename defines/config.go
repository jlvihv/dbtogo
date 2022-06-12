package defines

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	IP       string `toml:"ip"`
	Port     string `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Charset  string `toml:"charset"`
}
