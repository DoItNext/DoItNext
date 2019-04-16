package config

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-mods/zerolog-rotate/log"
	"github.com/spf13/viper"
)

// the ServerConfig variable stores the
// server configurations
var Configuration configuration

// Main configuration
type configuration struct {
	Server   server   `mapstructure:"server"`
	Database database `mapstructure:"database"`
}

// Server configuration
type server struct {
	Port         uint `mapstructure:"port";valid:"port,required"`
	Debug        bool `mapstructure:"debug";valid:"required"`
	ReadTimeout  int  `mapstructure:"read_timeout";valid:"required"`
	WriteTimeout int  `mapstructure:"write_timeout";valid:"required"`
}

// Database configuration
type database struct {
	Type     string `mapstructure:"type";valid:"in(mysql)"`
	Host     string `mapstructure:"host";valid:"required"`
	Port     uint   `mapstructure:"port";valid:"port,required"`
	Name     string `mapstructure:"name";valid:"required"`
	User     string `mapstructure:"user";valid:"required"`
	Password string `mapstructure:"password";valid:"required"`
	Charset  string `mapstructure:"charset";valid:"in(utf8)"`
}

// Load loads configuration from the given list of paths and populates it into the ServerConfig variable.
// The configuration file(s) should be named as server.yaml.
// Environment variables with the prefix "DOITNEXT_" in their names are also read automatically.
func Load(configPaths ...string) error {
	v := viper.New()
	// look for a config file named server.yaml
	v.SetConfigName("server")
	v.SetConfigType("yaml")
	// look for env variables that start with "DOITNEXT_".
	v.SetEnvPrefix("doitnext")
	v.AutomaticEnv()
	// Add paths to look for
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}
	// Find and read the config file
	if err := v.ReadInConfig(); err != nil {
		log.Error(err, "Failed to read the configuration file")
		return err
	}
	// Fill ServerConfig variable with data from server.yaml
	if err := v.Unmarshal(&Configuration); err != nil {
		log.Error(err, "Failed to unmarshal to ServerConfig")
		return err
	}
	//
	log.Debug("Server configuration loaded")
	// Validate imported data
	return Configuration.validate()
}

// Validate data from the config file
func (config configuration) validate() error {
	_, err := govalidator.ValidateStruct(&config)
	return err
}
