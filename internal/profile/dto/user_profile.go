package dto

import (
	"gochat/internal/profile/domain"
	"gochat/internal/shared/kernel"
	"time"
)

type UserProfile struct {
	ID          kernel.UserID      `json:"id"`
	Name        string             `json:"name"`
	Gender      kernel.Gender      `json:"gender"`
	Email       kernel.Email       `json:"email"`
	PhoneNumber kernel.PhoneNumber `json:"phone_number"`
	Address     kernel.Address     `json:"address"`
	Sign        string             `json:"sign"`
	SignedUpAt  time.Time          `json:"signed_up_at"`
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
