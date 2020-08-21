package helper

import (
	"crypto/md5"
	"fmt"
	"io"
)

// Md5 make a ma5
func Md5(s string) string {
	h := md5.New()
	_, _ = io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
