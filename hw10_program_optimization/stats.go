package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	easyjson "github.com/mailru/easyjson" //nolint:depguard
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	var i int
	for scanner.Scan() {
		var user User
		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return DomainStat{}, err
		}

		matched := strings.HasSuffix(user.Email, domain)
		if matched {
			s := strings.SplitN(user.Email, "@", 2)
			if len(s) <= 1 {
				continue
			}
			d := strings.ToLower(s[1])
			result[d]++
		}

		i++
	}

	if err := scanner.Err(); err != nil {
		return DomainStat{}, err
	}

	return result, nil
}
