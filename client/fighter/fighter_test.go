package fighter

import (
	"testing"
)

func TestMoveLeft(t *testing.T) {
	NewFighter(4, 4, 1, "enemy")

	subject := NewFighter(4, 5, 2, "me")
	subject.Left()

	if subject.X() == 4 && subject.Y() == 5 {
		t.Errorf("The subject should not be able to move into an ememy zone")
	}
}
