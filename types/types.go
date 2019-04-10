package types

// User describes the users
type User struct {
	Login    string `json:"login"    bson:"login"`
	Password string `json:"password" bson:"password"`
}

// Token describes the auth
type Token struct {
	AccessToken         string `json:"access_token"          bson:"access_token"`
	ExpiresAccessToken  int64  `json:"expires_access_token"  bson:"expires_access_token"`
	RefreshToken        string `json:"refresh_token"         bson:"refresh_token"`
	ExpiresRefreshToken int64  `json:"expires_refresh_token" bson:"expires_refresh_token"`
}
