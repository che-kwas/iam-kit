package util

import (
	"github.com/speps/go-hashids"
)

const minLen = 6

// GetInstanceID returns id format like: secret-2v69o5
func GetInstanceID(uid uint64, prefix string) string {
	hd := hashids.NewData()
	hd.Alphabet = "abcdefghijklmnopqrstuvwxyz1234567890"
	hd.Salt = "x20k5x"
	hd.MinLength = minLen

	h, _ := hashids.NewWithData(hd)
	hashed, _ := h.Encode([]int{int(uid)})

	return prefix + hashed
}
