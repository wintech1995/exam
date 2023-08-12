/*=======================================================
   __                 ____  _
  |   \  ___  __ __  |  _ || |  _ _   ___
  | |) |/ -_) \ V /  |  __|| | | | | / -_)
  |___/ \___|  \_/   |_|   |_| \_,_| \___| o o o

=======================================================*/

package main

import (
	"fmt"
	AuthController "wintech1995/my-app/controller/auth"
	BillController "wintech1995/my-app/controller/bill"
	ChatController "wintech1995/my-app/controller/chat"
	ProductController "wintech1995/my-app/controller/product"
	UserController "wintech1995/my-app/controller/user"

	"wintech1995/my-app/middleware"
	orm "wintech1995/my-app/orm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	orm.InitDB()

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/register", AuthController.Register)
	r.POST("/login", AuthController.Login)

	// r.GET("/product", ProductController.ProductList)
	// r.GET("/product/:id", ProductController.ProductData)

	authorized := r.Group("/users", middleware.JWTAuthen())
	authorized.GET("/readall", UserController.ReadAll)
	authorized.GET("/profile", UserController.Profile)
	authorized.POST("/uploadavatar", UserController.UploadAvatar)

	auth_product := r.Group("/product", middleware.JWTAuthen())
	auth_product.GET("", ProductController.ProductList)
	auth_product.POST("/edit", ProductController.UpdateStatus)
	auth_product.GET("/:id", ProductController.ProductData)

	auth_bill := r.Group("/bill", middleware.JWTAuthen())
	auth_bill.GET("/:id", BillController.ProductData)

	auth_chat := r.Group("/chat", middleware.JWTAuthen())
	auth_chat.POST("/send", ChatController.SendChat)

	r.Run("localhost:8080")
}
