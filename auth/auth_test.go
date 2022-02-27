package auth

import "testing"

func Test(t *testing.T) {

	var auth Auth
	auth.New()

	auth.Send(SendOpts{
		Path:   "/fills",
		Body:   nil,
		Method: "POST",
	})

}
