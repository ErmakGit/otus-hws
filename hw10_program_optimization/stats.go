package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	easyjson "github.com/mailru/easyjson"
)

var ErrorEmptyDomainProvided = errors.New("provided empty domain")

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat, 0)

	if strings.Trim(domain, " ") == "" {
		return result, ErrorEmptyDomainProvided
	}
	var user User

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return result, err
		}

		if strings.HasSuffix(user.Email, "."+domain) {
			key := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[key]++
		}
	}

	return result, nil
}
