package paginator

import (
	"encoding/base64"
	"errors"
	"strings"
)

var (
	ErrInvalidPageToken = errors.New("invalid page token")
)

func GeneratePageToken(id string, param string) string {
	return base64.StdEncoding.EncodeToString([]byte(id + "_" + param))
}

func ParsePageToken(pageToken string) (string, string, error) {
	decodedPageToken, err := base64.StdEncoding.DecodeString(pageToken)
	if err != nil {
		return "", "", ErrInvalidPageToken
	}

	splitPageToken := strings.Split(string(decodedPageToken), "_")
	if len(splitPageToken) != 2 {
		return "", "", ErrInvalidPageToken
	}

	return splitPageToken[0], splitPageToken[1], nil
}
