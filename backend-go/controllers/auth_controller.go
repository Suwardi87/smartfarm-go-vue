package controllers

import (
	"net/http"

	"smartfarm-api/dto"
	"smartfarm-api/services"
	"smartfarm-api/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req dto.RegisterRequest

	// bind JSON ke DTO (mapping + validasi)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := services.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Generate JWT for Auto-Login
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Set Cookie
	c.SetCookie(
		"access_token",
		token,
		3600*24,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "registrasi berhasil",
		"data":    user,
	})
}

func Login(c *gin.Context) {

	var req dto.LoginRequest

	// bind JSON (mapping + validasi request)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// login service (cek user + password + buat JWT)
	token, err := services.LoginUser(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// set cookie (simpan JWT di browser)
	c.SetCookie(
		"access_token", // nama cookie
		token,          // nilai JWT
		3600*24,        // maxAge (detik) = 1 hari
		"/",            // path
		"",             // domain (kosong = current domain)
		false,          // secure (true kalau HTTPS)
		true,           // httpOnly (tidak bisa diakses JS)
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login berhasil",
	})
}

func Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := services.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := services.UpdateUserProfile(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := services.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user, "message": "Profile updated successfully"})
}

func Logout(c *gin.Context) {

	// hapus cookie dengan maxAge negatif
	c.SetCookie(
		"access_token", // nama cookie
		"",             // value dikosongkan
		-1,             // maxAge negatif = hapus
		"/",            // path
		"",             // domain
		false,          // secure (true kalau HTTPS)
		true,           // httpOnly
	)

	c.JSON(200, gin.H{
		"message": "logout berhasil",
	})
}
