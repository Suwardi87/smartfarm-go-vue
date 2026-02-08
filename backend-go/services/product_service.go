package services

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"smartfarm-api/dto"
	"smartfarm-api/models"
	"smartfarm-api/repositories"
	"time"
)

type ProductService interface {
	CreateProduct(req dto.CreateProductRequest, farmerID uint) (dto.ProductResponse, error)
	FindAll(query string, page int, limit int) (dto.PaginatedProductResponse, error)
	FindByID(id uint) (dto.ProductResponse, error)
	FindProductsByFarmerID(farmerID uint, page int, limit int) (dto.PaginatedProductResponse, error)
	UpdateProduct(id uint, req dto.CreateProductRequest, farmerID uint) (dto.ProductResponse, error)
	DeleteProduct(id uint, farmerID uint) error
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo}
}

func (s *productService) CreateProduct(req dto.CreateProductRequest, farmerID uint) (dto.ProductResponse, error) {
	// Handle Image Upload
	imageURL := ""
	if req.Image != nil {
		filename, err := saveImage(req.Image)
		if err != nil {
			return dto.ProductResponse{}, err
		}
		imageURL = filename
	}

	// Parse Harvest Date if exists
	var harvestDate *time.Time
	if req.IsPreOrder && req.HarvestDate != "" {
		parsed, err := time.Parse("2006-01-02", req.HarvestDate)
		if err == nil {
			harvestDate = &parsed
		}
	}

	product := models.Product{
		Name:               req.Name,
		Description:        req.Description,
		Price:              req.Price,
		Stock:              req.Stock,
		ImageURL:           imageURL,
		Category:           req.Category,
		FarmerID:           farmerID,
		IsPreOrder:         req.IsPreOrder,
		HarvestDate:        harvestDate,
		IsSubscription:     req.IsSubscription,
		SubscriptionPeriod: req.SubscriptionPeriod,
	}

	err := s.repo.Create(&product)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return mapProductToResponse(product), nil
}

func (s *productService) FindAll(query string, page int, limit int) (dto.PaginatedProductResponse, error) {
	start := time.Now()
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 12
	}

	offset := (page - 1) * limit

	products, err := s.repo.FindAll(query, limit, offset)
	if err != nil {
		return dto.PaginatedProductResponse{}, err
	}
	log.Printf("[DEBUG] repo.FindAll took %v", time.Since(start))

	countStart := time.Now()
	total, err := s.repo.CountAll(query)
	if err != nil {
		return dto.PaginatedProductResponse{}, err
	}
	log.Printf("[DEBUG] repo.CountAll took %v", time.Since(countStart))

	var responses []dto.ProductResponse
	for _, p := range products {
		responses = append(responses, mapProductToResponse(p))
	}

	totalPages := int(total / int64(limit))
	if total%int64(limit) > 0 {
		totalPages++
	}

	return dto.PaginatedProductResponse{
		Data:       responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *productService) FindByID(id uint) (dto.ProductResponse, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	return mapProductToResponse(product), nil
}

func (s *productService) FindProductsByFarmerID(farmerID uint, page int, limit int) (dto.PaginatedProductResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 12
	}

	offset := (page - 1) * limit

	products, err := s.repo.FindByFarmerID(farmerID, limit, offset)
	if err != nil {
		return dto.PaginatedProductResponse{}, err
	}

	total, err := s.repo.CountByFarmerID(farmerID)
	if err != nil {
		return dto.PaginatedProductResponse{}, err
	}

	var responses []dto.ProductResponse
	for _, p := range products {
		responses = append(responses, mapProductToResponse(p))
	}

	totalPages := int(total / int64(limit))
	if total%int64(limit) > 0 {
		totalPages++
	}

	return dto.PaginatedProductResponse{
		Data:       responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *productService) UpdateProduct(id uint, req dto.CreateProductRequest, farmerID uint) (dto.ProductResponse, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	if product.FarmerID != farmerID {
		return dto.ProductResponse{}, fmt.Errorf("unauthorized to update this product")
	}

	// Update fields
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.Category = req.Category
	product.IsPreOrder = req.IsPreOrder
	product.IsSubscription = req.IsSubscription
	product.SubscriptionPeriod = req.SubscriptionPeriod

	if req.Image != nil {
		filename, err := saveImage(req.Image)
		if err == nil {
			product.ImageURL = filename
		}
	}

	if req.IsPreOrder && req.HarvestDate != "" {
		parsed, err := time.Parse("2006-01-02", req.HarvestDate)
		if err == nil {
			product.HarvestDate = &parsed
		}
	}

	err = s.repo.Update(&product)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return mapProductToResponse(product), nil
}

func (s *productService) DeleteProduct(id uint, farmerID uint) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if product.FarmerID != farmerID {
		return fmt.Errorf("unauthorized to delete this product")
	}

	return s.repo.Delete(id)
}

// Helpers
func saveImage(file *multipart.FileHeader) (string, error) {
	uploadDir := "uploads/products"

	// pastikan folder ada
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return "", err
	}

	return fmt.Sprintf("products/%s", filename), nil
}

func mapProductToResponse(p models.Product) dto.ProductResponse {
	var harvestDateStr string
	if p.HarvestDate != nil {
		harvestDateStr = p.HarvestDate.Format("2006-01-02")
	}

	return dto.ProductResponse{
		ID:                 p.ID,
		Name:               p.Name,
		Description:        p.Description,
		Price:              p.Price,
		Stock:              p.Stock,
		ImageURL:           p.ImageURL,
		Category:           p.Category,
		FarmerID:           p.FarmerID,
		FarmerName:         p.Farmer.Name,
		IsPreOrder:         p.IsPreOrder,
		HarvestDate:        harvestDateStr,
		IsSubscription:     p.IsSubscription,
		SubscriptionPeriod: p.SubscriptionPeriod,
	}
}
