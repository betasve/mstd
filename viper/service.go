package viper

import (
	"github.com/spf13/viper"
)

var Client ViperService = Viper{}

type ViperService interface {
	AddConfigPath(in string)
	AutomaticEnv()
	ConfigFileUsed() string
	GetString(key string) string
	Set(key string, value interface{})
	SetConfigFile(in string)
	SetConfigName(in string)
	ReadInConfig() error
	WriteConfig() error
}

type Viper struct{}

func (v Viper) AddConfigPath(in string) {
	viper.AddConfigPath(in)
}

func (v Viper) AutomaticEnv() {
	viper.AutomaticEnv()
}

func (v Viper) ConfigFileUsed() string {
	return viper.ConfigFileUsed()
}

func (v Viper) GetString(key string) string {
	return viper.GetString(key)
}

func (v Viper) Set(key string, value interface{}) {
	viper.Set(key, value)
}

func (v Viper) SetConfigFile(in string) {
	viper.SetConfigFile(in)
}

func (v Viper) SetConfigName(in string) {
	viper.SetConfigName(in)
}

func (v Viper) ReadInConfig() error {
	return viper.ReadInConfig()
}

func (v Viper) WriteConfig() error {
	return viper.WriteConfig()
}
