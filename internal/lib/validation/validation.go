package validation

import (
	"github.com/go-playground/validator/v10"
	ssov1 "github.com/whitxowl/contracts/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Validation struct {
	validate *validator.Validate
}

func NewValidator() *Validation {
	v := validator.New()
	return &Validation{validate: v}
}

type RegisterRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
	AppID    int64  `validate:"required,min=1"`
}

type IsAdminRequest struct {
	UserId int64 `validate:"required,min=1"`
}

func (v *Validation) ValidateRegisterRequest(req *ssov1.RegisterRequest) error {
	registerReq := RegisterRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	if err := v.validate.Struct(registerReq); err != nil {
		return status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	return nil
}

func (v *Validation) ValidateLoginRequest(req *ssov1.LoginRequest) error {
	loginReq := LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		AppID:    req.GetAppId(),
	}

	if err := v.validate.Struct(loginReq); err != nil {
		return status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	return nil
}

func (v *Validation) ValidateIsAdminRequest(req *ssov1.IsAdminRequest) error {
	isAdminReq := IsAdminRequest{
		UserId: req.GetUserId(),
	}

	if err := v.validate.Struct(isAdminReq); err != nil {
		return status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	return nil
}
