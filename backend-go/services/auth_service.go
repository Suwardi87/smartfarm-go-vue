package services

import (
	"errors"

	"smartfarm-api/dto"
	"smartfarm-api/models"
	"smartfarm-api/repositories"
	"smartfarm-api/utils"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(req dto.RegisterRequest) error {

	// cek email sudah ada
	existingUser, _ := repositories.FindUserByEmail(req.Email)
	if existingUser.ID != 0 {
		return errors.New("email sudah terdaftar")
	}

	// hash password (pengamanan password)
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	return repositories.CreateUser(&user)
}

func LoginUser(req dto.LoginRequest) (string, error) {

	user, err := repositories.FindUserByEmail(req.Email)
	if err != nil {
		return "", errors.New("email tidak ditemukan")
	}

	// compare password (bandingkan hash)
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return "", errors.New("password salah")
	}

	// generate JWT
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GetUserByID(id uint) (*models.User, error) {
	return repositories.FindUserByID(id)
}

func UpdateUserProfile(id uint, req dto.UpdateProfileRequest) error {
	user, err := repositories.FindUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	user.Name = req.Name
	user.Email = req.Email

	return repositories.UpdateUser(user)
}
