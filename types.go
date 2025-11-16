package main

type Backend struct {
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
}

type Server struct {
	Listen   int       `yaml:"listen"`
	Backends []Backend `yaml:"backends"`
}

type L4 struct {
	Servers []Server `yaml:"servers"`
}

type Config struct {
	L4 L4 `yaml:"l4"`
}
