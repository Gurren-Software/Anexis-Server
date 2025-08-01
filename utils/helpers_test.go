package utils_test

import (
	"crypto/rsa"
	"testing"

	"github.com/Treefle-labs/anexis-server/utils"
)

var privateKey *rsa.PrivateKey

const USER = "doni"

func TestGenerator(t *testing.T) {
	key, err := utils.GenerateRSAKeys(USER)
	if err != nil {
		t.Errorf("Failed to generate RSA keys: %v", err)
	}
	privateKey = key
}
