package config


type MongoConfig struct {
	Host string `json:"host"`
	Database string `json:"database"`
	Port string `json:"port"`
}
