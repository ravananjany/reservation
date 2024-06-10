package resources

type User struct {
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	EmailId   string `json:"email_id,omitempty"`
}

type Reservation struct {
	User   User    `json:"user,omitempty"`
	UserId string  `json:"user_id"`
	From   string  `json:"from,omitempty"`
	To     string  `json:"to,omitempty"`
	Price  float32 `json:"price_paid,omitempty"`
	Seat   string  `json:"seat,omitempty"`
}
