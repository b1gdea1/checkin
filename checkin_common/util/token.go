package util

import (
	"math/rand"
	"strconv"
	"time"
)

func init() {
	t := time.Now().Unix()
	rand.Seed(t)
}
func GetToken() string {
	return strconv.Itoa(rand.Int()) + strconv.Itoa(int(time.Now().Unix()))
}
