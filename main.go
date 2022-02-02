package main

import (
	"golang-api/config"
	"golang-api/controller"
	"golang-api/middleware"
	"golang-api/repository"
	"golang-api/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.SetupDatabaseConnection()
	bookRepository repository.BookRepository = repository.NewBookRepository(db)
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	jwtService     service.JWTService        = service.NewJwtService()
	bookService    service.BookService       = service.NewBookService(bookRepository)
	userService    service.UserService       = service.NewUserService(userRepository)
	authService    service.AuthService       = service.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
	bookController controller.BookController = controller.NewBookController(bookService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello",
		})
	})

	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}

	bookRoutes := r.Use(middleware.AuthorizeJWT(jwtService))
	bookRoutes.GET("/api/books", bookController.All)
	bookRoutes.POST("/api/books", bookController.Insert)
	bookRoutes.GET("/api/books/:id", bookController.FindByID)
	bookRoutes.PUT("/api/books/:id", bookController.Update)
	bookRoutes.DELETE("/api/books/:id", bookController.Delete)

	r.Run()
}
