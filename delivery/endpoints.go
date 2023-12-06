package delivery

import "github.com/gin-gonic/gin"

func InitEndPoints(router *gin.Engine, service *Service) {
	endpoints := router.Group("/tasks")
	endpoints.POST("/add", service.CreateTask)
	endpoints.GET("/all", service.GetTasks)
	endpoints.PUT("/:id", service.UpdateTask)
	endpoints.GET("/:id", service.GetTaskByID)
	endpoints.DELETE("/:id", service.DeleteTask)
	endpoints.GET("/report", service.GetTaskReport)
}
