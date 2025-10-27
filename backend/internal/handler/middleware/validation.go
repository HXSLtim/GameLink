package middleware

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 验证器实例
var validate = validator.New()

// ValidationError 表示验证失败的错误
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

// ValidateJSON 验证JSON请求体
//
// 使用方法：
// router.POST("/users", middleware.ValidateJSON(&CreateUserRequest{}), createUserHandler)
//
// 技术原理：
// 1. 将请求体JSON绑定到指定的结构体
// 2. 使用validator库进行字段验证
// 3. 如果验证失败，返回详细的错误信息
// 4. 验证通过，将验证后的结构体存储到gin.Context中
func ValidateJSON(req interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建请求结构体的副本
		reqType := reflect.TypeOf(req)
		if reqType.Kind() == reflect.Ptr {
			reqType = reqType.Elem()
		}

		reqValue := reflect.New(reqType).Interface()

		// 绑定JSON到结构体
		if err := c.ShouldBindJSON(reqValue); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    http.StatusBadRequest,
				"message": "无效的JSON格式: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 执行验证
		if err := validate.Struct(reqValue); err != nil {
			var validationErrors []ValidationError

			// 处理验证错误
			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				for _, e := range validationErrs {
					validationErrors = append(validationErrors, ValidationError{
						Field:   e.Field(),
						Tag:     e.Tag(),
						Message: getErrorMessage(e),
					})
				}
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    http.StatusBadRequest,
				"message": "输入验证失败",
				"errors":  validationErrors,
			})
			c.Abort()
			return
		}

		// 将验证后的数据存储到Context中
		c.Set("validated_request", reqValue)
		c.Next()
	}
}

// GetValidatedRequest 从Context中获取验证后的请求
func GetValidatedRequest(c *gin.Context, target interface{}) bool {
	if req, exists := c.Get("validated_request"); exists {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(req).Elem())
		return true
	}
	return false
}

// getErrorMessage 将验证标签转换为中文错误消息
func getErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()

	switch tag {
	case "required":
		return field + " 是必填字段"
	case "min":
		return field + " 长度不能少于 " + fe.Param() + " 个字符"
	case "max":
		return field + " 长度不能超过 " + fe.Param() + " 个字符"
	case "email":
		return field + " 必须是有效的邮箱地址"
	case "oneof":
		return field + " 必须是以下值之一: " + fe.Param()
	case "numeric":
		return field + " 必须是数字"
	case "alpha":
		return field + " 只能包含字母"
	case "alphanum":
		return field + " 只能包含字母和数字"
	default:
		return field + " 格式不正确"
	}
}

// ValidateQuery 验证查询参数
func ValidateQuery(rules map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var errors []ValidationError

		for field, rule := range rules {
			value := c.Query(field)

			// 检查必填字段
			if strings.Contains(rule, "required") && value == "" {
				errors = append(errors, ValidationError{
					Field:   field,
					Tag:     "required",
					Message: field + " 是必填的查询参数",
				})
				continue
			}

			// 如果字段为空且不是必填，跳过其他验证
			if value == "" {
				continue
			}

			// 检查长度限制
			if strings.Contains(rule, "min:") {
				minLen := 0
				if n, err := fmt.Sscanf(strings.TrimPrefix(rule, "min:"), "%d", &minLen); err == nil && n == 1 {
					if len(value) < minLen {
						errors = append(errors, ValidationError{
							Field:   field,
							Tag:     "min",
							Message: field + " 长度不能少于 " + fmt.Sprintf("%d", minLen) + " 个字符",
						})
					}
				}
			}
		}

		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    http.StatusBadRequest,
				"message": "查询参数验证失败",
				"errors":  errors,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 自定义验证函数
func init() {
	// 注册自定义验证规则
	_ = validate.RegisterValidation("phone", validatePhone)
	_ = validate.RegisterValidation("password", validatePassword)
}

// validatePhone 验证手机号格式
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// 简单的中国手机号验证：1开头，第二位3-9，总共11位数字
	if len(phone) != 11 {
		return false
	}
	if phone[0] != '1' {
		return false
	}
	if phone[1] < '3' || phone[1] > '9' {
		return false
	}
	for _, c := range phone {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// validatePassword 验证密码强度
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	// 检查是否包含至少一个数字
	hasNumber := false
	// 检查是否包含至少一个字母
	hasLetter := false

	for _, c := range password {
		if c >= '0' && c <= '9' {
			hasNumber = true
		}
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasLetter = true
		}
	}

	return hasNumber && hasLetter
}
