package repository

import (
	"gorm.io/gorm"
)

// Factory 仓库工厂接口
type Factory interface {
	Couple() CoupleRepository
	User() UserRepository
	Location() LocationRepository
	TimelineEvent() TimelineEventRepository
	PhotoVideo() PhotoVideoRepository
	Wishlist() WishlistRepository
}

// factory 仓库工厂实现
type factory struct {
	coupleRepository        CoupleRepository
	userRepository          UserRepository
	locationRepository      LocationRepository
	timelineEventRepository TimelineEventRepository
	photoVideoRepository    PhotoVideoRepository
	wishlistRepository      WishlistRepository
}

// NewFactory 创建仓库工厂
func NewFactory(db *gorm.DB) Factory {
	return &factory{
		coupleRepository:        NewCoupleRepository(db),
		userRepository:          NewUserRepository(db),
		locationRepository:      NewLocationRepository(db),
		timelineEventRepository: NewTimelineEventRepository(db),
		photoVideoRepository:    NewPhotoVideoRepository(db),
		wishlistRepository:      NewWishlistRepository(db),
	}
}

// Couple 获取情侣关系仓库
func (f *factory) Couple() CoupleRepository {
	return f.coupleRepository
}

// User 获取用户仓库
func (f *factory) User() UserRepository {
	return f.userRepository
}

// Location 获取地点仓库
func (f *factory) Location() LocationRepository {
	return f.locationRepository
}

// TimelineEvent 获取时间轴事件仓库
func (f *factory) TimelineEvent() TimelineEventRepository {
	return f.timelineEventRepository
}

// PhotoVideo 获取照片和视频仓库
func (f *factory) PhotoVideo() PhotoVideoRepository {
	return f.photoVideoRepository
}

// Wishlist 获取心愿清单仓库
func (f *factory) Wishlist() WishlistRepository {
	return f.wishlistRepository
}
