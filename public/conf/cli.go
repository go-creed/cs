package conf

import "github.com/micro/cli/v2"

var registryAddress string

func Init(c *cli.Context) {
	registryAddress = c.String("registry_address")

}

func RegistryAddress() string {
	return registryAddress
}
