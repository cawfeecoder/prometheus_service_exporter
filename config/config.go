package config

import (
	"errors"
	"github.com/spf13/viper"
	"regexp"
)

type Config struct {
	Supervisor string `yaml:"supervisor,omitempty"`
	ServiceWhitelist []*regexp.Regexp`yaml:"service_whitelist"`
	ServiceMetricFQDN string `yaml: "service_metric_fqdn"`
	PIDCollection    bool `yaml:"enable_pid"`
	PIDWhitelist []*regexp.Regexp `yaml:"pid_whitelist"`
	PIDMetricFQDN string `yaml: "pid_metric_fqdn"`
}

func LoadConfig(path string) (config Config, err error){
    viper.SetConfigFile(path)
    err = viper.ReadInConfig()
    if err != nil {
    	return
	}
    config = Config{}
    viper.Unmarshal(&config)
    config.ServiceMetricFQDN = viper.GetString("service_metric_fqdn")
    config.PIDCollection = viper.GetBool("enable_pid")
    config.PIDMetricFQDN = viper.GetString("service_metric_fqdn")
    if config.Supervisor == "" {
    	err = errors.New("a supervisor must be defined within your configuration")
    	return
	}
    for _, v := range viper.GetStringSlice("service_whitelist"){
    	r, err := regexp.Compile(v)
    	if err != nil {
    		err = errors.New("one or more whitelist values are not a valid string or regex")
    		break
		}
    	config.ServiceWhitelist = append(config.ServiceWhitelist, r)
	}
	for _, v := range viper.GetStringSlice("pid_whitelist"){
		r, err := regexp.Compile(v)
		if err != nil {
			err = errors.New("one or more whitelist values are not a valid string or regex")
			break
		}
		config.PIDWhitelist = append(config.PIDWhitelist, r)
	}
    if err != nil {
    	return
	}
    if len(config.ServiceWhitelist) == 0 {
    	err = errors.New("you must specify atleast one service to whitelist")
	}
    return
}
