package tokens

import "crypto/rand"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

func randToken(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // handle error properly in real code
	}

	for i := range n {
		b[i] = letters[int(b[i])%len(letters)]
	}

	return string(b)
}
