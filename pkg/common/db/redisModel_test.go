package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetTokenMapByUidPid(t *testing.T) {
	m := make(map[string]int, 0)
	m["test1"] = 1
	m["test2"] = 2
	m["test3"] = 3
	DB.SetTokenMapByUidPid("userID", 5, m)
}

func Test_GetTokenMapByUidPid(t *testing.T) {
	m, err := DB.GetTokenMapByUidPid("userID", "Web")
	assert.Nil(t, err)
	fmt.Println(m)
}
