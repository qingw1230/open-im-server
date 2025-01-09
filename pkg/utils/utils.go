package utils

import "github.com/jinzhu/copier"

func CopyStructFields(a, b interface{}, fields ...string) error {
	return copier.Copy(a, b)
}
