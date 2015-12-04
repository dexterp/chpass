package config

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
)

// {
//   User: username
//   Password: file
//   Hosts: [
//     { "Name": HOST, "Port": PORT },
//     { "Name": HOST, "Port": PORT }
//   ]
// }
type Host struct {
	Name string
	Port int
}

type Config struct {
	User string
	Newpassword string
	Oldpassword string
	Hosts []Host
}

func (config *Config) LoadConfig(path string) error {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Can not open %s: %v", path, err)
	}
	err = json.Unmarshal(jsonData,config)
	if err != nil {
		return fmt.Errorf("Can not load config file %s: %v", path, err)
	}
	return nil
}
