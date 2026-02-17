package chat

import (
	"math/rand"
)

func GenerateHashFromContent(content string, length int) string {
	return GenerateHashFromContentWithOffset(content, length, 0)
}

func GenerateHashFromContentWithOffset(content string, length, offset int) string {

	var charset []rune
	for _, char := range content {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			charset = append(charset, char)
		}
	}

	if len(charset) == 0 {
		charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	}

	seed := int64(len(content) + offset)
	for i, char := range content {
		if i < 100 {
			seed += int64(char) * int64(i+offset+1)
		}
	}
	r := rand.New(rand.NewSource(seed))

	hash := make([]rune, length)
	for i := range hash {
		hash[i] = charset[r.Intn(len(charset))]
	}

	return string(hash)
}
