// XSS 防护示例实现
// 此文件展示如何在现有的 profile.go 中添加 XSS 防护

package player

import (
	"github.com/gin-gonic/gin"

	"gamelink/pkg/apierr"
	"gamelink/pkg/sanitize"
	serviceplayer "gamelink/internal/service/player"
)

// updatePlayerProfileHandlerWithXSS 更新陪玩师资料（带 XSS 防护）
// @Summary      更新陪玩师资料
// @Description  更新陪玩师个人资料（包含 XSS 防护）
// @Tags         Player - Profile
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  UpdatePlayerProfileRequest  true  "更新信息"
// @Success      200      {object}  model.SuccessResponse
// @Failure      400      {object}  apierr.APIError
// @Failure      401      {object}  apierr.APIError
// @Failure      404      {object}  apierr.APIError
// @Failure      422      {object}  apierr.APIError
// @Failure      500      {object}  apierr.APIError
// @Router       /player/profile [put]
func updatePlayerProfileHandlerWithXSS(c *gin.Context, svc *serviceplayer.PlayerService) {
	userID := getUserIDFromContext(c)

	var req serviceplayer.UpdatePlayerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	// ==================== XSS 防护 ====================
	// 1. 清理昵称（限制长度、移除 HTML、转义特殊字符）
	req.Nickname = sanitize.SanitizeNickname(req.Nickname)

	// 2. 清理个人简介（移除危险模式、验证 UTF-8、限制长度）
	req.Bio = sanitize.SanitizeMessage(req.Bio)

	// 3. 清理标签数组中的每个标签
	for i, tag := range req.Tags {
		req.Tags[i] = sanitize.SanitizeNickname(tag) // 标签使用昵称规则
	}
	// ==================================================

	if err := svc.UpdatePlayerProfile(c.Request.Context(), userID, req); err != nil {
		if err == serviceplayer.ErrNotFound {
			respondAPIError(c, apierr.NotFound(err.Error()))
			return
		}
		if err == serviceplayer.ErrValidation {
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("更新资料失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "资料更新成功", struct{}{})
}

// ====================================================================
// 以下展示如何集成到现有的 profile.go 文件中
// ====================================================================

/*
集成步骤：

1. 在 profile.go 文件顶部添加导入：
   import (
       "gamelink/pkg/sanitize"  // 添加此行
       ... 其他导入
   )

2. 在 updatePlayerProfileHandler 函数中，找到这行代码：
   if err := c.ShouldBindJSON(&req); err != nil {

3. 在该代码块之后，添加以下代码：
   // XSS 防护
   req.Nickname = sanitize.SanitizeNickname(req.Nickname)
   req.Bio = sanitize.SanitizeMessage(req.Bio)
   for i, tag := range req.Tags {
       req.Tags[i] = sanitize.SanitizeNickname(tag)
   }

4. 完成！
*/

// ====================================================================
// 完整的修改示例（替换 profile.go 中的 updatePlayerProfileHandler 函数）
// ====================================================================

/*
func updatePlayerProfileHandler(c *gin.Context, svc *serviceplayer.PlayerService) {
	userID := getUserIDFromContext(c)

	var req serviceplayer.UpdatePlayerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	// ========== 新增：XSS 防护 ==========
	// 清理用户输入，防止 XSS 攻击
	req.Nickname = sanitize.SanitizeNickname(req.Nickname)
	req.Bio = sanitize.SanitizeMessage(req.Bio)
	for i, tag := range req.Tags {
		req.Tags[i] = sanitize.SanitizeNickname(tag)
	}
	// ================================

	if err := svc.UpdatePlayerProfile(c.Request.Context(), userID, req); err != nil {
		if err == serviceplayer.ErrNotFound {
			respondAPIError(c, apierr.NotFound(err.Error()))
			return
		}
		if err == serviceplayer.ErrValidation {
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("更新资料失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "资料更新成功", struct{}{})
}
*/
