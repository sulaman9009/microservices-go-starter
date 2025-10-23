package domain

type Driver struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profilePicture"`
	CarPlate       string `json:"carPlate"`
	PackageSlug    string `json:"packageSlug"`
}
