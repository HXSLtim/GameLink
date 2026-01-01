// XSS 防护示例：聊天和评价
// 此文件展示如何在聊天和评价功能中添加 XSS 防护

package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/pkg/sanitize"
	orderservice "gamelink/internal/service/order"
)

// createReviewHandlerWithXSS 创建评价（带 XSS 防护）
// @Summary      创建评价
// @Description  为已完成订单创建评价（包含 XSS 防护）
// @Tags         User - Reviews
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                        true  "Bearer {token}"
// @Param        request        body      orderservice.CreateReviewRequest    true  "创建评价请求"
// @Success      200            {object}  model.APIResponse[CreateReviewResponse]
// @Failure      400            {object}  model.ErrorResponse
// @Failure      401            {object}  model.ErrorResponse
// @Router       /user/reviews [post]
func createReviewHandlerWithXSS(c *gin.Context, svc *orderservice.ReviewService) {
	userID := getUserIDFromContext(c)

	var req orderservice.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// ==================== XSS 防护 ====================
	// 1. 清理评价内容（移除危险模式、验证 UTF-8、限制长度）
	req.Comment = sanitize.SanitizeReview(req.Comment)

	// 2. 清理评价标签数组中的每个标签
	for i, tag := range req.Tags {
		req.Tags[i] = sanitize.SanitizeNickname(tag) // 标签使用昵称规则
	}
	// ==================================================

	resp, err := svc.CreateReview(c.Request.Context(), userID, req)
	if err != nil {
		if err == orderservice.ErrAlreadyReviewed {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		if err == orderservice.ErrOrderNotCompleted {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		if err == orderservice.ErrReviewUnauthorized {
			respondError(c, http.StatusForbidden, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, model.APIResponse[orderservice.CreateReviewResponse]{
		Success: true,
		Code:    http.StatusOK,
		Message: "评价创建成功",
		Data:    *resp,
	})
}

// ====================================================================
// 以下展示如何集成到现有的 review.go 文件中
// ====================================================================

/*
集成步骤：

1. 在 review.go 文件顶部添加导入：
   import (
       "gamelink/pkg/sanitize"  // 添加此行
       ... 其他导入
   )

2. 在 createReviewHandler 函数中，找到这行代码：
   if err := c.ShouldBindJSON(&req); err != nil {

3. 在该代码块之后，添加以下代码：
   // XSS 防护
   req.Comment = sanitize.SanitizeReview(req.Comment)
   for i, tag := range req.Tags {
       req.Tags[i] = sanitize.SanitizeNickname(tag)
   }

4. 完成！
*/

// ====================================================================
// 聊天功能的 XSS 防护示例
// ====================================================================

/*
聊天功能需要在 WebSocket 处理器中添加 XSS 防护。

假设你在 internal/handler/user/chat.go 或类似的文件中有 WebSocket 消息处理：

import (
    "gamelink/pkg/sanitize"  // 添加导入
)

// 在发送消息的处理函数中添加：
func handleSendMessage(c *gin.Context, ws *websocket.Conn) {
    var msg struct {
        Content string `json:"content"`
        RoomID  uint64 `json:"roomId"`
    }

    if err := c.ShouldBindJSON(&msg); err != nil {
        // 错误处理
        return
    }

    // ========== 新增：XSS 防护 ==========
    // 清理聊天消息内容
    msg.Content = sanitize.SanitizeMessage(msg.Content)
    // ================================

    // 发送清理后的消息
    // ...
}

// 如果是 WebSocket 直接接收的 JSON 消息：
func handleWSMessage(message []byte) {
    var msg struct {
        Content string `json:"content"`
    }

    json.Unmarshal(message, &msg)

    // ========== 新增：XSS 防护 ==========
    msg.Content = sanitize.SanitizeMessage(msg.Content)
    // ================================

    // 处理消息...
}
*/

// ====================================================================
// 举报功能的 XSS 防护示例
// ====================================================================

/*
举报功能通常在 dispute.go 或类似文件中。

import (
    "gamelink/pkg/sanitize"  // 添加导入
)

func createDisputeHandler(c *gin.Context, svc *DisputeService) {
    var req CreateDisputeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // 错误处理
        return
    }

    // ========== 新增：XSS 防护 ==========
    // 清理举报理由和详细说明
    req.Reason = sanitize.SanitizeNickname(req.Reason)
    req.Description = sanitize.SanitizeReport(req.Description)
    // ================================

    // 创建举报...
}
*/

// ====================================================================
// 动态/Feed 功能的 XSS 防护示例
// ====================================================================

/*
动态功能在 feed.go 中。

import (
    "gamelink/pkg/sanitize"  // 添加导入
)

func createFeedHandler(c *gin.Context, svc *FeedService) {
    var req CreateFeedRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // 错误处理
        return
    }

    // ========== 新增：XSS 防护 ==========
    // 清理动态内容
    req.Content = sanitize.SanitizeMessage(req.Content)
    // 清理图片/视频标题（如果有）
    for i, media := range req.Media {
        if media.Title != "" {
            req.Media[i].Title = sanitize.SanitizeNickname(media.Title)
        }
    }
    // 清理提及的用户（如果有）
    for i, mention := range req.Mentions {
        req.Mentions[i].DisplayName = sanitize.SanitizeNickname(mention.DisplayName)
    }
    // ================================

    // 创建动态...
}
*/

// ====================================================================
// 全局中间件方式（可选）
// ====================================================================

/*
如果你想要对所有请求体自动进行 XSS 清理，可以创建一个中间件：

package middleware

import (
    "github.com/gin-gonic/gin"
    "gamelink/pkg/sanitize"
)

// XSSProtectionMiddleware 自动清理请求体中的字符串字段
func XSSProtectionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 只处理 POST/PUT/PATCH 请求
        if c.Request.Method != "POST" && c.Request.Method != "PUT" && c.Request.Method != "PATCH" {
            c.Next()
            return
        }

        // 读取请求体
        body, err := c.GetRawData()
        if err != nil {
            c.Next()
            return
        }

        // 解析为 map
        var data map[string]interface{}
        if err := json.Unmarshal(body, &data); err != nil {
            c.Next()
            return
        }

        // 清理所有字符串字段
        cleanedData := sanitize.EscapeAll(data)

        // 将清理后的数据重新设置到上下文
        c.Set("cleaned_body", cleanedData)

        c.Next()
    }
}

然后在路由中使用：
router.Use(middleware.XSSProtectionMiddleware())

注意：这种方式是可选的，建议在每个 handler 中明确地进行清理，以获得更好的控制。
*/
