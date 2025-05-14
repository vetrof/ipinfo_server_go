package models

type User struct {
	ID       int
	Username string
	Password string
	Token    string
}

type IPInfo struct {
	IP       string
	Hostname string
	City     string
	Region   string
	Country  string
	Loc      string
	Org      string
	Postal   string
	Timezone string
	Readme   string
	UserID   int
}
