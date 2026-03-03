package configs

type AppConfig struct {
	Port     string
	DBConfig DBConfig
}

type DBConfig struct {
	DBUrl string
}
