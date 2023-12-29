package router

import (
	"task-5-pbi-btpns-daffa_satria/controllers"
	"task-5-pbi-btpns-daffa_satria/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	public := r.Group("/api/v1")
	{
		public.GET("/users/login", controllers.Login)
		public.POST("/users/register", controllers.Register)
	}

	protected := r.Group("/api/v1")
	protected.Use(middlewares.Authentication())
	{
		protected.GET("/users/:id", controllers.GetUserById)
		protected.PUT("/users/:id", controllers.UpdateUser)
		protected.DELETE("/users/:id", controllers.DeleteUser)
		
		protected.GET("/photos", controllers.GetAllPhotos)
		protected.GET("/photos/:id", controllers.GetPhotoById)
		protected.POST("/photos", controllers.UploadPhoto)
		protected.PUT("/photos/:id", controllers.UpdatePhoto)
		protected.DELETE("/photos/:id", controllers.DeletePhoto)
	}

	return r
}