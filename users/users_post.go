package users

import (
	"encoding/json"
	"fmt"
	"sub-ui/random"
	"sub-ui/setup"
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
			ConfigData.Inbounds[i].Port = fmt.Sprintf("%d", tagData.Port)

			if !ConfigData.Inbounds[i].FixedSecurity {
				ConfigData.Inbounds[i].Security = tagData.Security
			}

			if tagData.Alpn != "" && ConfigData.Inbounds[i].Security != "reality" {
				if ConfigData.Inbounds[i].Tls == nil {
					ConfigData.Inbounds[i].Tls = new(Tls)
				}
				ConfigData.Inbounds[i].Tls.Alpn = tagData.Alpn
			}

		}
		toggleContent = "set"
	} else {
		inboundsLen := len(ConfigData.Inbounds)
		i := tagData.Index

		if inboundsLen < i || 0 > i {

			return response, toggleContent, "网络错误!"
		}

		if ConfigData.Inbounds[i].Tag != tagData.Tag {
			return response, toggleContent, "网络错误!"
		}

		ConfigData.Inbounds[i].Addr = tagData.Addr
		ConfigData.Inbounds[i].Port = fmt.Sprintf("%d", tagData.Port)
		if !ConfigData.Inbounds[i].FixedSecurity {
			ConfigData.Inbounds[i].Security = tagData.Security
		}

		if tagData.Alpn != "" && ConfigData.Inbounds[i].Security != "reality" {
			if ConfigData.Inbounds[i].Tls == nil {
				ConfigData.Inbounds[i].Tls = new(Tls)
			}
			ConfigData.Inbounds[i].Tls.Alpn = tagData.Alpn
		}

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

func (re RenewUsers) SetUsersUrl() error {

	inboundsLen := len(ConfigData.Inbounds)

	for i := range re.Users {
		x := re.Users[i].X
		y := re.Users[i].Y

		if inboundsLen < x || 0 > x {
			break
		}

		UsersLen := len(ConfigData.Inbounds[x].Users)

		if UsersLen < y || 0 > y {
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

		if !setup.ConfigData.Static.Enabled || !ConfigData.Inbounds[x].Users[y].Static {
			ConfigData.Inbounds[x].Users[y].UserPath = path
		}

	}

	err := ConfigData.SavedConfig()
	if err != nil {
		return err
	}

	return nil
}

func (re RenewUsers) SetStaticUsers() {

	var newConsts []setup.Consts
	var users setup.ConstUser

	inboundsLen := len(ConfigData.Inbounds)

	for i := range re.Users {
		x := re.Users[i].X
		y := re.Users[i].Y

		if inboundsLen < x || 0 > x {
			break
		}

		UsersLen := len(ConfigData.Inbounds[x].Users)

		if UsersLen < y || 0 > y {
			break
		}

		if ConfigData.Inbounds[x].Users[y].Name != re.Users[i].Name {
			break
		}

		tag := ConfigData.Inbounds[x].Tag

		isNew := true

		for j := range newConsts {
			if tag == newConsts[j].Tag {
				newConsts[j].Users = append(newConsts[j].Users, setup.ConstUser{
					Name: ConfigData.Inbounds[x].Users[y].Name,
					Path: ConfigData.Inbounds[x].Users[y].UserPath,
				})
				isNew = false
				break
			}
		}

		if !isNew {
			continue
		}

		users = setup.ConstUser{
			Name: ConfigData.Inbounds[x].Users[y].Name,
			Path: ConfigData.Inbounds[x].Users[y].UserPath,
		}

		newConsts = append(newConsts, setup.Consts{
			Tag:   tag,
			Users: []setup.ConstUser{users},
		})

	}

	//fmt.Println(newConsts)

	setup.ConfigData.Static.ConstList = newConsts
	setup.SavedConfig()

}

func (bac BackupInfo) AddUsers() {

	var newExcludes []setup.Exclude

	inboundsLen := len(ConfigData.Inbounds)

	for i := range bac.Users {
		x := bac.Users[i].X
		y := bac.Users[i].Y

		if inboundsLen < x || 0 > x {
			break
		}

		UsersLen := len(ConfigData.Inbounds[x].Users)

		if UsersLen < y || 0 > y {
			break
		}

		if ConfigData.Inbounds[x].Users[y].Name != bac.Users[i].Name {
			break
		}

		tag := ConfigData.Inbounds[x].Tag

		isNew := true

		for j := range newExcludes {
			if tag == newExcludes[j].Tag {
				newExcludes[j].Users = append(newExcludes[j].Users, bac.Users[i].Name)
				isNew = false
				break
			}
		}

		if !isNew {
			continue
		}

		newExcludes = append(newExcludes, setup.Exclude{
			Tag:   tag,
			Users: []string{bac.Users[i].Name},
		})
	}

	setup.ConfigData.Backup.Excludes = newExcludes
	setup.SavedConfig()
}
