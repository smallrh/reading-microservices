package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
}

const (
	SUCCESS              = 0
	ERROR                = 1
	ERROR_INVALID_PARAMS = 400
	ERROR_UNAUTHORIZED   = 401
	ERROR_FORBIDDEN      = 403
	ERROR_NOT_FOUND      = 404
	ERROR_INTERNAL       = 500
)

var CodeMsg = map[int]string{
	SUCCESS:              "success",
	ERROR:                "error",
	ERROR_INVALID_PARAMS: "invalid params",
	ERROR_UNAUTHORIZED:   "unauthorized",
	ERROR_FORBIDDEN:      "forbidden",
	ERROR_NOT_FOUND:      "not found",
	ERROR_INTERNAL:       "internal server error",
}

func GetMsg(code int) string {
	if msg, ok := CodeMsg[code]; ok {
		return msg
	}
	return CodeMsg[ERROR]
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    SUCCESS,
		Message: GetMsg(SUCCESS),
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    SUCCESS,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

func ErrorWithCode(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: GetMsg(code),
	})
}

func PageSuccess(c *gin.Context, data interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, PageResponse{
		Code:    SUCCESS,
		Message: GetMsg(SUCCESS),
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}
