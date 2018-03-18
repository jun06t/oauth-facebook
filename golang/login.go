package main

func Login(fbs Facebook, code string) (string, error) {
	// get access token
	tok, err := fbs.ExchangeCode(code)
	if err != nil {
		return "", err
	}

	account := struct {
		ID string `json:"id"`
	}{}
	// get profile with access token
	err = fbs.GetMe(tok, &account)
	if err != nil {
		return "", err
	}

	return account.ID, nil
}
