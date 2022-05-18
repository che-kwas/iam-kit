package util

import "crypto/rand"

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// RandString returns a random string of a specified length.
func RandString(length int) string {
	output := make([]byte, length)

	// We will take length bytes, one byte for each character of output.
	randomness := make([]byte, length)

	// read all random
	if _, err := rand.Read(randomness); err != nil {
		panic(err)
	}

	l := len(alphabet)
	// fill output
	for pos := range output {
		// get random item
		random := randomness[pos]

		// random % 64
		randomPos := random % uint8(l)

		// put into output
		output[pos] = alphabet[randomPos]
	}

	return string(output)
}
