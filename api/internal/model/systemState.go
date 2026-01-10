package model

import (
	"encoding/json"
	"time"
)

// SystemState represents system initialization and version state
// It tracks when the system was last initialized with menu/permission sync
type SystemState struct {
	Base
	Key         string    `json:"key" gorm:"size:64;uniqueIndex;not null"` // e.g., "admin_init"
	Value       string    `json:"value" gorm:"type:text;not null"`           // JSON string containing state data
	ExpiresAt   *time.Time `json:"expiresAt,omitempty" gorm:"index"`          // Optional expiration for cache-like behavior
	Version     string    `json:"version" gorm:"size:32"`                    // Version hash for change detection
	LastSyncAt  time.Time `json:"lastSyncAt" gorm:"index"`                   // When sync was performed
	SyncedBy    uint64    `json:"syncedBy"`                                  // User ID who performed the sync
	SyncedByIP  string    `json:"syncedByIp" gorm:"size:64"`                 // IP address of sync requester
	Description string    `json:"description" gorm:"type:text"`              // Human-readable description
}

// TableName specifies the table name for SystemState
func (SystemState) TableName() string {
	return "system_states"
}

// SystemInitData represents the initialization data stored in Value field
type SystemInitData struct {
	MenuCount      int       `json:"menuCount"`
	PermissionCount int      `json:"permissionCount"`
	MenuVersion    string    `json:"menuVersion"`
	PermVersion    string    `json:"permVersion"`
}

// GetInitData parses the Value JSON into SystemInitData
func (s *SystemState) GetInitData() (*SystemInitData, error) {
	var data SystemInitData
	if err := json.Unmarshal([]byte(s.Value), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// SetInitData sets the Value field from SystemInitData
func (s *SystemState) SetInitData(data *SystemInitData) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	s.Value = string(bytes)
	return nil
}

// SystemStateKey constants
const (
	SystemStateKeyAdminInit = "admin_init" // Admin panel initialization state
)

// NewAdminInitState creates a new admin initialization state
func NewAdminInitState(menuCount, permCount int, menuVersion, permVersion string, userID uint64, ip string) *SystemState {
	data := &SystemInitData{
		MenuCount:       menuCount,
		PermissionCount: permCount,
		MenuVersion:     menuVersion,
		PermVersion:     permVersion,
	}
	value, _ := json.Marshal(data)
	return &SystemState{
		Key:         SystemStateKeyAdminInit,
		Value:       string(value),
		Version:     menuVersion + ":" + permVersion,
		LastSyncAt:  time.Now(),
		SyncedBy:    userID,
		SyncedByIP:  ip,
		Description: "Admin panel menu and permission synchronization",
	}
}
