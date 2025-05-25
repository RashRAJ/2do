package models

type DBConn struct {
	DbConfig DbConfig `yaml:"dbconfig"`
}

type DbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
}
