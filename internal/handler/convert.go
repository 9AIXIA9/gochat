package handler

import (
	"gochat/internal/domain"
	"time"
)

func (r SignupRequest) ToDomain() *domain.SignupRequest {
	return &domain.SignupRequest{
		Password: r.Body.Password,
	}
}

func (r LoginRequest) ToDomain() *domain.LoginRequest {
	return &domain.LoginRequest{
		Number:   r.Body.Number,
		Password: r.Body.Password,
	}
}

func (r RefreshRequest) ToDomain() *domain.RefreshTokenRequest {
	return &domain.RefreshTokenRequest{
		RefreshToken: r.Cookie.RefreshToken,
	}
}

func (r CreateRoomRequest) ToDomain() *domain.CreateRoomRequest {
	return &domain.CreateRoomRequest{
		AuthInfo: r.AuthInfo,
		Secret:   r.Body.Secret,
		MaxUsers: r.Body.MaxUsers,
	}
}

func (r JoinRoomRequest) ToDomain() *domain.JoinRoomRequest {
	return &domain.JoinRoomRequest{
		AuthInfo: r.AuthInfo,
		Number:   r.URI.Number,
		Secret:   r.Body.Secret,
	}
}

func (r LeaveRoomRequest) ToDomain() *domain.LeaveRoomRequest {
	return &domain.LeaveRoomRequest{
		AuthInfo: r.AuthInfo,
		Number:   r.URI.RoomNumber,
	}
}

func (r SendMessageRequest) ToDomain() *domain.SendMessageRequest {
	return &domain.SendMessageRequest{
		AuthInfo: r.AuthInfo,
		To:       r.URI.To,
		Content:  r.Body.Content,
		SentAt:   time.Unix(r.Body.SentAt, 0),
	}
}
