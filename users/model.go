package users



type User struct{
	Id int64 `json:"id,pk"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token string `json:"token"`
}
