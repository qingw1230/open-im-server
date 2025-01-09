package base_info

type UserRegisterReq struct {
	Secret   string `json:"secret" binding:"required,max=32"`
	Platform int32  `json:"platform" binding:"required,min=1,max=7"`
	ApiUserInfo
	OperationID string `json:"operationID" binding:"required"`
}

type UserRegisterResp struct {
	CommResp
	UserToken UserTokenInfo `json:"data"`
}

type UserTokenInfo struct {
	UserID      string `json:"userID"`
	Token       string `json:"token"`
	ExpiredTime int64  `json:"expiredTime"`
}

type UserTokenReq struct {
	Secret      string `json:"secret" binding:"required,max=32"`
	Platform    int32  `json:"platform" binding:"required,min=1,max=8"`
	UserID      string `json:"userID" binding:"required,min=1,max=64"`
	OperationID string `json:"operationID" binding:"required"`
}

type UserTokenResp struct {
	CommResp
	UserToken UserTokenInfo `json:"data"`
}
