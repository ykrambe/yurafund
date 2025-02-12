package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"yurafund/auth"
	"yurafund/campaign"
	"yurafund/handler"
	"yurafund/helper"
	"yurafund/payment"
	"yurafund/transaction"
	"yurafund/user"

	webHandler "yurafund/web/handler"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
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

	//repository
	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	//service
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)
	paymentSercice := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentSercice)

	// handler
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	userWebhandler := webHandler.NewUserHandler()

	router := gin.Default()
	router.Use(cors.Default())

	router.HTMLRender = loadTemplates("./web/templates")

	router.Static("/images", "./images")
	// list endpoint
	api := router.Group("/api/v1")

	// endpoint user
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users", authMiddleware(authService, userService), userHandler.GetUsers)
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)

	// endpoint campaign
	api.GET("/campaigns", authMiddleware(authService, userService), campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.Uploadimage)

	// endpoint transaction
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	router.GET("/users", userWebhandler.Index)

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

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
