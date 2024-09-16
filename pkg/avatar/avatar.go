package avatar

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// GetAvatarFromID получает изображение в PNG аватарки пользователя Discord по его ID.
func GetAvatarFromID(userID string) ([]byte, error) {
	// Отправка GET-запроса
	resp, err := http.Get(fmt.Sprintf("https://avatar.cdev.shop/%s", userID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusNotModified && resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid response status:" + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
