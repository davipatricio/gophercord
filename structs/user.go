package structs

type User struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Discriminator string `json:"discriminator"`
}
