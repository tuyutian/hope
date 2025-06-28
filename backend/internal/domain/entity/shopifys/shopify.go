package shopifys

type Token struct {
	Token string `json:"access_token"`
	Scope string `json:"scope"`
}
