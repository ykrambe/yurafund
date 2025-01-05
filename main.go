package main

import (
	"fmt"
	"log"
	"yurafund/auth"
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
	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)

	// userService.SaveAvatar(13, "images/avatar-13.png")
	token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNn0.nLdnBx_wPlE7Tb2mEUIqe_ePOEdUL-pXNoNS3340PPQ")
	if err != nil {
		fmt.Println("Error", err)
	}

	if token.Valid {
		fmt.Println("Token Valid", token)
	} else {
		fmt.Println("Token Invalid", token)
	}

	router := gin.Default()
	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)

	router.Run()
}
