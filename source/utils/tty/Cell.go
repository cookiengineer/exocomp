package tty

type Cell struct {
	Text  string
	Width int
	Style Style
}

func (cell Cell) IsEmpty() bool {
	return cell.Text == ""
}

func (cell Cell) Equals(other Cell) bool {
	return cell.Text == other.Text && cell.Width == other.Width && cell.Style.Equals(other.Style)
}
