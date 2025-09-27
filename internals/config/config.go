package config

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
)

type Config struct {
	Connection string
	IpAddress  string
	Port       string
	Secret     string
	Directory  string
	Silent     bool
}

func (c *Config) IsValid() bool {
	var err error

	validIp, err := regexp.Match(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`, []byte(c.IpAddress))
	if err != nil {
		return false
	}

	validPort, err := regexp.Match(`\d{2,6}`, []byte(c.Port))
	if err != nil {
		return false
	}

	validDir := len(c.Directory) != 0

	return validIp && validPort && validDir
}

func Parse() Config {
	ipAddress := flag.String("ip", "", "ip address")
	port := flag.String("p", "5090", "port")
	secret := flag.String("s", "Cd25+_0v{.18", "secret")
	connection := flag.String("c", "ws", "connection type")
	directory := flag.String("dir", "", "directory")
	silent := flag.Bool("q", false, "silent logging")
	flag.Parse()

	config := Config{}
	config.IpAddress = *ipAddress
	config.Port = *port
	config.Secret = *secret
	config.Connection = *connection
	config.Silent = *silent
	config.Directory = parseDirectory(*directory)

	return config
}

func parseDirectory(dirPath string) string {
	if len(dirPath) == 0 {
		ex, err := os.Executable()
		if err != nil {
			return ""
		}

		dirPath = filepath.Dir(ex)
	}

	return dirPath
}
