package config

import "time"

type Config struct {
	//Server   ServerConfig `yaml:"server"`
	DbConfig DbConfig `yaml:"dbConfig"`
}

type DbConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port,omitempty" json:"port,omitempty"`
	User     string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	DbName   string `yaml:"dbname" json:"dbname"`

	MaxOpenConns    int           `yaml:"maxOpenConns" json:"maxOpenConns"` // Maximum number of open connections
	MaxIdleConns    int           `yaml:"maxIdleConns" json:"maxIdleConns"` // Maximum number of idle connections
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime" json:"connMaxLifetime"`
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime" json:"connMaxIdleTime"`

	MaxRetries   int           `yaml:"maxRetries" json:"maxRetries"`
	InitialDelay time.Duration `yaml:"initialDelay" json:"initialDelay"`
	MaxDelay     time.Duration `yaml:"maxDelay" json:"maxDelay"`
	TotalTimeout time.Duration `yaml:"totalTimeout" json:"totalTimeout"`
}

//type ServerConfig struct {
//	Port int `yaml:"port"`
//}
