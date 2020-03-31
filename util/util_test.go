package util

import (
	"fmt"
	"testing"
)

func TestSha1Stream_Update(t *testing.T) {
	ss := &Sha1Stream{}
	ss.Update([]byte("hello"))
	fmt.Println(ss.Sum())
}
