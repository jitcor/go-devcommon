package dcdebug

import (
	"errors"
	"log"
	"testing"
)

func TestLogI(t *testing.T) {
	log.Println()
	LogE(errors.New("Test error"))
}
