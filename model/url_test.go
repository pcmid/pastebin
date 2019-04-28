package model

import (
	"testing"
)

func Test_url(t *testing.T) {
	var id uint32 = 300000

	s := GenSortUrl(id)

	i := UrlToID(s)

	if i != id {
		panic("err")
	}
}
