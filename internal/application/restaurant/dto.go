package restaurantapp

type CreateRestaurantRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Category    string `json:"category"`
	City        string `json:"city"`
	District    string `json:"district"`
	LogoUrl     string `json:"logo_url"`
	BannerUrl   string `json:"banner_url"`
	PhoneNumber string `json:"phone_number"`
	WebsiteUrl  string `json:"website_url"`
	Email       string `json:"email"`
}

type UpLoadLogoRequest struct {
}
