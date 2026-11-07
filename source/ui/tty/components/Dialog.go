package components

import utils_tty "exocomp/utils/tty"

type Dialog interface {
	Component
	GetQuestion() string
	GetCursor() (int, int, bool)
	HandleKey(key utils_tty.Key) (string, bool, bool)
}

var _ Dialog = (*QuestionDialog)(nil)
var _ Dialog = (*ChoiceDialog)(nil)
