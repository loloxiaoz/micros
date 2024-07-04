package common

type BizError struct {
	//错误码
	Code int
	//错误信息
	Message string
}

func (err BizError) Error() string {
	return err.Message
}

/*
错误码设计
第一位表示错误级别, 0 为未知错误 1 为系统错误, 2 为业务错误
第二三位表示服务模块代码
第四五位表示具体错误代码
*/

var (
	OK = &BizError{Code: 0, Message: "OK"}

	// 通用服务端错误
	CommonServerError = &BizError{Code: 00010, Message: "服务端错误"}
	BindJSONError     = &BizError{Code: 00020, Message: "参数解析错误"}

	// 系统错误
	InternalServerError = &BizError{Code: 10001, Message: "内部服务器错误"}
	FileSaveError       = &BizError{Code: 10101, Message: "保存文件错误"}
	FileUploadError     = &BizError{Code: 10102, Message: "上传文件错误"}

	// 业务错误
	NameEmptyError  = &BizError{Code: 20401, Message: "名称不能为空"}
	NameExistsError = &BizError{Code: 20402, Message: "名称已存在"}
)
