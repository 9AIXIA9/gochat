package domain

import (
	"gochat/internal/shared/kernel"
	kernelEvent "gochat/internal/shared/kernel/event"
	"time"
)

var _ kernelEvent.Getter = (*User)(nil)

type User struct {
	id             kernel.UserID
	email          kernel.Email
	number         UserNumber
	passwordHash   string
	lastLoggedInAt time.Time // UTC
	signedUpAt     time.Time // UTC
	eventManager   *kernelEvent.Manager
}

func NewUser(id kernel.UserID, email kernel.Email, number UserNumber, passwordHash string, lastLoggedInAt time.Time, signedUpAt time.Time) *User {
	return &User{
		id:             id,
		email:          email,
		number:         number,
		passwordHash:   passwordHash,
		lastLoggedInAt: lastLoggedInAt,
		signedUpAt:     signedUpAt,
		eventManager:   kernelEvent.NewEventManager(),
	}
}

// behavior

// SignUp 注册返回新用户
func SignUp(
	email kernel.Email,
	password Password,
	encryptor kernel.HashEncryptor,
	idGenerator kernel.IDGenerator,
	numberGenerator kernel.NumberGenerator,
) (*User, error) {
	if err := password.Validate(); err != nil {
		return nil, err
	}
	if err := email.Validate(); err != nil {
		return nil, err
	}

	hash, err := password.Hash(encryptor)
	if err != nil {
		return nil, err
	}

	id := idGenerator.Generate()
	number := numberGenerator.Generate()

	now := time.Now().UTC()
	user := NewUser(kernel.UserID(id), email, UserNumber(number), hash, now, now)

	user.eventManager.RecordEvent(kernelEvent.CreateUserSignedUpEvent(user.id, email, idGenerator))

	return user, nil
}

// Login 登录返回RefreshToken
func (u *User) Login(
	password Password,
	generator kernel.IDGenerator,
	comparator kernel.HashComparator,
	refreshGenerator RandomStringGenerator,
	accessGenerator kernel.AccessTokenGenerator,
) (*RefreshToken, kernel.AccessToken, error) {
	if err := password.Validate(); err != nil {
		return nil, "", err
	}

	if err := comparator.Compare(u.passwordHash, password.String()); err != nil {
		return nil, "", err
	}

	now := time.Now()

	u.lastLoggedInAt = now

	u.eventManager.RecordEvent(kernelEvent.CreateUserLoggedInEvent(u.id, generator))

	refreshTokenString, err := refreshGenerator.Generate()
	if err != nil {
		return nil, "", err
	}

	refreshToken := CreateRefreshToken(refreshTokenString, u.id)

	accessToken, err := accessGenerator.Generate(u.id)
	if err != nil {
		return nil, "", err
	}

	return refreshToken, accessToken, nil
}

// getter

func (u *User) ID() kernel.UserID {
	return u.id
}

func (u *User) SignedUpAt() time.Time {
	return u.signedUpAt
}

func (u *User) LastLoggedInAt() time.Time {
	return u.lastLoggedInAt
}

func (u *User) Entity() *kernelEvent.Manager {
	return u.eventManager
}

func (u *User) PasswordHash() string {
	return u.passwordHash
}
func (u *User) Email() kernel.Email {
	return u.email
}

func (u *User) Number() UserNumber {
	return u.number
}

func (u *User) GetEvents() []kernelEvent.Event {
	return u.eventManager.GetEvents()
}
