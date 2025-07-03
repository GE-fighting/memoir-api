package email

import (
	"fmt"
	"strings"
)

// 邮箱验证邮件模板
func renderVerificationEmailTemplate(data map[string]string) string {
	template := `
<div style="max-width:600px;margin:0 auto;font-family:Arial,sans-serif;">
    <div style="background:#f8f9fa;padding:20px;text-align:center;">
        <h1 style="color:#333;">{{AppName}} - 邮箱验证</h1>
    </div>
    <div style="padding:30px;">
        <p>您好 <strong>{{Username}}</strong>，</p>
        <p>请使用以下验证码完成邮箱验证：</p>
        <div style="background:#007bff;color:white;padding:15px;text-align:center;font-size:24px;font-weight:bold;margin:20px 0;">
            {{VerificationCode}}
        </div>
        <p style="color:#666;">验证码将在 <strong>{{ExpireMinutes}}</strong> 分钟后过期，请及时使用。</p>
        <p style="color:#999;font-size:12px;">如果您没有注册账户，请忽略此邮件。</p>
    </div>
    <div style="background:#f8f9fa;padding:15px;text-align:center;font-size:12px;color:#666;">
        <p>&copy; {{AppName}}. 保留所有权利。</p>
    </div>
</div>`

	return renderTemplate(template, data)
}

// 密码重置邮件模板
func renderPasswordResetEmailTemplate(data map[string]string) string {
	template := `
<div style="max-width:600px;margin:0 auto;font-family:Arial,sans-serif;">
    <div style="background:#f8f9fa;padding:20px;text-align:center;">
        <h1 style="color:#333;">{{AppName}} - 密码重置</h1>
    </div>
    <div style="padding:30px;">
        <p>您好，</p>
        <p>我们收到了您的密码重置请求。请点击下面的链接重置您的密码：</p>
        <div style="text-align:center;margin:30px 0;">
            <a href="{{ResetLink}}" style="background:#28a745;color:white;padding:12px 30px;text-decoration:none;border-radius:5px;">重置密码</a>
        </div>
        <p style="color:#666;">此链接将在 <strong>30</strong> 分钟后过期。</p>
        <p style="color:#999;font-size:12px;">如果您没有请求重置密码，请忽略此邮件。</p>
    </div>
    <div style="background:#f8f9fa;padding:15px;text-align:center;font-size:12px;color:#666;">
        <p>&copy; {{AppName}}. 保留所有权利。</p>
    </div>
</div>`

	return renderTemplate(template, data)
}

// 系统通知邮件模板
func renderNotificationEmailTemplate(data map[string]string) string {
	template := `
<div style="max-width:600px;margin:0 auto;font-family:Arial,sans-serif;">
    <div style="background:#f8f9fa;padding:20px;text-align:center;">
        <h1 style="color:#333;">{{AppName}} - 系统通知</h1>
    </div>
    <div style="padding:30px;">
        <p>您好 <strong>{{Username}}</strong>，</p>
        <p>{{Message}}</p>
        <p style="color:#666;margin-top:30px;">此邮件由系统自动发送，请勿回复。</p>
    </div>
    <div style="background:#f8f9fa;padding:15px;text-align:center;font-size:12px;color:#666;">
        <p>&copy; {{AppName}}. 保留所有权利。</p>
    </div>
</div>`

	return renderTemplate(template, data)
}

// 欢迎邮件模板
func renderWelcomeEmailTemplate(data map[string]string) string {
	template := `
<div style="max-width:600px;margin:0 auto;font-family:Arial,sans-serif;">
    <div style="background:#f8f9fa;padding:20px;text-align:center;">
        <h1 style="color:#333;">欢迎加入 {{AppName}}</h1>
    </div>
    <div style="padding:30px;">
        <p>您好 <strong>{{Username}}</strong>，</p>
        <p>感谢您注册 {{AppName}}！我们很高兴您加入我们的社区。</p>
        <p>您现在可以使用您的账号访问所有功能：</p>
        <div style="text-align:center;margin:30px 0;">
            <a href="{{AppURL}}" style="background:#007bff;color:white;padding:12px 30px;text-decoration:none;border-radius:5px;">访问 {{AppName}}</a>
        </div>
        <p>如果您有任何问题，请随时联系我们的支持团队。</p>
    </div>
    <div style="background:#f8f9fa;padding:15px;text-align:center;font-size:12px;color:#666;">
        <p>&copy; {{AppName}}. 保留所有权利。</p>
    </div>
</div>`

	return renderTemplate(template, data)
}

