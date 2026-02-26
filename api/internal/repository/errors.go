package repository

import (
	"errors"

	"gorm.io/gorm"
)

// ErrNotFound 表示记录不存在。
var ErrNotFound = errors.New("record not found")

// ErrInvalidStatusTransition indicates a status transition violates domain rules.
var ErrInvalidStatusTransition = errors.New("invalid status transition")

// WrapNotFound 将 gorm.ErrRecordNotFound 转换为 repository.ErrNotFound。
// 如果 err 是 gorm.ErrRecordNotFound，返回 ErrNotFound；
// 如果 err 为 nil，返回 nil；
// 否则返回原始错误。
func WrapNotFound(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}

// IsNotFound 检查错误是否为 "记录不存在" 错误。
// 同时检查 repository.ErrNotFound 和 gorm.ErrRecordNotFound。
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound)
}

// HandleGetError 处理单条记录查询的错误。
// 如果记录不存在返回 (nil, ErrNotFound)；
// 如果有其他错误返回 (nil, err)；
// 如果成功返回 (result, nil)。
func HandleGetError[T any](result *T, err error) (*T, error) {
	if err != nil {
		return nil, WrapNotFound(err)
	}
	return result, nil
}
