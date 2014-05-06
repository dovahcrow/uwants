package myfans

import (
	"testing"
)

func test_post(t *testing.T) {
	cl := New(`doomsplayer`, `1cd3599df`)
	cl.Login()
	t.Log(cl.SendReply(`3608797`, `哈哈哈`))
}

func Test_NewTopic(t *testing.T) {
	cl := New(`doomsplayer`, `1cd3599df`)
	cl.Login()
	t.Log(cl.NewTopic(`216`, `test`, `哈哈哈`))
}
