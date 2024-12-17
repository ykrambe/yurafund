package main

import (
	"fmt"
	"log"
	"yurafund/handler"
	"yurafund/user"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//koneksi ke database mysql with gorm
	dsn := "freedb_yuraroot:N5C@FK6PDChn&Yh@tcp(sql.freedb.tech:3306)/freedb_yuradb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connection Opened to Database")

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	// userInput := user.RegisterUserInput{}
	// userInput.Name = "Test simpan"
	// userInput.Email = "test@mail.com"
	// userInput.Occupation = "Programmer"
	// userInput.Password = "password"
	// userService.RegisterUser(userInput)

	// input := user.LoginInput{
	// 	Email:    "test@mail.com",
	// 	Password: "passwordqwe",
	// }
	// user, err := userService.Login(input)
	// if err != nil {
	// 	fmt.Println("terjadi kesalahan")
	// 	fmt.Println(err.Error())

	// } else {
	// 	fmt.Println(user)
	// }

	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkkers", userHandler.CheckEmailAvailability)

	router.Run()
}
