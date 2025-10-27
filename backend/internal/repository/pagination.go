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
