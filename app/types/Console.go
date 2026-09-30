package types

type Console struct {
	Messages []ConsoleMessage `json:"messages"`
}

func NewConsole() *Console {

	console := Console{
		Messages: make([]ConsoleMessage, 0),
	}

	return &console

}

func (console *Console) GetMessages(from int) []ConsoleMessage {

	if from < 0 || from >= len(console.Messages) {
		return make([]ConsoleMessage, 0)
	}

	return console.Messages[from:]

}

func (console *Console) Length() int {
	return len(console.Messages)
}
