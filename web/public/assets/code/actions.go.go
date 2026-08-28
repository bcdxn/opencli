package cliapp

import (
	"context"
	"fmt"

	"gen.com/example/internal/gencli"
	"github.com/bcdxn/opencli/spec"
)

func NewActions(version string) Actions {
	return Actions{
		version: version,
	}
}

type Actions struct {
	version string
}

func (a Actions) PetstoreList(ctx context.Context, args gencli.PetstoreListArgs) error {
	fmt.Println("action ran")
	return nil
}
func (a Actions) PetstorePetAdd(ctx context.Context, args gencli.PetstorePetAddArgs, flags gencli.PetstorePetAddFlags) error {
	fmt.Println("action ran", args, flags)
	return nil
}
func (a Actions) PetstorePetUpdate(ctx context.Context, args gencli.PetstorePetUpdateArgs, flags gencli.PetstorePetUpdateFlags) error {
	fmt.Println("action ran", args, flags)
	return nil
}
func (a Actions) PetstorePetFindByStatus(ctx context.Context, flags gencli.PetstorePetFindByStatusFlags) error {
	fmt.Println("action ran", flags)
	return nil
}
func (a Actions) PetstorePetFindByTags(ctx context.Context, flags gencli.PetstorePetFindByTagsFlags) error {
	fmt.Println("action ran", flags)
	return nil
}
func (a Actions) PetstorePetGet(ctx context.Context, flags gencli.PetstorePetGetFlags) error {
	fmt.Println("action ran", flags)
	return nil
}
func (a Actions) PetstorePetUpdateForm(ctx context.Context, flags gencli.PetstorePetUpdateFormFlags) error {
	fmt.Println("action ran", flags)
	return nil
}
func (a Actions) PetstorePetDelete(ctx context.Context, flags gencli.PetstorePetDeleteFlags) error {
	fmt.Println("action ran", flags)
	return nil
}
func (a Actions) PetstorePetUploadImage(ctx context.Context, args gencli.PetstorePetUploadImageArgs, flags gencli.PetstorePetUploadImageFlags) error {
	fmt.Println("action ran", args, flags)
	return nil
}
func (a Actions) PetstoreStoreInventory(ctx context.Context) error {
	fmt.Println("action ran")
	return nil
}
func (a Actions) PetstoreStoreOrderPlace(ctx context.Context, args gencli.PetstoreStoreOrderPlaceArgs, flags gencli.PetstoreStoreOrderPlaceFlags) error {
	fmt.Println("action ran", args, flags)
	return nil
}
func (a Actions) PetstoreStoreOrderGet(ctx context.Context, flags gencli.PetstoreStoreOrderGetFlags) error {
	fmt.Println("action ran", flags)
	return nil
}
func (a Actions) PetstoreStoreOrderDelete(ctx context.Context, args gencli.PetstoreStoreOrderDeleteArgs) error {
	fmt.Println("action ran", args)
	return nil
}
func (a Actions) PetstoreUserCreate(ctx context.Context, args gencli.PetstoreUserCreateArgs, flags gencli.PetstoreUserCreateFlags) error {
	fmt.Println("action ran", args, flags)
	return nil
}
func (a Actions) PetstoreUserCreateWithList(ctx context.Context, args gencli.PetstoreUserCreateWithListArgs) error {
	fmt.Println("action ran", args)
	return nil
}
func (a Actions) PetstoreUserLogin(ctx context.Context, flags gencli.PetstoreUserLoginFlags) error {
	fmt.Println("action ran", flags)
	return nil
}
func (a Actions) PetstoreUserLogout(ctx context.Context) error {
	fmt.Println("action ran")
	return nil
}
func (a Actions) PetstoreUserGet(ctx context.Context, flags gencli.PetstoreUserGetFlags) error {
	fmt.Println("action ran", flags)
	return nil
}
func (a Actions) PetstoreUserUpdate(ctx context.Context, flags gencli.PetstoreUserUpdateFlags) error {
	fmt.Println("action ran", flags)
	return nil
}
func (a Actions) PetstoreUserDelete(ctx context.Context, flags gencli.PetstoreUserDeleteFlags) error {
	fmt.Println("action ran", flags)
	return nil
}
func (a Actions) HelpFunc(cmd *spec.CommandItem) {
	gencli.DefaultHelpFunc(a, cmd)
}
func (a Actions) UsageFunc(cmd *spec.CommandItem) error {
	gencli.DefaultUsageFunc(a, cmd)
	return nil
}
func (a Actions) IOStreams() gencli.IOStreams {
	return gencli.DefaultIOS()
}
func (a Actions) Version() string {
	return a.version
}
