package http_comm

// QueryUserByNameReq 查询用户名字
type QueryUserByNameReq struct {
	UserId uint64 `json:"user_id"`
}

// QueryUserByNameRsq 查询用户名字
type QueryUserByNameRsq struct {
	UserName string `json:"user_name"`
}
