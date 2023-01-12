package response

import (
	"github.com/gin-gonic/gin"
	"github.com/wenccc/myskeleton/logger"
	"gorm.io/gorm"
	"net/http"
)

type responseData struct {
	Code    Code        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Code int

const (
	CodeSuccess          Code = 0
	CodeFail             Code = 1
	CodeNotFound         Code = 2
	CodeRequestParseFail Code = 3
	CodeServerError      Code = 4
	CodeValidateError    Code = 5
	CodeForbidden        Code = 6
	CodeUnauthorized     Code = 7
)

func defaultMessage(df string, msg ...string) string {
	if len(msg) > 0 {
		return msg[0]
	}
	return df
}

func Json(ctx *gin.Context, code int, data interface{}) {
	ctx.JSON(code, data)
}

func SuccessWithoutData(c *gin.Context) {
	Json(c, http.StatusOK, responseData{
		Code:    CodeSuccess,
		Message: "操作成功",
		Data:    nil,
	})
}

func Success(c *gin.Context, data interface{}) {
	Json(c, http.StatusOK, data)
}

func FailWithoutData(c *gin.Context, msg ...string) {
	Json(c, http.StatusBadRequest, responseData{
		Code:    CodeFail,
		Message: defaultMessage("操作失败", msg...),
		Data:    nil,
	})
}

func Fail(c *gin.Context, data interface{}, msg ...string) {
	Json(c, http.StatusBadRequest, responseData{
		Code:    CodeFail,
		Message: defaultMessage("操作失败", msg...),
		Data:    data,
	})
}
func Abort404WithoutData(c *gin.Context, msg ...string) {

	c.AbortWithStatusJSON(http.StatusNotFound, responseData{
		Code:    CodeNotFound,
		Message: defaultMessage("操作失败", msg...),
		Data:    nil,
	})
}

func Abort404(c *gin.Context, data interface{}, msg ...string) {
	c.AbortWithStatusJSON(http.StatusNotFound, responseData{
		Code:    CodeNotFound,
		Message: defaultMessage("操作失败", msg...),
		Data:    data,
	})
}
func Abort403WithoutData(c *gin.Context, msg ...string) {

	c.AbortWithStatusJSON(http.StatusForbidden, responseData{
		Code:    CodeForbidden,
		Message: defaultMessage("权限不足，请确定您有对应的权限", msg...),
		Data:    nil,
	})
}

func Abort403(c *gin.Context, data interface{}, msg ...string) {
	c.AbortWithStatusJSON(http.StatusForbidden, responseData{
		Code:    CodeForbidden,
		Message: defaultMessage("权限不足，请确定您有对应的权限", msg...),
		Data:    data,
	})
}

func Abort500WithoutData(c *gin.Context, msg ...string) {

	c.AbortWithStatusJSON(http.StatusInternalServerError, responseData{
		Code:    CodeServerError,
		Message: defaultMessage("服务器内部错误，请稍后再试", msg...),
		Data:    nil,
	})
}

func Abort500(c *gin.Context, data interface{}, msg ...string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, responseData{
		Code:    CodeServerError,
		Message: defaultMessage("服务器内部错误，请稍后再试", msg...),
		Data:    data,
	})
}

func BadRequest(c *gin.Context, err error, msg ...string) {

	logger.LogErrorIf(err)
	c.AbortWithStatusJSON(http.StatusBadRequest, responseData{
		Code:    CodeRequestParseFail,
		Message: defaultMessage("请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。", msg...),
		Data:    nil,
	})
}

// Error 响应 404 或 422，未传参 msg 时使用默认消息
// 处理请求时出现错误 err，会附带返回 error 信息，如登录错误、找不到 ID 对应的 Model
func Error(c *gin.Context, err error, msg ...string) {
	logger.LogErrorIf(err)

	// error 类型为『数据库未找到内容』
	if err == gorm.ErrRecordNotFound {
		Abort404WithoutData(c, msg...)
		return
	}

	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, responseData{
		Code:    CodeFail,
		Message: defaultMessage(err.Error(), msg...),
		Data:    nil,
	})
}

func ValidationError(c *gin.Context, errors map[string][]string) {

	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, responseData{
		Code:    CodeValidateError,
		Message: "数据验证失败",
		Data:    errors,
	})
}

func Unauthorized(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, responseData{
		Code:    CodeUnauthorized,
		Message: "未认证",
		Data:    nil,
	})
}
