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

	// TODO login by facebook account ID.

	return account.ID, nil
}
