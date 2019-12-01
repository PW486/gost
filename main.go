package main

import (
	"net/http"

	"github.com/PW486/gost/db"
	"github.com/PW486/gost/dto"
	"github.com/PW486/gost/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func getHandler(c *gin.Context) {
	var accounts []model.Account
	db.Service().Find(&accounts)

	c.JSON(200, gin.H{"data": accounts})
}

func postHandler(c *gin.Context) {
	var createAccountDTO dto.CreateAccountDTO
	if err := c.ShouldBindJSON(&createAccountDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newAccount model.Account
	newAccount.ID, _ = uuid.NewUUID()
	newAccount.Email = createAccountDTO.Email
	newAccount.Name = createAccountDTO.Name
	newAccount.Password, _ = bcrypt.GenerateFromPassword([]byte(createAccountDTO.Password), 10)

	db.Service().Create(&newAccount)

	c.JSON(201, gin.H{"data": newAccount})
}

func logInHandler(c *gin.Context) {
	var logInDTO dto.LogInDTO
	if err := c.ShouldBindJSON(&logInDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var account model.Account
	db.Service().Where("Email = ?", logInDTO.Email).First(&account)

	if err := bcrypt.CompareHashAndPassword(account.Password, []byte(logInDTO.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mySigningKey := []byte("AllYourBase")

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)

	c.JSON(200, gin.H{"token": ss})
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", getHandler)
	r.POST("/", postHandler)
	r.POST("/login", logInHandler)

	return r
}

func main() {
	db.Open()
	db.Migration()

	r := setupRouter()
	r.Run(":8080")
}
