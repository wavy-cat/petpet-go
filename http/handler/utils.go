package handler

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/wavy-cat/petpet-go/pkg/avatar"
	"github.com/wavy-cat/petpet-go/pkg/cache"
	"net/http"
	"strings"
)

func getAvatarUsingCdev(userId string) ([]byte, string, error) {
	avatarImage, err := avatar.GetAvatarFromID(userId)
	if err != nil {
		return nil, "", err
	}

	hash := md5.Sum(avatarImage)
	hashString := hex.EncodeToString(hash[:])

	return avatarImage, hashString, nil
}

func checkError(err error) (string, string) {
	switch {
	case strings.Contains(err.Error(), "10013"):
		return "Not Found", "User not found"
	case strings.Contains(err.Error(), "50035"):
		return "Incorrect ID", "Check your ID for correctness"
	}

	return "Unknown Error", "Something went wrong"
}

func responseFromCache(w http.ResponseWriter, cache cache.BytesCache, avatarId string) (bool, error) {
	cachedImage, err := cache.Pull(avatarId)
	if err != nil {
		if err.Error() != "not exists" {
			return false, nil
		}
		return false, err
	}

	_, err = w.Write(cachedImage)
	return true, err
}
