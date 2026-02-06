package repositories

import (
	"smartfarm-api/config"
	"smartfarm-api/models"

	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(address *models.Address) error
	GetUserAddresses(userID uint) ([]models.Address, error)
	FindByID(id uint) (*models.Address, error)
	Update(address *models.Address) error
	Delete(id uint) error
	SetDefault(userID uint, addressID uint) error
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db}
}

func (r *addressRepository) Create(address *models.Address) error {
	return r.db.Create(address).Error
}

func (r *addressRepository) GetUserAddresses(userID uint) ([]models.Address, error) {
	var addresses []models.Address
	err := r.db.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error
	return addresses, err
}

func (r *addressRepository) FindByID(id uint) (*models.Address, error) {
	var address models.Address
	err := r.db.First(&address, id).Error
	return &address, err
}

func (r *addressRepository) Update(address *models.Address) error {
	return r.db.Save(address).Error
}

func (r *addressRepository) Delete(id uint) error {
	return r.db.Delete(&models.Address{}, id).Error
}

func (r *addressRepository) SetDefault(userID uint, addressID uint) error {
	// Unset previous default
	if err := r.db.Model(&models.Address{}).Where("user_id = ? AND id != ?", userID, addressID).Update("is_default", false).Error; err != nil {
		return err
	}
	// Set new default
	return r.db.Model(&models.Address{}).Where("id = ?", addressID).Update("is_default", true).Error
}

// Keep legacy function wrappers for backward compatibility
func CreateAddress(address *models.Address) error {
	return config.DB.Create(address).Error
}

func GetUserAddresses(userID uint) ([]models.Address, error) {
	var addresses []models.Address
	err := config.DB.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error
	return addresses, err
}

func GetAddressByID(id uint) (*models.Address, error) {
	var address models.Address
	err := config.DB.First(&address, id).Error
	return &address, err
}

func UpdateAddress(address *models.Address) error {
	return config.DB.Save(address).Error
}

func DeleteAddress(id uint) error {
	return config.DB.Delete(&models.Address{}, id).Error
}

func SetDefaultAddress(userID uint, addressID uint) error {
	// Unset previous default
	if err := config.DB.Model(&models.Address{}).Where("user_id = ? AND id != ?", userID, addressID).Update("is_default", false).Error; err != nil {
		return err
	}
	// Set new default
	return config.DB.Model(&models.Address{}).Where("id = ?", addressID).Update("is_default", true).Error
}
