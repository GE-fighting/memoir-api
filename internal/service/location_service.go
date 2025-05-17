package service

import (
	"context"
	"errors"
	"fmt"

	"memoir-api/internal/models"
	"memoir-api/internal/repository"
)

var (
	ErrLocationNotFound = errors.New("地点不存在")
)

// LocationService 地点服务接口
type LocationService interface {
	Service
	CreateLocation(ctx context.Context, location *models.Location) (*models.Location, error)
	GetLocationByID(ctx context.Context, id int64) (*models.Location, error)
	ListLocationsByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.Location, int64, error)
	UpdateLocation(ctx context.Context, location *models.Location) error
	DeleteLocation(ctx context.Context, id int64) error
	FindNearbyLocations(ctx context.Context, lat, lng float64, radiusMeters float64, limit int) ([]*models.Location, error)
}

// locationService 地点服务实现
type locationService struct {
	*BaseService
	locationRepo repository.LocationRepository
}

// NewLocationService 创建地点服务
func NewLocationService(locationRepo repository.LocationRepository) LocationService {
	return &locationService{
		BaseService:  NewBaseService(locationRepo),
		locationRepo: locationRepo,
	}
}

// CreateLocation 创建地点
func (s *locationService) CreateLocation(ctx context.Context, location *models.Location) (*models.Location, error) {
	if err := s.locationRepo.Create(ctx, location); err != nil {
		return nil, fmt.Errorf("创建地点失败: %w", err)
	}
	return location, nil
}

// GetLocationByID 通过ID获取地点
func (s *locationService) GetLocationByID(ctx context.Context, id int64) (*models.Location, error) {
	location, err := s.locationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			return nil, ErrLocationNotFound
		}
		return nil, fmt.Errorf("获取地点失败: %w", err)
	}
	return location, nil
}

// ListLocationsByCoupleID 获取情侣关系下的所有地点
func (s *locationService) ListLocationsByCoupleID(ctx context.Context, coupleID int64, offset, limit int) ([]*models.Location, int64, error) {
	return s.locationRepo.ListByCoupleID(ctx, coupleID, offset, limit)
}

// UpdateLocation 更新地点
func (s *locationService) UpdateLocation(ctx context.Context, location *models.Location) error {
	// 检查地点是否存在
	_, err := s.locationRepo.GetByID(ctx, location.ID)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			return ErrLocationNotFound
		}
		return fmt.Errorf("查询地点失败: %w", err)
	}

	// 更新地点
	if err := s.locationRepo.Update(ctx, location); err != nil {
		return fmt.Errorf("更新地点失败: %w", err)
	}
	return nil
}

// DeleteLocation 删除地点
func (s *locationService) DeleteLocation(ctx context.Context, id int64) error {
	// 检查地点是否存在
	_, err := s.locationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			return ErrLocationNotFound
		}
		return fmt.Errorf("查询地点失败: %w", err)
	}

	// 删除地点
	if err := s.locationRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除地点失败: %w", err)
	}
	return nil
}

// FindNearbyLocations 查找附近地点
func (s *locationService) FindNearbyLocations(ctx context.Context, lat, lng float64, radiusMeters float64, limit int) ([]*models.Location, error) {
	locations, err := s.locationRepo.FindNearby(ctx, lat, lng, radiusMeters, limit)
	if err != nil {
		return nil, fmt.Errorf("查找附近地点失败: %w", err)
	}
	return locations, nil
}
