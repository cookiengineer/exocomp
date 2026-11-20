package components

import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"

func drawSidebarScrollbar(buffer *utils_tty.Buffer, x int, top int, height int, total int, scroll int) {

	if height <= 0 || total <= height {
		return
	}

	track := height * 2
	thumb := max(2, track*height/total)
	max_scroll := max(1, total-height)
	start := scroll * (track - thumb) / max_scroll

	for row := 0; row < height; row++ {

		cell_start := row * 2
		cell_end := cell_start + 2
		overlap_start := max(start, cell_start)
		overlap_end := min(start+thumb, cell_end)

		if overlap_end > overlap_start {

			glyph := "█"

			if overlap_end-overlap_start < 2 {
				glyph = "▄"
			}

			if overlap_end-overlap_start < 2 && overlap_start == cell_start {
				glyph = "▀"
			}

			buffer.Set(x, top+row, glyph, styles.ScrollbarThumb)

		}

	}

}
