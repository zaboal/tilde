package tilde

func Subscribe(username string) (password string, error error) {
	return userAdd(username)
}
