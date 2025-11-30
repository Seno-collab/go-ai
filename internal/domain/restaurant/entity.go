package restaurant

import "github.com/google/uuid"

type Entity struct {
	Name        string
	Description *string
	Address     *string
	Category    *string
	City        *string
	District    *string
	LogoUrl     *string
	BannerUrl   *string
	PhoneNumber *string
	WebsiteUrl  *string
	Email       *string
	UserID      uuid.UUID
}
