package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNotificationTemplate_TableName(t *testing.T) {
	template := NotificationTemplate{}
	assert.Equal(t, "notification_templates", template.TableName())
}

func TestUserNotification_TableName(t *testing.T) {
	notification := UserNotification{}
	assert.Equal(t, "user_notifications", notification.TableName())
}

func TestUserNotification_IsRead(t *testing.T) {
	tests := []struct {
		name     string
		notif    UserNotification
		expected bool
	}{
		{
			name: "status is read",
			notif: UserNotification{
				Status: NotificationStatusRead,
			},
			expected: true,
		},
		{
			name: "readAt is set",
			notif: UserNotification{
				Status: NotificationStatusSent,
				ReadAt: func() *time.Time { t := time.Now(); return &t }(),
			},
			expected: true,
		},
		{
			name: "not read - pending status",
			notif: UserNotification{
				Status: NotificationStatusPending,
			},
			expected: false,
		},
		{
			name: "not read - sent status without readAt",
			notif: UserNotification{
				Status: NotificationStatusSent,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.notif.IsRead())
		})
	}
}

func TestUserNotificationSetting_TableName(t *testing.T) {
	setting := UserNotificationSetting{}
	assert.Equal(t, "user_notification_settings", setting.TableName())
}

func TestUserNotificationSetting_IsInDoNotDisturbPeriod(t *testing.T) {
	tests := []struct {
		name     string
		setting  UserNotificationSetting
		expected bool
	}{
		{
			name: "disabled",
			setting: UserNotificationSetting{
				DoNotDisturbEnabled: false,
				DoNotDisturbStart:   "22:00",
				DoNotDisturbEnd:     "08:00",
			},
			expected: false,
		},
		{
			name: "empty start time",
			setting: UserNotificationSetting{
				DoNotDisturbEnabled: true,
				DoNotDisturbStart:   "",
				DoNotDisturbEnd:     "08:00",
			},
			expected: false,
		},
		{
			name: "empty end time",
			setting: UserNotificationSetting{
				DoNotDisturbEnabled: true,
				DoNotDisturbStart:   "22:00",
				DoNotDisturbEnd:     "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.setting.IsInDoNotDisturbPeriod())
		})
	}
}

func TestNotificationConfig_TableName(t *testing.T) {
	config := NotificationConfig{}
	assert.Equal(t, "notification_configs", config.TableName())
}

func TestNotificationSchedule_TableName(t *testing.T) {
	schedule := NotificationSchedule{}
	assert.Equal(t, "notification_schedules", schedule.TableName())
}

func TestNotificationTypeConstants(t *testing.T) {
	assert.Equal(t, NotificationType("order_status"), NotificationTypeOrderStatus)
	assert.Equal(t, NotificationType("vip_expire"), NotificationTypeVipExpire)
	assert.Equal(t, NotificationType("coupon_expire"), NotificationTypeCouponExpire)
	assert.Equal(t, NotificationType("activity_start"), NotificationTypeActivityStart)
	assert.Equal(t, NotificationType("activity_end"), NotificationTypeActivityEnd)
	assert.Equal(t, NotificationType("system"), NotificationTypeSystem)
	assert.Equal(t, NotificationType("promotion"), NotificationTypePromotion)
	assert.Equal(t, NotificationType("chat"), NotificationTypeChat)
}

func TestNotificationChannelConstants(t *testing.T) {
	assert.Equal(t, NotificationChannel("in_app"), NotificationChannelInApp)
	assert.Equal(t, NotificationChannel("push"), NotificationChannelPush)
	assert.Equal(t, NotificationChannel("sms"), NotificationChannelSMS)
	assert.Equal(t, NotificationChannel("wechat"), NotificationChannelWechat)
	assert.Equal(t, NotificationChannel("email"), NotificationChannelEmail)
}

func TestNotificationStatusConstants(t *testing.T) {
	assert.Equal(t, NotificationStatus("pending"), NotificationStatusPending)
	assert.Equal(t, NotificationStatus("sent"), NotificationStatusSent)
	assert.Equal(t, NotificationStatus("read"), NotificationStatusRead)
	assert.Equal(t, NotificationStatus("failed"), NotificationStatusFailed)
	assert.Equal(t, NotificationStatus("canceled"), NotificationStatusCanceled)
}

func TestNotificationConfigConstants(t *testing.T) {
	assert.Equal(t, "vip_expire_days", NotificationConfigVipExpireDays)
	assert.Equal(t, "coupon_expire_days", NotificationConfigCouponExpireDays)
	assert.Equal(t, "push_provider", NotificationConfigPushProvider)
	assert.Equal(t, "sms_provider", NotificationConfigSMSProvider)
}
