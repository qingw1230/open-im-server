package constant

import "errors"

// key = errCode, string = errMsg
type ErrInfo struct {
	ErrCode int32
	ErrMsg  string
}

var (
	OK                   = ErrInfo{0, ""}
	ErrServer            = ErrInfo{500, "server error"}
	ErrParseToken        = ErrInfo{700, ParseTokenMsg.Error()}
	ErrTencentCredential = ErrInfo{400, ThirdPartyMsg.Error()}

	ErrTokenExpired     = ErrInfo{701, TokenExpiredMsg.Error()}
	ErrTokenInvalid     = ErrInfo{702, TokenInvalidMsg.Error()}
	ErrTokenMalformed   = ErrInfo{703, TokenMalformedMsg.Error()}
	ErrTokenNotValidYet = ErrInfo{704, TokenNotValidYetMsg.Error()}
	ErrTokenUnknown     = ErrInfo{705, TokenUnknownMsg.Error()}

	ErrAccess = ErrInfo{ErrCode: 801, ErrMsg: AccessMsg.Error()}
	ErrDB     = ErrInfo{ErrCode: 802, ErrMsg: DBMsg.Error()}
	ErrArgs   = ErrInfo{ErrCode: 8003, ErrMsg: ArgsMsg.Error()}
)

var (
	ParseTokenMsg       = errors.New("parse token failed")
	TokenExpiredMsg     = errors.New("token is timed out, please log in again")
	TokenInvalidMsg     = errors.New("token has been invalidated")
	TokenNotValidYetMsg = errors.New("token not active yet")
	TokenMalformedMsg   = errors.New("that's not even a token")
	TokenUnknownMsg     = errors.New("couldn't handle this token")

	AccessMsg = errors.New("no permission")
	DBMsg     = errors.New("db failed")
	ArgsMsg   = errors.New("args failed")

	ThirdPartyMsg = errors.New("third party error")
)

const (
	NoError              = 0
	FormattingError      = 10001
	HasRegistered        = 10002
	NotRegistered        = 10003
	PasswordErr          = 10004
	GetIMTokenErr        = 10005
	RepeatSendCode       = 10006
	MailSendCodeErr      = 10007
	SmsSendCodeErr       = 10008
	CodeInvalidOrExpired = 10009
	RegisterFailed       = 10010
	ResetPasswordFailed  = 10011
	DatabaseError        = 10002
	ServerError          = 10004
	HttpError            = 10005
	IoError              = 10006
	IntentionalError     = 10007
)

func (e ErrInfo) Error() string {
	return e.ErrMsg
}

func (e *ErrInfo) Code() int32 {
	return e.ErrCode
}
