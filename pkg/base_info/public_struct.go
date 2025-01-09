package base_info

type CommResp struct {
	ErrCode int32  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type CommDataResp struct {
	CommResp
	Data []map[string]interface{} `json:"data"`
}

type ApiUserInfo struct {
	UserID      string `json:"userID" binding:"required,min=1,max=64"`
	Nickname    string `json:"nickname" binding:"omitempty,min=1,max=64"`
	FaceURL     string `json:"faceURL" binding:"omitempty,max=1024"`
	Gender      int32  `json:"gender" binding:"omitempty,oneof=0 1 2"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,max=32"`
	Birth       uint32 `json:"birth" binding:"omitempty"`
	Email       string `json:"email" binding:"omitempty,max=64"`
	Ex          string `json:"ex" binding:"omitempty,max=1024"`
}
