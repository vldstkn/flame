package dto

type AccountRegisterReq struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=50"`
	Location string `json:"location,omitempty"`
}

type AccountRegisterRes struct {
	AccessToken string `json:"access_token"`
}

type AccountLoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=50"`
	Location string `json:"location,omitempty"`
}

type AccountLoginRes struct {
	AccessToken string `json:"access_token"`
}

type AccountGetTokensRes struct {
	AccessToken string `json:"access_token"`
}

type AccountUpdateProfileReq struct {
	Name      *string `json:"name,omitempty"`
	BirthDate *string `json:"birth_date,omitempty"`
	City      *string `json:"city,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	Gender    *string `json:"gender,omitempty"`
}

type DeletePhotoReq struct {
	PhotoId int64 `json:"photo_id" validate:"required,min=0,numeric"`
}

type GetMatchingReq struct {
	Location string `json:"location" validate:"required"`
}

type UpdateLocation struct {
	Location string `json:"location" validate:"required"`
}
