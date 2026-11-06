package tty

import "math"
import "time"

type ActivityStyle int

const (
	ActivityBlocks ActivityStyle = iota
	ActivityDiamonds
)

type Activity struct {
	Style      ActivityStyle
	Color      RGBA
	Background RGBA
	TrailSteps int
	HoldStart  int
	HoldEnd    int
	Interval   time.Duration
}

func NewActivity(color RGBA, background RGBA) Activity {

	return Activity{
		Style:      ActivityBlocks,
		Color:      color,
		Background: background,
		TrailSteps: 6,
		HoldStart:  30,
		HoldEnd:    9,
		Interval:   80 * time.Millisecond,
	}

}

func (activity Activity) Render(buffer *Buffer, x int, y int, width int, elapsed time.Duration) {

	if width <= 0 {
		return
	}

	frame := activity.frame(elapsed)
	position := activity.activePosition(frame, width)
	forward := activity.isForward(frame, width)
	head, trail, inactive := activity.glyphs()
	trail_colors := activity.trailColors()

	for index := 0; index < width; index++ {

		distance := position - index

		if forward == false {
			distance = index - position
		}

		if distance == 0 {

			style := Style{
				Foreground:    activity.Color,
				Background:    activity.Background,
				HasForeground: true,
				HasBackground: true,
			}

			buffer.Set(x+index, y, head, style)

		} else if distance > 0 && distance < activity.TrailSteps {

			color := trail_colors[distance]
			blended := color.Blend(activity.Background)

			style := Style{
				Foreground:    blended,
				Background:    activity.Background,
				HasForeground: true,
				HasBackground: true,
			}

			buffer.Set(x+index, y, trail, style)

		} else {

			style := Style{
				Foreground:    activity.Color.WithAlpha(0.20).Blend(activity.Background),
				Background:    activity.Background,
				HasForeground: true,
				HasBackground: true,
			}

			buffer.Set(x+index, y, inactive, style)

		}

	}

}

func (activity Activity) trailColors() []RGBA {

	steps := max(activity.TrailSteps, 2)
	result := make([]RGBA, steps)

	for index := 0; index < steps; index++ {

		alpha := 1.0
		brightness := 1.0

		if index == 0 {
			alpha = 1.0
		} else if index == 1 {
			alpha = 0.9
			brightness = 1.15
		} else {
			alpha = math.Pow(0.65, float64(index-1))
		}

		result[index] = NewRGBA(
			math.Min(1.0, activity.Color.R*brightness),
			math.Min(1.0, activity.Color.G*brightness),
			math.Min(1.0, activity.Color.B*brightness),
			alpha,
		)

	}

	return result

}

func (activity Activity) glyphs() (string, string, string) {

	if activity.Style == ActivityDiamonds {
		return "◆", "⬩", "·"
	} else {
		return "■", "▪", "⬝"
	}

}

func (activity Activity) interval() time.Duration {

	if activity.Interval <= 0 {
		return 80 * time.Millisecond
	} else {
		return activity.Interval
	}

}

func (activity Activity) frame(elapsed time.Duration) int {

	return int(elapsed / activity.interval())

}

func (activity Activity) isForward(frame int, width int) bool {

	hold_end := max(0, activity.HoldEnd)
	hold_start := max(0, activity.HoldStart)
	forward := width
	backward := max(1, width-1)
	total := forward + hold_end + backward + hold_start
	position := frame % max(1, total)

	return position < forward || (position >= forward && position < forward+hold_end)

}

func (activity Activity) activePosition(frame int, width int) int {

	if width <= 1 {
		return 0
	}

	hold_end := max(0, activity.HoldEnd)
	hold_start := max(0, activity.HoldStart)
	forward := width
	backward := max(1, width-1)
	total := forward + hold_end + backward + hold_start
	position := frame % max(1, total)

	if position < forward {
		return position
	} else if position < forward+hold_end {
		return width - 1
	} else if position < forward+hold_end+backward {
		index := position - forward - hold_end
		return width - 2 - index
	} else {
		return 0
	}

}

func (activity Activity) Frames(width int) []string {

	if width <= 0 {
		return []string{}
	}

	hold_end := max(0, activity.HoldEnd)
	hold_start := max(0, activity.HoldStart)
	total := width + hold_end + max(1, width-1) + hold_start
	frames := make([]string, total)
	head, _, inactive := activity.glyphs()

	for frame := 0; frame < total; frame++ {

		position := activity.activePosition(frame, width)
		cells := make([]rune, width)

		for index := 0; index < width; index++ {

			if index == position {
				cells[index] = []rune(head)[0]
			} else {
				cells[index] = []rune(inactive)[0]
			}

		}

		frames[frame] = string(cells)

	}

	return frames

}
