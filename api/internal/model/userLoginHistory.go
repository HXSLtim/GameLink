package model

import "time"

// UserLoginHistory 登录历史
// 记录用户登录事件与设备信息
type UserLoginHistory struct {
    Base
    UserID      uint64     `json:"userId" gorm:"column:user_id;index;not null;comment:用户ID"`
    IPAddress   string     `json:"ipAddress" gorm:"column:ip_address;size:45;comment:登录IP地址"`
    UserAgent   string     `json:"userAgent" gorm:"column:user_agent;type:text;comment:UA字符串"`
    DeviceType  string     `json:"deviceType" gorm:"column:device_type;size:32;comment:设备类型"`
    DeviceInfo  string     `json:"deviceInfo" gorm:"column:device_info;type:text;comment:设备详情"`
    Location    string     `json:"location" gorm:"size:255;comment:地理位置"`
    LoginResult string     `json:"loginResult" gorm:"column:login_result;size:32;comment:登录结果"`
    SessionID   string     `json:"sessionId" gorm:"column:session_id;size:128;comment:会话ID"`
    LogoutAt    *time.Time `json:"logoutAt,omitempty" gorm:"column:logout_at;comment:登出时间"`
}

// TableName 指定表名
func (UserLoginHistory) TableName() string { return "user_login_histories" }

