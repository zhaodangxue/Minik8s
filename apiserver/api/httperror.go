package api

type HttpError struct {
	Code    ApiserverErrorCode `json:"code"`
	Message string             `json:"message"`
}

type ApiserverErrorCode int32

const (
	// 0 为非法错误码，不应该出现在实际的错误码中
	ApiserverErrorCode_NO_ERROR ApiserverErrorCode = 0
	// 1 是通用错误码
	ApiserverErrorCode_GENERIC_ERROR ApiserverErrorCode = 1
	// 2 是请求错误码，代表请求的参数不合法
	ApiserverErrorCode_INVALID_REQUEST ApiserverErrorCode = 2
	// 大于1000的错误码为自定义错误码

	// 1101 代表Update Pod对应的Pod不存在
	ApiserverErrorCode_UPDATE_POD_NOT_FOUND ApiserverErrorCode = 1101
)
