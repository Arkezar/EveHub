package main

import (
	"github.com/spf13/viper"
)

var Configuration = initializeConfiguration()

func initializeConfiguration() *viper.Viper {
	viper.SetConfigName("evehub")
	viper.AddConfigPath("/etc/evehub")
	viper.AddConfigPath("$HOME/.evehub")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return viper.GetViper()
}
