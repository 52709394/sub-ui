package read

import (
	"encoding/json"
	"fmt"
	"os"
)

func GetJsonData(name string, data any) error {
	file, err := os.ReadFile(name)
	if err != nil {
		fmt.Println("文件:", name)
		fmt.Println("读取错误:", err)
		return err
	}

	err = json.Unmarshal(file, data)
	if err != nil {
		fmt.Println("文件:", name)
		fmt.Println("JSON解析错误:", err)
		return err
	}
	return nil
}

func CheckExistence(path string) string {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "nil"
	}

	if err == nil && info.IsDir() {
		return "folder"
	}

	if err == nil && !info.IsDir() {
		return "file"
	}

	return "nil"
}
