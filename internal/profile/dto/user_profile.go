package dto

import (
	"gochat/internal/profile/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type UserProfile struct {
	ID          kernel.UserID      `json:"id" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	Name        string             `json:"name" example:"Jack"`
	Gender      kernel.Gender      `json:"gender" example:"1"`
	Email       kernel.Email       `json:"email" example:"email@demo.com"`
	PhoneNumber kernel.PhoneNumber `json:"phone_number" example:"13310001000"`
	Address     kernel.Address     `json:"address" example:"China"`
	Sign        string             `json:"sign" example:"I'm Jack!'"`
	SignedUpAt  time.Time          `json:"signed_up_at" example:"2025-12-26 05:38:19.740"`
}

func ToUserProfileDTO(profile *domain.UserProfile) *UserProfile {
	return &UserProfile{
		ID:          profile.ID(),
		Name:        profile.Name(),
		Gender:      profile.Gender(),
		Email:       profile.Email(),
		PhoneNumber: profile.PhoneNumber(),
		Address:     profile.Address(),
		Sign:        profile.Sign(),
		SignedUpAt:  profile.SignedUpAt(),
	}
}
