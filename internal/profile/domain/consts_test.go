package domain_test

import (
	"gochat/internal/profile/domain"
	"gochat/internal/shared/kernel"
)

const (
	fixedUserID       kernel.UserID      = "user-123"
	fixedRoomID       kernel.RoomID      = "room-456"
	fixedRoomshipID   domain.RoomshipID  = "roomship-789"
	fixedName                            = "Test Name"
	fixedIntroduction                    = "This is a test introduction."
	fixedEmail        kernel.Email       = "demo@test.com"
	fixedPhoneNumber  kernel.PhoneNumber = "123-456-7890"
	fixedSign                            = "This is a test sign."
	fixedGender                          = kernel.MaleGender
	fixedCountry                         = "Test Country"
	fixedProvince                        = "Test Province"
	fixedCity                            = "Test City"
	fixedDistrict                        = "Test District"
	fixedStreet                          = "123 Test St."
)
