package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"reading-microservices/shared/utils"
	"reading-microservices/user-service/models"
	"reading-microservices/user-service/services"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "注册信息"
// @Success 200 {object} utils.Response{data=models.LoginResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	response, err := h.userService.Register(&req)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, response)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登录信息"
// @Success 200 {object} utils.Response{data=models.LoginResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	// 获取客户端IP
	req.Platform = getClientPlatform(c)

	response, err := h.userService.Login(&req)
	if err != nil {
		utils.Error(c, utils.ERROR_UNAUTHORIZED, err.Error())
		return
	}

	utils.Success(c, response)
}

// GetProfile 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的个人信息
// @Tags 用户信息
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response{data=models.UserInfo}
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		utils.Error(c, utils.ERROR_NOT_FOUND, err.Error())
		return
	}

	utils.Success(c, profile)
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户的个人信息
// @Tags 用户信息
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.UpdateProfileRequest true "更新信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	if err := h.userService.UpdateProfile(userID, &req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Profile updated successfully", nil)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 用户信息
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.ChangePasswordRequest true "密码信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	if err := h.userService.ChangePassword(userID, &req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Password changed successfully", nil)
}

// RefreshToken 刷新token
// @Summary 刷新Token
// @Description 使用旧Token刷新获取新Token
// @Tags 用户认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response{data=object{token=string,expires_in=int}}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	token := extractTokenFromHeader(c)
	if token == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	newToken, err := h.userService.RefreshToken(token)
	if err != nil {
		utils.Error(c, utils.ERROR_UNAUTHORIZED, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"token":      newToken,
		"expires_in": 86400,
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口，使当前Token失效
// @Tags 用户认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /user/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	token := extractTokenFromHeader(c)
	if token == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	if err := h.userService.Logout(token); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Logout successfully", nil)
}

// ValidateToken 验证token
// @Summary 验证Token
// @Description 验证Token的有效性
// @Tags 用户认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response{data=object{user_id=string,username=string,valid=bool}}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/validate [post]
func (h *UserHandler) ValidateToken(c *gin.Context) {
	token := extractTokenFromHeader(c)
	if token == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	claims, err := h.userService.ValidateToken(token)
	if err != nil {
		utils.Error(c, utils.ERROR_UNAUTHORIZED, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"user_id": claims.UserID,
		"valid":   true,
	})
}

// Health 健康检查
// @Summary 健康检查
// @Description 服务健康检查接口
// @Tags 系统管理
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=object{status=string,service=string}}
// @Router /health [get]
func (h *UserHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "user-service",
		"time":    gin.H{},
	})
}

// extractTokenFromHeader 从请求头中提取token
func extractTokenFromHeader(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token == "" {
		return ""
	}

	// Remove Bearer prefix
	if len(token) > 7 && strings.ToLower(token[:7]) == "bearer " {
		token = token[7:]
	}
	return token
}

// getClientPlatform 从请求中获取客户端平台信息
func getClientPlatform(c *gin.Context) string {
	platform := c.GetHeader("X-Platform")
	if platform == "" {
		userAgent := c.GetHeader("User-Agent")
		if userAgent != "" {
			ua := strings.ToLower(userAgent)
			if strings.Contains(ua, "android") {
				return "android"
			}
			if strings.Contains(ua, "iphone") || strings.Contains(ua, "ios") {
				return "ios"
			}
			return "web"
		}
		return "h5"
	}
	return platform
}
