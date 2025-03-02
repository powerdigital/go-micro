package config

import "fmt"

type GRPC struct {
	NetworkType string `default:"tcp"  envconfig:"GRPC_NETWORK_TYPE"`
	Port        uint16 `default:"8081" envconfig:"GRPC_PORT"`
}

func (c *Config) GRPCNetworkType() string {
	return c.GRPC.NetworkType
}

func (c *Config) GRPCAddress() string {
	return fmt.Sprintf(":%d", c.GRPC.Port)
}
