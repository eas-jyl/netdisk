package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID   uint
	Name string
	Age  int
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	fmt.Println("MySQL 连接成功")

	// 更改结果
	var update_user User
	db.First(&update_user)

	update_user.Age = 20
	update_user.Name = "baijiyi"

	db.Save(&update_user)

	// 查询结果
	var query_user User
	db.First(&query_user)
	fmt.Printf("查询结果： ID = %d , Name = %s , Age = %d\n", query_user.ID, query_user.Name, query_user.Age)

}
