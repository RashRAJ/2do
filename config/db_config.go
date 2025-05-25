package config

type Config struct {
	//Server   ServerConfig `yaml:"server"`
	DbConfig DbConfig `yaml:"dbConfig"`
}

type DbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
}

//type ServerConfig struct {
//	Port int `yaml:"port"`
//}
