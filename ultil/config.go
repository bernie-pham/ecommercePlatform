package ultils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPServerAddr  string        `mapstructure:"HTTP_SERVER_ADDR"`
	Environment     string        `mapstructure:"ENVIRONMENT"`
	DBDriverName    string        `mapstructure:"DB_DRIVER_NAME"`
	DBSource        string        `mapstructure:"DB_DATA_SOURCE"`
	SymKey          string        `mapstructure:"SYMMENTRIC_KEY"`
	AccessTimeout   time.Duration `mapstructure:"ACCESS_TIMEOUT"`
	RefreshTimeout  time.Duration `mapstructure:"REFRESH_TIMEOUT"`
	VerifyTimeout   time.Duration `mapstructure:"VERIFY_TINMEOUT"`
	RedisAddr       string        `mapstructure:"REDIS_ADDRESS"`
	MailSender      string        `mapstructure:"EMAIL_SENDER"`
	MailTempltePath string        `mapstructure:"EMAIL_TEMPLATE_PATH"`
	GRPCServerAddr  string        `mapstructure:"GRPC_SERVER_ADDR"`
}

func LoadConfig() (Config, error) {
	var config Config
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}
	viper.Unmarshal(&config)

	return config, nil
}
