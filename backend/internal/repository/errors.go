package repository

import "errors"

// ErrNotFound 表示记录不存在。
var ErrNotFound = errors.New("record not found")
