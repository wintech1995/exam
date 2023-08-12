package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"wintech1995/my-app/orm"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type RegisterBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
}

type LoginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var hmacSampleSecret []byte

func Register(c *gin.Context) {
	var json RegisterBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userExist orm.User
	orm.Db.Where("username = ?", json.Username).First(&userExist)
	if userExist.ID > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User exists"})
		return
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	user := orm.User{Username: json.Username, Password: string(encryptedPassword), Fullname: json.Fullname}
	orm.Db.Create(&user)
	if user.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "User create success",
			"userId":  user.ID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User create faild"})
	}
}

func Login(c *gin.Context) {
	var json LoginBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userExist orm.User
	orm.Db.Where("username = ?", json.Username).First(&userExist)
	if userExist.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User does not exists"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(json.Password))
	if err == nil {
		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": userExist.ID,
			"exp":    time.Now().Add(24 * time.Hour).Unix(),
		})
		tokenString, err := token.SignedString(hmacSampleSecret)
		fmt.Println(tokenString, err)

		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Login success",
			"token":   tokenString,
			"data":    userExist,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Login faild"})
	}
}
