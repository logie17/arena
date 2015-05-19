package fighter

import (
	"testing"
)

func TestMoveLeft(t *testing.T) {
	enemyReply := make(chan CommandData, 1)
	NewFighter(4, 4, 1, "enemy", enemyReply)

	reply := make(chan CommandData)
	go func() {
		subject := NewFighter(4, 5, 2, "me", reply)
		subject.Left()
	}()
	if response := <-reply; response.Action != "DRAW" && response.Value[1] != 4 && response.Value[2] != 5 {
		t.Errorf("A draw action should come back in the reply")
	}

	if response := <-reply; response.Action != "HIDE" && response.Value[1] != 4 && response.Value[2] != 5 {
		t.Errorf("A draw action should come back in the reply")
	}

	if response := <-reply; response.Action != "DRAW" && response.Value[1] != 4 && response.Value[2] != 5 {
		t.Errorf("A draw action should come back in the reply")
	}

	close(reply)
}
