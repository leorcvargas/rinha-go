package config

import "github.com/leorcvargas/rinha-2023-q3/pkg/env"

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	Cache struct {
		Host string
		Port string
	}

	Server struct {
		Port     string
		UseSonic bool
		Prefork  bool
	}

	Profiling struct {
		Enabled bool
		CPU     string
		Mem     string
	}
}

func NewConfig() *Config {
	return &Config{
		Database: struct {
			Host     string
			Port     string
			User     string
			Password string
			Name     string
		}{
			Host:     env.GetEnvOrDie("DB_HOST"),
			Port:     env.GetEnvOrDie("DB_PORT"),
			User:     env.GetEnvOrDie("DB_USER"),
			Password: env.GetEnvOrDie("DB_PASSWORD"),
			Name:     env.GetEnvOrDie("DB_NAME"),
		},

		Cache: struct {
			Host string
			Port string
		}{
			Host: env.GetEnvOrDie("CACHE_HOST"),
			Port: env.GetEnvOrDie("CACHE_PORT"),
		},

		Server: struct {
			Port     string
			UseSonic bool
			Prefork  bool
		}{
			Port:     env.GetEnvOrDie("SERVER_PORT"),
			UseSonic: env.GetEnvOrDie("ENABLE_SONIC_JSON") == "1",
			Prefork:  env.GetEnvOrDie("ENABLE_PREFORK") == "1",
		},

		Profiling: struct {
			Enabled bool
			CPU     string
			Mem     string
		}{
			Enabled: env.GetEnvOrDie("ENABLE_PROFILING") == "1",
			CPU:     env.GetEnvOrDie("CPU_PROFILE"),
			Mem:     env.GetEnvOrDie("MEM_PROFILE"),
		},
	}
}
