package auth

type SendOTPRequest struct {
	Phone string `json:"phone"`
}

type SendOTPResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"` // dev only
}

type VerifyOTPRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
	Name  string `json:"name,omitempty"`
}

type VerifyOTPResponse struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
}