// 纪念日邮件模板
func renderAnniversaryEmailTemplate(data map[string]string) string {
	template := `
<div style="max-width:600px;margin:0 auto;font-family:Arial,sans-serif;">
    <div style="background:#f8f9fa;padding:20px;text-align:center;">
        <h1 style="color:#e91e63;">❤️ 甜蜜纪念日 ❤️</h1>
    </div>
    <div style="padding:30px;text-align:center;">
        <div style="margin-bottom:20px;">
            <img src="https://img.icons8.com/color/96/000000/hearts.png" alt="Hearts" style="width:80px;height:80px;">
        </div>
        <h2 style="color:#e91e63;margin-bottom:20px;">今天是您和{{PartnerName}}在一起的第{{Days}}天！</h2>
        <p style="font-size:18px;color:#555;margin-bottom:20px;">亲爱的 <strong>{{Username}}</strong>，</p>
        <p style="font-size:16px;color:#555;margin-bottom:20px;">在这特别的日子里，{{AppName}}想要送上我们最真挚的祝福！</p>
        <div style="background:#ffe8f0;border-radius:8px;padding:20px;margin:25px 0;text-align:left;">
            <p style="font-size:16px;line-height:1.6;color:#333;">🌹 <strong>恋爱是一场美丽的旅程</strong>，而您们已经一同走过了{{Days}}天。每一天都是珍贵的回忆，每一刻都值得铭记和庆祝。</p>
            <p style="font-size:16px;line-height:1.6;color:#333;">💫 希望您们能够用心记录这美好的时光，创造更多动人的瞬间。</p>
        </div>
        <div style="text-align:center;margin:30px 0;">
            <a href="{{AppURL}}" style="background:#e91e63;color:white;padding:12px 30px;text-decoration:none;border-radius:5px;">记录美好时光</a>
        </div>
        <p style="font-style:italic;color:#888;">记得和{{PartnerName}}分享这一刻，一起庆祝你们的爱情故事！</p>
    </div>
    <div style="background:#f8f9fa;padding:15px;text-align:center;font-size:12px;color:#666;">
        <p>❤️ {{AppName}} 祝您们爱情甜蜜，幸福长久！</p>
        <p>&copy; {{AppName}}. 保留所有权利。</p>
    </div>
</div>`

	return renderTemplate(template, data)
}

// 节日邮件模板
func renderFestivalEmailTemplate(data map[string]string) string {
	template := `
<div style="max-width:600px;margin:0 auto;font-family:Arial,sans-serif;">
    <div style="background:#f8f9fa;padding:20px;text-align:center;">
        <h1 style="color:#9c27b0;">💖 {{FestivalName}}快乐 💖</h1>
    </div>
    <div style="padding:30px;text-align:center;">
        <div style="margin-bottom:20px;">
            <img src="https://img.icons8.com/color/96/000000/gift.png" alt="Gift" style="width:80px;height:80px;">
        </div>
        <h2 style="color:#9c27b0;margin-bottom:20px;">亲爱的 {{Username}}</h2>
        <p style="font-size:18px;color:#555;margin-bottom:20px;">{{AppName}} 祝您和 {{PartnerName}} {{FestivalName}}快乐！</p>
        <div style="background:#f3e5f5;border-radius:8px;padding:20px;margin:25px 0;text-align:left;">
            <p style="font-size:16px;line-height:1.6;color:#333;">🌟 在这个特别的日子里，愿你们的爱情如星光般闪耀，温暖彼此的心灵。</p>
            <p style="font-size:16px;line-height:1.6;color:#333;">🎁 每一个节日都是庆祝爱情的机会，希望这一天能为你们的感情增添美好的回忆。</p>
            <p style="font-size:16px;line-height:1.6;color:#333;">💕 珍惜当下，用心感受彼此的陪伴，这是最珍贵的礼物。</p>
        </div>
        <div style="text-align:center;margin:30px 0;">
            <a href="{{AppURL}}" style="background:#9c27b0;color:white;padding:12px 30px;text-decoration:none;border-radius:5px;">浪漫相册</a>
        </div>
        <p style="font-style:italic;color:#888;">希望您们能一起度过一个难忘的{{FestivalName}}！</p>
    </div>
    <div style="background:#f8f9fa;padding:15px;text-align:center;font-size:12px;color:#666;">
        <p>💖 {{AppName}} 祝您们幸福美满！</p>
        <p>&copy; {{AppName}}. 保留所有权利。</p>
    </div>
</div>`

	return renderTemplate(template, data)
}

// 通用模板渲染函数
func renderTemplate(template string, data map[string]string) string {
	result := template
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}
