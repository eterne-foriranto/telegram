package main

import (
	"crypto/rand"
	"strings"
)

const KeyLength = 2

func char(rune byte) string {
	if rune == 0 {
		return string(45)
	}

	if rune < 11 {
		return string(47 + rune)
	}

	if rune < 37 {
		return string(54 + rune)
	}

	if rune > 37 {
		return string(59 + rune)
	}

	return string(95)
}

func randomDigits(b []byte) string {
	chars := make([]string, 4)
	chars[0] = char(b[0] >> 2)
	left1 := b[0] << 6
	left1 = left1 >> 2
	right1 := b[1] >> 4
	chars[1] = char(left1 | right1)
	left2 := b[1] << 4
	left2 = left2 >> 2
	right2 := b[2] >> 6
	chars[2] = char(left2 | right2)
	chars3 := b[2] << 2
	chars[3] = char(chars3 >> 2)
	return strings.Join(chars, "")
}

func inviteKey() string {
	b := make([]byte, 3)
	key := ""
	for i := 0; i < KeyLength; i++ {
		_, err := rand.Read(b)
		handleError(err)
		key += randomDigits(b)
	}
	return key
}
