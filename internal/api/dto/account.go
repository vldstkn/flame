package dto

type AccountRegisterReq struct {
	Name       string `json:"name" validate:"required,min=2,max=50"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6,max=50"`
	Gender     string `json:"gender" validate:"required"`
	LookingFor string `json:"looking_for" validate:"required"`
}

type AccountRegisterRes struct {
	AccessToken string `json:"access_token"`
}

type AccountLoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

type AccountLoginRes struct {
	AccessToken string `json:"access_token"`
}

type GetTokensRes struct {
	AccessToken string `json:"access_token"`
}
