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

// 把数据写入 json 文件中
func main() {
	user := User{
		Name: "张三",
		Age: 25,
	}

	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = os.WriteFile("data.json", data, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("写入成功")
}
