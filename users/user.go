package users

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type WatchList struct {
	ID             string `json:"id"`
	UserID         string `json:"userId"`
	Symbol         string `json:"symbol"`
	SequenceNumber int    `json:"sequenceNumber"`
}
