package core

type (

	// System stores system information.
	System struct {
		ClientID     string
		ClientSecret string
		Version      string
		Addresses    map[string]*Address
	}
)
