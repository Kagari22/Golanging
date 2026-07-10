package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 读取 json 文件数据并反序列化输出
func main() {
	data, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("姓名: %s, 年龄: %d\n", user.Name, user.Age)
}	
