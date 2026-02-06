package services

import (
	"errors"
	"smartfarm-api/dto"
	"smartfarm-api/models"
	"smartfarm-api/repositories"
)

func CreateAddress(userID uint, req dto.CreateAddressRequest) (*models.Address, error) {
	address := models.Address{
		UserID:        userID,
		Label:         req.Label,
		RecipientName: req.RecipientName,
		PhoneNumber:   req.PhoneNumber,
		Street:        req.Street,
		City:          req.City,
		Province:      req.Province,
		PostalCode:    req.PostalCode,
		IsDefault:     req.IsDefault,
	}

	if err := repositories.CreateAddress(&address); err != nil {
		return nil, err
	}

	return &address, nil
}

func GetUserAddresses(userID uint) ([]models.Address, error) {
	return repositories.GetUserAddresses(userID)
}

func UpdateAddressService(userID uint, addressID uint, req dto.UpdateAddressRequest) (*models.Address, error) {
	// Verify address belongs to user
	address, err := repositories.GetAddressByID(addressID)
	if err != nil {
		return nil, errors.New("address not found")
	}

	if address.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Update fields
	address.Label = req.Label
	address.RecipientName = req.RecipientName
	address.PhoneNumber = req.PhoneNumber
	address.Street = req.Street
	address.City = req.City
	address.Province = req.Province
	address.PostalCode = req.PostalCode
	address.IsDefault = req.IsDefault

	if err := repositories.UpdateAddress(address); err != nil {
		return nil, err
	}

	return address, nil
}

func DeleteAddressService(userID uint, addressID uint) error {
	// Verify address belongs to user
	address, err := repositories.GetAddressByID(addressID)
	if err != nil {
		return errors.New("address not found")
	}

	if address.UserID != userID {
		return errors.New("unauthorized")
	}

	return repositories.DeleteAddress(addressID)
}

func SetDefaultAddressService(userID uint, addressID uint) error {
	// Verify address belongs to user
	address, err := repositories.GetAddressByID(addressID)
	if err != nil {
		return errors.New("address not found")
	}

	if address.UserID != userID {
		return errors.New("unauthorized")
	}

	return repositories.SetDefaultAddress(userID, addressID)
}
