package service

import (
	"context"
	"fmt"
	"memoir-api/internal/logger"
	"memoir-api/internal/repository"
	"time"
)

// 中国传统节日和西方节日
var festivals = map[string]string{
	"01-14": "情人节", // 2月14日，但月份从0开始
	"05-20": "520",    // 5月20日，中国网络情人节
	"06-07": "七夕",   // 农历七月初七，这里使用公历近似值
}

// 重要纪念日天数
var anniversaryDays = []int{
	7,    // 一周
	30,   // 一个月
	100,  // 100天
	365,  // 一年
	730,  // 两年
	1095, // 三年
	1825, // 五年
	3650, // 十年
}

// CoupleReminderService 情侣纪念日服务接口
type CoupleReminderService interface {
	Service
	// 检查并发送纪念日邮件
	CheckAndSendAnniversaryReminders(ctx context.Context) error
	// 检查并发送节日邮件
	CheckAndSendFestivalReminders(ctx context.Context) error
	// 计算恋爱天数
	CalculateCoupleDays(anniversaryDate time.Time) int
}

// coupleReminderService 情侣纪念日服务实现
type coupleReminderService struct {
	*BaseService
	coupleRepo repository.CoupleRepository
	userRepo   repository.UserRepository
	emailSvc   EmailService
	log        logger.Logger
}

// NewCoupleReminderService 创建情侣纪念日服务实例
func NewCoupleReminderService(
	coupleRepo repository.CoupleRepository,
	userRepo repository.UserRepository,
	emailSvc EmailService,
) CoupleReminderService {
	return &coupleReminderService{
		BaseService: NewBaseService(coupleRepo),
		coupleRepo:  coupleRepo,
		userRepo:    userRepo,
		emailSvc:    emailSvc,
		log:         logger.GetLogger("couple-reminder-service"),
	}
}

// CheckAndSendAnniversaryReminders 检查并发送纪念日邮件
func (s *coupleReminderService) CheckAndSendAnniversaryReminders(ctx context.Context) error {
	s.log.Info("开始检查情侣纪念日")

	// 获取所有情侣关系
	couples, _, err := s.coupleRepo.List(ctx, 0, -1)
	if err != nil {
		s.log.Error(err, "获取情侣列表失败")
		return err
	}

	today := time.Now()

	for _, couple := range couples {
		// 计算恋爱天数
		days := s.CalculateCoupleDays(couple.AnniversaryDate)

		// 检查是否是重要纪念日
		isSpecialDay := false
		for _, specialDay := range anniversaryDays {
			if days == specialDay {
				isSpecialDay = true
				break
			}
		}

		if !isSpecialDay {
			continue
		}

		// 获取情侣用户
		users, err := s.userRepo.ListByCoupleID(ctx, couple.ID)
		if err != nil {
			s.log.Error(err, "获取情侣用户失败", "coupleID", couple.ID)
			continue
		}

		if len(users) != 2 {
			s.log.Warn("情侣用户数量异常", "coupleID", couple.ID, "userCount", len(users))
			continue
		}

		// 发送纪念日邮件给两个用户
		dateStr := today.Format("2006-01-02")

		// 给第一个用户发送邮件
		err = s.emailSvc.SendAnniversaryEmail(
			ctx,
			users[0].Email,
			users[0].Username,
			users[1].Username,
			days,
			dateStr,
		)
		if err != nil {
			s.log.Error(err, "发送纪念日邮件失败", "userID", users[0].ID)
		}

		// 给第二个用户发送邮件
		err = s.emailSvc.SendAnniversaryEmail(
			ctx,
			users[1].Email,
			users[1].Username,
			users[0].Username,
			days,
			dateStr,
		)
		if err != nil {
			s.log.Error(err, "发送纪念日邮件失败", "userID", users[1].ID)
		}

		s.log.Info("已发送纪念日邮件", "coupleID", couple.ID, "days", days)
	}

	s.log.Info("情侣纪念日检查完成")
	return nil
}

// CheckAndSendFestivalReminders 检查并发送节日邮件
func (s *coupleReminderService) CheckAndSendFestivalReminders(ctx context.Context) error {
	s.log.Info("开始检查节日提醒")

	today := time.Now()
	monthDay := fmt.Sprintf("%02d-%02d", today.Month(), today.Day())

	festivalName, isFestival := festivals[monthDay]
	if !isFestival {
		s.log.Info("今天不是特殊节日", "date", monthDay)
		return nil
	}

	// 获取所有情侣关系
	couples, _, err := s.coupleRepo.List(ctx, 0, 1000)
	if err != nil {
		s.log.Error(err, "获取情侣列表失败")
		return err
	}

	for _, couple := range couples {
		// 获取情侣用户
		users, err := s.userRepo.ListByCoupleID(ctx, couple.ID)
		if err != nil {
			s.log.Error(err, "获取情侣用户失败", "coupleID", couple.ID)
			continue
		}

		if len(users) != 2 {
			s.log.Warn("情侣用户数量异常", "coupleID", couple.ID, "userCount", len(users))
			continue
		}

		// 给第一个用户发送节日邮件
		err = s.emailSvc.SendFestivalEmail(
			ctx,
			users[0].Email,
			users[0].Username,
			users[1].Username,
			festivalName,
		)
		if err != nil {
			s.log.Error(err, "发送节日邮件失败", "userID", users[0].ID)
		}

		// 给第二个用户发送节日邮件
		err = s.emailSvc.SendFestivalEmail(
			ctx,
			users[1].Email,
			users[1].Username,
			users[0].Username,
			festivalName,
		)
		if err != nil {
			s.log.Error(err, "发送节日邮件失败", "userID", users[1].ID)
		}

		s.log.Info("已发送节日邮件", "coupleID", couple.ID, "festival", festivalName)
	}

	s.log.Info("节日提醒检查完成")
	return nil
}

// CalculateCoupleDays 计算恋爱天数
func (s *coupleReminderService) CalculateCoupleDays(anniversaryDate time.Time) int {
	now := time.Now()
	duration := now.Sub(anniversaryDate)
	days := int(duration.Hours() / 24)
	return days
}
