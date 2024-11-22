package config

type DatabaseConfig struct {
	Host     string `json:"host" env:"WEATHERDB_HOST" env-default:"localhost"`
	Port     string `json:"port" env:"WEATHERDB_PORT" env-default:"5432"`
	Username string `json:"username" env:"WEATHERDB_USERNAME" env-default:"weather"`
	Password string `json:"password" env:"WEATHERDB_PASSWORD" env-default:"weather"`
	Name     string `json:"name" env:"WEATHERDB_NAME" env-default:"weather"`
}

type APIConfig struct {
	WritePassword   string `json:"password" env:"WEATHERDB_API_PASS" env-default:"weather"`
	DefaultLocation string `json:"location" env:"WEATHERDB_API_LOCATION" env-default:"unset"`
}

type Config struct {
	Database      DatabaseConfig `json:"database"`
	API           APIConfig      `json:"API"`
	ListenAddress string         `json:"listen" env:"WEATHERDB_LISTEN" env-default:":3000"`
}
