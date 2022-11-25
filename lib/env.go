package lib

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	DBUsername string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASS"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
	DBType     string `mapstructure:"DB_TYPE"`
	StorageBucketName  string `mapstructure:"STORAGE_BUCKET_NAME"`

}

var globalEnv = Env{}

func GetEnv() Env {
	return globalEnv
}

func NewEnv() *Env {
	viper.SetConfigFile(".env");

	if err:=viper.ReadInConfig(); err!=nil{
		log.Println("Cannot Read Configuration",err)
	}
	if err:=viper.Unmarshal(&globalEnv); err!=nil{
		log.Println("Environmen Cannot Be Loaded",err)
	}
	
	return &globalEnv
}