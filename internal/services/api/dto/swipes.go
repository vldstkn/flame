package dto

type CreateSwipesReq struct {
	UserId int64 `json:"user_id" validate:"required,number"`
	IsLike *bool `json:"is_like" validate:"required,boolean"`
}

type GetUnreadSwipes struct {
	Users []int64 `json:"users"`
}
