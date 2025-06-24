package route

import (
	ctrl "github.com/DevisArya/BE-challenge-syn/UserService/internal/delivery/http"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *gin.Engine
	UserController ctrl.UserController
	AuthMiddleware fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.POST("/api/user", c.UserController.Create)
	c.App.POST("/api/user/login", c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute() {

	//use auth
	// c.App.Use(c.AuthMiddleware)
	api := c.App.Group("/api/user")

	api.GET("s", c.UserController.FindAll)
	// api.GET("/:id", c.UserController.FindById)
	// api.DELETE("", c.UserController.Delete)
	// api.PATCH("/email", c.UserController.UpdateEmail)
	// api.PATCH("/name", c.UserController.UpdateNameProfile)
	// api.PATCH("/password", c.UserController.UpdatePassword)
	// api.DELETE("/logout", c.UserController.Logout)
}
