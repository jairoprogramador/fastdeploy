package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func Generate(valueOne string, valueTwo string) string {
	valueOneToLower := strings.ToLower(strings.TrimSpace(valueOne))
	valueTwoToLower := strings.ToLower(strings.TrimSpace(valueTwo))

	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%s-%s-%d", valueOneToLower, valueTwoToLower, timestamp)

	hash := sha256.Sum256([]byte(data))
	hashString := hex.EncodeToString(hash[:])

	return hashString
}