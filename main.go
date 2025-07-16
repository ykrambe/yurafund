package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func validationEnv(vars ...string) {
	for _, v := range vars {
		if os.Getenv(v) == "" {
			log.Fatalf("Error: Missing required environment variable: %s", v)
		}
	}
}

func main() {
	//koneksi ke database mysql with gorm
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file not found")
	}

	validationEnv("PORT", "MIDTRANS_SERVER_KEY", "MIDTRANS_CLIENT_KEY", "JWT_SECRET", "SESSION_NAME", "DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME")

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
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

	// handler web
	userWebHandler := webHandler.NewUserHandler(userService)
	campaignWebHanlder := webHandler.NewCampaignHandler(campaignService, userService)
	transactionWebHandler := webHandler.NewTransactionHandler(transactionService)
	sessionWebHandler := webHandler.NewSessionHandler(userService)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Atur URL frontend kamu
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	cookieStore := cookie.NewStore([]byte(auth.SECRET_KEY))
	router.Use(sessions.Sessions(os.Getenv("SESSION_NAME"), cookieStore))

	router.HTMLRender = loadTemplates("./web/templates")

	router.Static("/user_images", "./user_images")
	router.Static("/campaign_images", "./campaign_images")
	router.Static("/css", "./web/assets/css")
	router.Static("/js", "./web/assets/js")
	router.Static("/webfonts", "./web/assets/webfonts")
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
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns-limit", campaignHandler.GetCampaignsWithLimit)
	api.GET("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.Uploadimage)

	// endpoint transaction
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	//admin cms dashboards
	router.GET("/users", userWebHandler.Index)
	router.GET("/users/new", userWebHandler.New)
	router.POST("/users", userWebHandler.Create)
	router.GET("/users/edit/:id", userWebHandler.Edit)
	router.POST("/users/update/:id", userWebHandler.Update)
	router.GET("/users/avatar/:id", userWebHandler.NewAvatar)
	router.POST("/users/avatar/:id", userWebHandler.CreateAvatar)

	router.GET("/campaigns", authAdminMiddleware(), campaignWebHanlder.Index)
	router.GET("/campaigns/new", authAdminMiddleware(), campaignWebHanlder.New)
	router.POST("/campaigns", authAdminMiddleware(), campaignWebHanlder.Create)
	router.GET("/campaigns/image/:id", authAdminMiddleware(), campaignWebHanlder.NewImage)
	router.POST("/campaigns/image/:id", authAdminMiddleware(), campaignWebHanlder.CreateImage)
	router.GET("/campaigns/edit/:id", authAdminMiddleware(), campaignWebHanlder.Edit)
	router.POST("/campaigns/update/:id", authAdminMiddleware(), campaignWebHanlder.Update)
	router.GET("/campaigns/show/:id", authAdminMiddleware(), campaignWebHanlder.Show)
	router.GET("/transactions", authAdminMiddleware(), transactionWebHandler.Index)

	router.GET("/login", sessionWebHandler.New)
	router.POST("/session", sessionWebHandler.Create)
	router.GET("/logout", sessionWebHandler.Destroy)

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

func authAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		userIDSession := session.Get("userID")
		fmt.Println("userIDSession", userIDSession)
		if userIDSession == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
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
