package entity

type Config struct {
	Namespace string
	Redis     RedisConfig
	SMTP      SMTPConfig
}

type SMTPConfig struct {
	Host     string `env:"HOST" envDefault:"smtp.hostinger.com"`
	Port     string `env:"PORT" envDefault:"465"`
	Password string `env:"PASSWORD" envDefault:"AGA_hrms1"`
}

type RedisConfig struct {
	Host string
	Port string
}
