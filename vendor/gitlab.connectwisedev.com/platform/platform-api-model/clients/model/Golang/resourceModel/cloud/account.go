package cloud

//RegistrationTransaction returns the transaction ID
type RegistrationTransaction struct {
	TransactionID string `json:"transactionid"`
}

//Client returns the clientID
type Client struct {
	ID string `json:"clientid"`
}

//User return the user detail
type User struct {
	UserName string `json:"username"`
}

//AuthorizationStatus returns the Authorization status
type AuthorizationStatus struct {
	UserName string `json:"username"`
	Status   string `json:"status"`
}
