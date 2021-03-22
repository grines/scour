package awsdata

import (
	"math/rand"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func SetTracking() string {
	rando := String(10)
	os.Setenv("AWS_EXECUTION_ENV", rando)
	return rando
}

func SetTrackingAction(t string, call string) string {
	rando := String(10)
	os.Setenv("AWS_EXECUTION_ENV", t+"/"+rando+"/"+call)
	return t + "/" + rando + "/" + call
}
