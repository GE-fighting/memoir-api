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
	PersonalMedia() PersonalMediaRepository
	CoupleAlbum() CoupleAlbumRepository
	TimelineEventLocation() TimelineEventLocationRepository
	TimelineEventPhotoVideo() TimelineEventPhotoVideoRepository
	Attachment() AttachmentRepository
	WishlistAttachment() WishlistAttachmentRepository
	GetDB() *gorm.DB
}

// factory 仓库工厂实现
type factory struct {
	db                                *gorm.DB
	coupleRepository                  CoupleRepository
	userRepository                    UserRepository
	locationRepository                LocationRepository
	timelineEventRepository           TimelineEventRepository
	photoVideoRepository              PhotoVideoRepository
	wishlistRepository                WishlistRepository
	personalMediaRepository           PersonalMediaRepository
	coupleAlbumRepository             CoupleAlbumRepository
	timelineEventLocationRepository   TimelineEventLocationRepository
	timelineEventPhotoVideoRepository TimelineEventPhotoVideoRepository
	attachmentRepository              AttachmentRepository
	wishlistAttachmentRepository      WishlistAttachmentRepository
}

func (f *factory) TimelineEventLocation() TimelineEventLocationRepository {
	return f.timelineEventLocationRepository
}

func (f *factory) TimelineEventPhotoVideo() TimelineEventPhotoVideoRepository {
	return f.timelineEventPhotoVideoRepository
}

// NewFactory 创建仓库工厂
func NewFactory(db *gorm.DB) Factory {
	return &factory{
		db:                                db,
		coupleRepository:                  NewCoupleRepository(db),
		userRepository:                    NewUserRepository(db),
		locationRepository:                NewLocationRepository(db),
		timelineEventRepository:           NewTimelineEventRepository(db),
		photoVideoRepository:              NewPhotoVideoRepository(db),
		wishlistRepository:                NewWishlistRepository(db),
		personalMediaRepository:           NewGormPersonalMediaRepository(db),
		coupleAlbumRepository:             NewCoupleAlbumRepository(db),
		timelineEventLocationRepository:   NewTimelineEventLocationRepository(db),
		timelineEventPhotoVideoRepository: NewTimelineEventPhotoVideoRepository(db),
		attachmentRepository:              NewAttachmentRepository(db),
		wishlistAttachmentRepository:      NewWishlistAttachmentRepository(db),
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

// PersonalMedia 获取个人媒体仓库
func (f *factory) PersonalMedia() PersonalMediaRepository {
	return f.personalMediaRepository
}

// CoupleAlbum 获取情侣相册仓库
func (f *factory) CoupleAlbum() CoupleAlbumRepository {
	return f.coupleAlbumRepository
}

// Attachment 获取附件仓库
func (f *factory) Attachment() AttachmentRepository {
	return f.attachmentRepository
}

// WishlistAttachment 获取心愿附件关联仓库
func (f *factory) WishlistAttachment() WishlistAttachmentRepository {
	return f.wishlistAttachmentRepository
}

// GetDB 获取数据库连接
func (f *factory) GetDB() *gorm.DB {
	return f.db
}
