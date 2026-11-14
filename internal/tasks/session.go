package tasks

import "net/http"

func Authorize(r *http.Request, users *UserRepo) error {
	_, err := AuthorizeUser(r, users)
	if err != nil {
		return ErrNotFound
	}

	return nil
}
