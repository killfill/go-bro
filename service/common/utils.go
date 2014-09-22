package utils

import (
	"crypto/rand"
	"strconv"
	"strings"
)

//Taken from http://stackoverflow.com/questions/12771930/what-is-the-fastest-way-to-generate-a-long-random-string-in-go
func Rand_str(str_size int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func GetHostAndPortFromURL(dsn string) (host string, port int) {
	at := strings.Index(dsn, "@")
	addr := strings.Split(dsn[at+1:], "/")[0]
	s := strings.Split(addr, ":")

	//What a mess
	if len(s) > 1 {
		host = s[0]
		port, _ = strconv.Atoi(s[1])
	} else {
		host, port = s[0], 5432
	}
	return
}
