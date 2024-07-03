package common

const (
	TimeFormat    = "2006-01-02 15:04:05"
	DateFormat    = "2006-01-02"
	TimeTplDay    = "2006-01-02 00:00:00"
	TimeTplHour   = "2006-01-02 15:00:00"
	TimeTplMinute = "2006-01-02 15:04:00"

	ErrorInvalidArgument    = "INVALID_ARGUMENT"    //客户端指定了无效参数
	ErrorFailedPrecondition = "FAILED_PRECONDITION" //请求无法在当前系统状态下执行，例如删除非空目录
	ErrorOutOfRange         = "OUT_OF_RANGE"        //客户端指定了无效范围
	ErrorUnauthenticated    = "UNAUTHENTICATED"     //由于 OAuth 令牌丢失、无效或过期，请求未通过身份验证
	ErrorPermissionDenied   = "PERMISSION_DENIED"   //客户端权限不足。可能的原因包括 OAuth 令牌的覆盖范围不正确、客户端没有权限或者尚未为客户端项目启用 API。
	ErrorNotFound           = "NOT_FOUND"           //找不到指定的资源，或者请求由于未公开的原因（例如白名单）而被拒绝
	ErrorAborted            = "ABORTED"             //并发冲突，例如读取/修改/写入冲突
	ErrorAlreadyExists      = "ALREADY_EXISTS"      //客户端尝试创建的资源已存在
	ErrorResourceExhausted  = "RESOURCE_EXHAUSTED"  //资源配额不足或达到速率限制
	ErrorCancelled          = "CANCELLED"           //请求被客户端取消
	ErrorDataLoss           = "DATA_LOSS"           //出现不可恢复的数据丢失或数据损坏
	ErrorUnknown            = "UNKNOWN"             //出现未知的服务器错误。通常是服务器错误
	ErrorInternal           = "INTERNAL"            //出现未知的服务器错误。通常是服务器错误。
	ErrorNotImplemented     = "NOT_IMPLEMENTED"     //API 方法未通过服务器实现。
	ErrorUnavailable        = "UNAVAILABLE"         //服务不可用。通常是服务器已关闭。
	ErrorDeadlineExceeded   = "DEADLINE_EXCEEDED"   //超时
)


//ErrCodeMap 错误类型 和 错误码 对应关系
var ErrCodeMap = map[string]int{
	ErrorInvalidArgument:    400, //客户端指定了无效参数
	ErrorFailedPrecondition: 400, //请求无法在当前系统状态下执行，例如删除非空目录
	ErrorOutOfRange:         400, //客户端指定了无效范围
	ErrorUnauthenticated:    401, //由于 OAuth 令牌丢失、无效或过期，请求未通过身份验证
	ErrorPermissionDenied:   403, //客户端权限不足。可能的原因包括 OAuth 令牌的覆盖范围不正确、客户端没有权限或者尚未为客户端项目启用 API。
	ErrorNotFound:           404, //找不到指定的资源，或者请求由于未公开的原因（例如白名单）而被拒绝
	ErrorAborted:            409, //并发冲突，例如读取/修改/写入冲突
	ErrorAlreadyExists:      409, //客户端尝试创建的资源已存在
	ErrorResourceExhausted:  429, //资源配额不足或达到速率限制
	ErrorCancelled:          499, //请求被客户端取消
	ErrorDataLoss:           500, //出现不可恢复的数据丢失或数据损坏
	ErrorUnknown:            500, //出现未知的服务器错误。通常是服务器错误
	ErrorInternal:           500, //出现未知的服务器错误。通常是服务器错误。
	ErrorNotImplemented:     501, //API 方法未通过服务器实现。
	ErrorUnavailable:        503, //服务不可用。通常是服务器已关闭。
	ErrorDeadlineExceeded:   504, //超时
}
