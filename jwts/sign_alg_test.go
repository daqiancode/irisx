package jwts_test

import (
	"fmt"
	"testing"

	"my_iam/jwts"

	"github.com/stretchr/testify/assert"
)

func TestGenerateEdDSAKeyPair(t *testing.T) {
	pub, pri, err := jwts.GenerateEdDSAKeyPair()
	assert.Nil(t, err)
	fmt.Println(pub)
	fmt.Println(pri)
}
