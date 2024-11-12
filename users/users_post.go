package users

import (
	"encoding/json"
	"fmt"
	"sub-ui/random"
)

func SetTagData(body []byte) (map[string]string, string, string) {

	var response map[string]string
	var toggleContent string
	var err error

	var tagData TagData

	err = json.Unmarshal(body, &tagData)
	if err != nil {
		return response, toggleContent, "无法解析数据!"
	}

	var tagSrtA, tagSrtB, tagInfoSrt string

	response = map[string]string{"status": "success", "message": "successfully",
		"tagSrtA": tagSrtA, "tagSrtB": tagSrtB, "tagInfoSrt": tagInfoSrt}

	if tagData.Index == -1 {
		for i := range ConfigData.Inbounds {
			ConfigData.Inbounds[i].Addr = tagData.Addr
			ConfigData.Inbounds[i].Port = tagData.Port

			if !ConfigData.Inbounds[i].FixedSecurity {
				ConfigData.Inbounds[i].Security = tagData.Security
			}

			ConfigData.Inbounds[i].Alpn = tagData.Alpn
		}
		toggleContent = "set"
	} else {
		inboundsLen := len(ConfigData.Inbounds)
		i := tagData.Index

		if inboundsLen < i && 1 > i {

			return response, toggleContent, "网络错误!"
		}

		if ConfigData.Inbounds[i].Tag != tagData.Tag {
			return response, toggleContent, "网络错误!"
		}

		ConfigData.Inbounds[i].Addr = tagData.Addr
		ConfigData.Inbounds[i].Port = tagData.Port
		if !ConfigData.Inbounds[i].FixedSecurity {
			ConfigData.Inbounds[i].Security = tagData.Security
		}
		ConfigData.Inbounds[i].Alpn = tagData.Alpn

		tagSrtA, tagSrtB, tagInfoSrt = TagHttpString(ConfigData.Inbounds[i])

	}

	err = ConfigData.SavedConfig()

	if err != nil {
		return response, toggleContent, "无法保存数据,请重试!"
	}

	response = map[string]string{"status": "success", "message": "successfully",
		"tagSrtA": tagSrtA, "tagSrtB": tagSrtB, "tagInfoSrt": tagInfoSrt}
	return response, toggleContent, ""
}

func (re ResetUrl) SetUserstUrl() error {

	inboundsLen := len(ConfigData.Inbounds)

	for i := range re.Users {
		x := re.Users[i].X
		y := re.Users[i].Y

		if inboundsLen < x && 1 > x {
			break
		}

		UsersLen := len(ConfigData.Inbounds[x].Users)

		if UsersLen < y && 1 > y {
			break
		}

		if ConfigData.Inbounds[x].Users[y].Name != re.Users[i].Name {
			break
		}

		path, err := random.GenerateStrings(16)
		if err != nil {
			fmt.Println("重置Url警告:随机Url路径错误!")
			fmt.Println("错误提示:", err)
			return err
		}
		ConfigData.Inbounds[x].Users[y].UserPath = path

	}

	err := ConfigData.SavedConfig()
	if err != nil {
		return err
	}

	return nil
}
