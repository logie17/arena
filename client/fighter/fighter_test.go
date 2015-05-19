package fighter

import (
	"testing"
)

func TestMoveLeft(t *testing.T) {
	enemyReply := make(chan Line, 1)
	NewFighter(4, 4, 1, "enemy", enemyReply)

	reply := make(chan Line)
	go func() {
		subject := NewFighter(4, 5, 2, "me", reply)
		subject.Left()
	}()
	if response := <-reply; response.Action != "draw_me" && response.X != 4 && response.Y != 5 {
		t.Errorf("A draw action should come back in the reply")
	}

	if response := <-reply; response.Action != "hide" && response.X != 4 && response.Y != 5 {
		t.Errorf("A draw action should come back in the reply")
	}

	if response := <-reply; response.Action != "draw_me" && response.X != 4 && response.Y != 5 {
		t.Errorf("A draw action should come back in the reply")
	}

	close(reply)
}
