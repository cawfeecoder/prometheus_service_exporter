package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Supervisor string
	Service_Metric_Series_Name string
	Pid_Metric_Series_Name string
	Collect_Pids bool
	Targets []TargetGroup
}

type TargetGroup struct {
	Name string
	Daemon_Whitelist []string
	Pid_Whitelist []PidTarget
}

type PidTarget struct {
	Name string
	Service string
}

func LoadConfig(path string) (config Config, err error){
    viper.SetConfigFile(path)
    viper.SetConfigType("yaml")
    err = viper.ReadInConfig()
    if err != nil {
    	return
	}
    config = Config{}
    viper.Unmarshal(&config)
    return
}
