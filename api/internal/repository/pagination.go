package repository

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// NormalizePage 返回合法的页码。
func NormalizePage(page int) int {
	if page <= 0 {
		return defaultPage
	}
	return page
}

// NormalizePageSize 返回合法的分页大小。
func NormalizePageSize(size int) int {
	if size <= 0 {
		return defaultPageSize
	}
	if size > maxPageSize {
		return maxPageSize
	}
	return size
}

// BuildPagination 构建分页响应对象。
// 自动规范化 page 和 pageSize 参数。
func BuildPagination(page, pageSize int, total int64) (int, int, int, int, bool, bool) {
	page = NormalizePage(page)
	pageSize = NormalizePageSize(pageSize)

	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	hasNext := page < totalPages
	hasPrev := page > 1

	return page, pageSize, int(total), totalPages, hasNext, hasPrev
}
