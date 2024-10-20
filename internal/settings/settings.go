package settings

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Mode    string `mapstructure:"mode"`
	Version string `mapstructure:"version"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`

	*SnowFlakeConfig `mapstructure:"snowflake"`
	*AuthConfig      `mapstructure:"auth"`
	*LogConfig       `mapstructure:"log"`
	*MySQLConfig     `mapstructure:"mysql"`
	*RedisConfig     `mapstructure:"redis"`
	*RabbitMQConfig  `mapstructure:"rabbitmq"`
	*EsConfig        `mapstructure:"elasticsearch"`
	*RpcServerConfig `mapstructure:"rpc_server"`
	*JaegerConfig    `mapstructure:"jaeger"`
}

type SnowFlakeConfig struct {
	StartTime   string                          `mapstructure:"start_time"`
	TableConfig map[string]TableSnowFlakeConfig `mapstructure:"table_config"`
}

type TableSnowFlakeConfig struct {
	DataCenterID int `mapstructure:"datacenter_id"`
	WorkerID     int `mapstructure:"worker_id"`
}

type AuthConfig struct {
	JwtExpire int `mapstructure:"jwt_expire"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Vhost    string `mapstructure:"vhost"`
}

type EsConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type RpcServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JaegerConfig struct {
	Host              string `mapstructure:"host"`
	HttpCollectorPort int    `mapstructure:"http_collector_port"`
}

func Init() (err error) {
	//viper.SetConfigName("config") //指定配置文件名称（不需要制定配置文件的扩展名）
	//viper.SetConfigType("yaml")   //指定配置文件类型（专用于从远程配置信息指定配置文件类型）
	//viper.AddConfigPath(".")      //指定查找配置文件的路径（这里使用相对路径）
	viper.SetConfigFile("./conf/config.yaml")
	err = viper.ReadInConfig() //读取配置信息
	if err != nil {
		//读取失败
		zap.L().Error("viper.ReadInConfig() failed", zap.Error(err))
		return err
	}
	viper.WatchConfig() //监听配置文件变化
	viper.OnConfigChange(func(e fsnotify.Event) {
		zap.L().Info("配置文件修改了...")
	})
	if err := viper.Unmarshal(&Conf); err != nil {
		zap.L().Error("viper.Unmarshal failed", zap.Error(err))
		return err
	}
	return
}
