package users

import (
	"encoding/json"
	"fmt"
	"os"
	"sub-ui/setup"
)

func setOldUrlPath(newUsers *[]User, oldUsers []User) {

	for i := range *newUsers {
		for j := range oldUsers {
			if (*newUsers)[i].Name == oldUsers[j].Name {
				(*newUsers)[i].UserPath = oldUsers[j].UserPath
			}
		}
	}
}

func (config *Config) SetOldData() {

	for i := range config.Inbounds {
		for j := range ConfigData.Inbounds {
			if config.Inbounds[i].Tag == ConfigData.Inbounds[j].Tag {
				config.Inbounds[i].Addr = ConfigData.Inbounds[j].Addr
				config.Inbounds[i].Port = ConfigData.Inbounds[j].Port

				if !config.Inbounds[i].FixedSecurity {
					config.Inbounds[i].Security = ConfigData.Inbounds[j].Security
				}

				config.Inbounds[i].Alpn = ConfigData.Inbounds[j].Alpn

				setOldUrlPath(&config.Inbounds[i].Users, ConfigData.Inbounds[j].Users)

			}
		}
	}
}

func (config Config) SavedConfig() error {
	nowData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("文件:", setup.ConfigData.Users.Config)
		fmt.Println("JSON格式化错误:", err)
		return err
	}

	err = os.WriteFile(setup.ConfigData.Users.Config, nowData, 0644)
	if err != nil {
		fmt.Println("文件:", setup.ConfigData.Users.Config)
		fmt.Println("文件写入错误:", err)
		return err
	}
	ConfigData = config
	return nil
}
