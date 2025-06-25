package router

import (
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/delivery/rest/controller"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                    *fiber.App
	NotificationController controller.NotificationController
	AuthMiddleware         fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
}

func (c *RouteConfig) SetupAuthRoute() {

	//use auth
	c.App.Use(c.AuthMiddleware)
	api := c.App.Group("/api/notifications")

	api.Get("", c.NotificationController.GetNotifications)
}
