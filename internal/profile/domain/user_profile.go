package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"time"
)

const (
	maxNameLen    = 32
	maxAddressLen = 100
	maxSignLen    = 100
)

type UserProfile struct {
	id          kernel.UserID
	name        string
	gender      kernel.Gender
	email       kernel.Email
	phoneNumber kernel.PhoneNumber
	address     kernel.Address
	sign        string
	signedUpAt  time.Time
}

func LoadUserProfile(
	id kernel.UserID,
	name string,
	gender kernel.Gender,
	email kernel.Email,
	phoneNumber kernel.PhoneNumber,
	address kernel.Address,
	sign string,
	signedUpAt time.Time,
) *UserProfile {
	return &UserProfile{
		id:          id,
		name:        name,
		gender:      gender,
		email:       email,
		phoneNumber: phoneNumber,
		address:     address,
		sign:        sign,
		signedUpAt:  signedUpAt,
	}
}

func CreateUserProfile(
	id kernel.UserID,
	email kernel.Email,
	signedUpAt time.Time,
) *UserProfile {
	return &UserProfile{
		id:         id,
		gender:     kernel.UnknownGender,
		email:      email,
		signedUpAt: signedUpAt,
	}
}

func (p *UserProfile) UpdateName(name string) error {
	if len(name) > maxNameLen {
		return myErrors.NewBusiness("name length exceeds limit: max=%d", maxNameLen)
	}
	p.name = name
	return nil
}

func (p *UserProfile) UpdateGender(gender kernel.Gender) {
	p.gender = gender
}

func (p *UserProfile) UpdatePhoneNumber(phoneNumber kernel.PhoneNumber) {
	p.phoneNumber = phoneNumber
}

func (p *UserProfile) UpdateEmail(email kernel.Email) {
	p.email = email
}

func (p *UserProfile) UpdateAddress(address kernel.Address) error {
	if len(address) > maxAddressLen {
		return myErrors.NewBusiness("address length exceeds limit: max=%d", maxAddressLen)
	}
	p.address = address
	return nil
}

func (p *UserProfile) UpdateSign(sign string) error {
	if len(sign) > maxSignLen {
		return myErrors.NewBusiness("sign length exceeds limit: max=%d", maxSignLen)
	}
	p.sign = sign
	return nil
}

func (p *UserProfile) ID() kernel.UserID {
	return p.id
}

func (p *UserProfile) Name() string {
	return p.name
}

func (p *UserProfile) Gender() kernel.Gender {
	return p.gender
}

func (p *UserProfile) Email() kernel.Email {
	return p.email
}

func (p *UserProfile) PhoneNumber() kernel.PhoneNumber {
	return p.phoneNumber
}

func (p *UserProfile) Address() kernel.Address {
	return p.address
}

func (p *UserProfile) Sign() string {
	return p.sign
}

func (p *UserProfile) SignedUpAt() time.Time {
	return p.signedUpAt
}
