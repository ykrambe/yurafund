package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"yurafund/auth"
	"yurafund/campaign"
	"yurafund/handler"
	"yurafund/helper"
	"yurafund/user"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//koneksi ke database mysql with gorm
	// dsn := "freedb_yuraroot:N5C@FK6PDChn&Yh@tcp(sql.freedb.tech:3306)/freedb_yuradb?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := ""
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connection Opened to Database")

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)

	campaignService := campaign.NewService(campaignRepository)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	campaign, _ := campaignService.GetCampaigns(1)
	for _, a := range campaign {
		fmt.Println(a.Name)
	}

	// userService.SaveAvatar(13, "images/avatar-13.png")
	token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.v64cDEe2k9U0bbF0S1NmUhJhDR3CBZoqreD_TcWifK0")
	if err != nil {
		fmt.Println("Error", err)
	}

	if token.Valid {
		fmt.Println("Token Valid", token)
	} else {
		fmt.Println("Token Invalid", token)
	}

	router := gin.Default()
	router.Static("/images", "./images")

	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	api.GET("campaigns", campaignHandler.GetCampaigns)

	router.Run()
}

func authMiddleware(authService auth.Service, userservice user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response) //AbortWithStatusJSON means if unauthorized, process is stopped in this process, not coninuous to next process
			return
		}

		// get token only
		tokenString := ""
		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) == 2 {
			tokenString = splitToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userservice.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
