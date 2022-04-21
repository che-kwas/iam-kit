package util

import (
	"github.com/speps/go-hashids"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz1234567890"
	salt     = "x20k5x"
	minLen   = 6
)

// GetInstanceID returns id format like: secret-2v69o5
func GetInstanceID(uid uint64, prefix string) string {
	hd := hashids.NewData()
	hd.Alphabet = alphabet
	hd.Salt = salt
	hd.MinLength = minLen

	h, _ := hashids.NewWithData(hd)
	hashed, _ := h.Encode([]int{int(uid)})

	return prefix + hashed
}
