package application

import (
	"context"
	"gochat/internal/profile/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type UpdateRoomProfileUseCase kernel.UseCase[*UpdateRoomProfileInput, *kernel.NoOutput]

type UpdateRoomProfileInput struct {
	UserID       kernel.UserID
	RoomID       kernel.RoomID
	Name         string
	Introduction string
}

func (r *UpdateRoomProfileInput) Validate() error {
	if len(r.RoomID) == 0 || len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id or room id is empty")
	}

	if len(r.Name) == 0 && len(r.Introduction) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "everything is empty")
	}

	return nil
}

type updateRoomProfileUseCase struct {
	updater        domain.RoomProfileUpdater
	profileFinder  domain.RoomProfileFinder
	roomshipFinder domain.RoomshipFinderByUserIDAndRoomID
}

func NewUpdateRoomProfileUseCase(
	updater domain.RoomProfileUpdater,
	profileFinder domain.RoomProfileFinder,
	roomshipFinder domain.RoomshipFinderByUserIDAndRoomID,
) (UpdateRoomProfileUseCase, error) {
	if err := utils.CheckInterfaces(
		updater,
		profileFinder,
		roomshipFinder,
	); err != nil {
		return nil, err
	}

	return &updateRoomProfileUseCase{
		updater:        updater,
		profileFinder:  profileFinder,
		roomshipFinder: roomshipFinder,
	}, nil
}

func (uc *updateRoomProfileUseCase) Execute(ctx context.Context, input *UpdateRoomProfileInput) (*kernel.NoOutput, error) {
	//确认权限
	roomship, err := uc.roomshipFinder.FindByUserIDAndRoomID(ctx, input.UserID, input.RoomID)
	if err != nil {
		return nil, err
	}

	if !roomship.IsOwner() {
		return nil, domain.ErrNoPermission
	}

	profile, err := uc.profileFinder.FindByID(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	uc.updatesProfile(profile, input)

	if err := uc.updater.Update(ctx, profile); err != nil {
		return nil, err
	}

	//更新后也可发布通知事件来提醒房间成员

	return nil, nil
}

func (uc *updateRoomProfileUseCase) updatesProfile(profile *domain.RoomProfile, input *UpdateRoomProfileInput) {
	if input.Name != "" {
		profile.UpdateName(input.Name)
	}
	if input.Introduction != "" {
		profile.UpdateIntroduction(input.Introduction)
	}
}
