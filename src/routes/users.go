package routes

import (
	"net/http"

	"github.com/MaxiOtero6/go-auth-rest/controller"
	"github.com/MaxiOtero6/go-auth-rest/service"
	"github.com/gin-gonic/gin"
)

func UsersRoutes(router *gin.Engine) {
	{
		users_routes := router.Group("/users")
		users_routes.GET("/", getUsers)
		users_routes.GET("/:username", getUser)
	}
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, service.GetAllUsers())
}

func getUser(c *gin.Context) {
	username := c.Param("username")

	controller := controller.UserController{}

	controller.ValidateString(username, "username")

	user, err := service.GetUser(username)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}
