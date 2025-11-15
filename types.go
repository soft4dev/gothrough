package main

type Backend struct {
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
}

type Server struct {
	Listen []int `yaml:"listen"`
}

type Config struct {
	Server   Server    `yaml:"server"`
	Backends []Backend `yaml:"backends"`
}
