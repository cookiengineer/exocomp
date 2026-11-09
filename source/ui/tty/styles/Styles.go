package styles

import "exocomp/ui/tty/colors"
import utils_tty "exocomp/utils/tty"

var (
	Background   = utils_tty.NewStyle(colors.BackgroundBase, colors.BackgroundBase)
	Text         = utils_tty.NewStyle(colors.TextBase, colors.BackgroundBase)
	TextMuted    = utils_tty.NewStyle(colors.TextMuted, colors.BackgroundBase)
	Border       = utils_tty.NewStyle(colors.BorderBase, colors.BackgroundBase)
	BorderRaised = utils_tty.NewStyle(colors.BorderBase, colors.BackgroundRaisedBase)
	Selection    = utils_tty.NewStyle(colors.TextActionPrimaryFocused, colors.BackgroundActionPrimaryFocused)

	SidebarHeader = utils_tty.NewStyle(colors.TextFeedbackInfoBase, colors.BackgroundBase).WithAttributes(utils_tty.AttrBold)
	SidebarError  = utils_tty.NewStyle(colors.TextFeedbackErrorBase, colors.BackgroundBase)

	Dialog       = utils_tty.NewStyle(colors.TextBase, colors.BackgroundRaisedBase)
	DialogBorder = utils_tty.NewStyle(colors.BorderBase, colors.BackgroundRaisedBase)
	DialogTitle  = utils_tty.NewStyle(colors.TextFeedbackInfoBase, colors.BackgroundRaisedBase).WithAttributes(utils_tty.AttrBold)
	DialogInput  = utils_tty.NewStyle(colors.TextBase, colors.BackgroundRaisedHigh)
	DialogMuted  = utils_tty.NewStyle(colors.TextMuted, colors.BackgroundRaisedBase)

	Palette = utils_tty.NewStyle(colors.TextBase, colors.BackgroundRaisedHigh)

	MarkdownText            = utils_tty.NewStyle(colors.MarkdownText, colors.BackgroundBase)
	MarkdownHeading         = utils_tty.NewStyle(colors.MarkdownHeading, colors.BackgroundBase).WithAttributes(utils_tty.AttrBold)
	MarkdownStrong          = utils_tty.NewStyle(colors.MarkdownStrong, colors.BackgroundBase).WithAttributes(utils_tty.AttrBold)
	MarkdownEmphasis        = utils_tty.NewStyle(colors.MarkdownEmphasis, colors.BackgroundBase).WithAttributes(utils_tty.AttrItalic)
	MarkdownCode            = utils_tty.NewStyle(colors.MarkdownCode, colors.BackgroundRaisedHigh)
	MarkdownCodeBlock       = utils_tty.NewStyle(colors.MarkdownCodeBlock, colors.BackgroundRaisedBase)
	MarkdownLink            = utils_tty.NewStyle(colors.MarkdownLinkText, colors.BackgroundBase).WithAttributes(utils_tty.AttrUnderline)
	MarkdownLinkURL         = utils_tty.NewStyle(colors.TextMuted, colors.BackgroundBase)
	MarkdownBlockQuote      = utils_tty.NewStyle(colors.MarkdownBlockQuote, colors.BackgroundBase).WithAttributes(utils_tty.AttrItalic)
	MarkdownListItem        = utils_tty.NewStyle(colors.MarkdownListItem, colors.BackgroundBase)
	MarkdownListEnumeration = utils_tty.NewStyle(colors.MarkdownListEnumeration, colors.BackgroundBase)
	MarkdownHorizontalRule  = utils_tty.NewStyle(colors.MarkdownHorizontalRule, colors.BackgroundBase)
	MarkdownImage           = utils_tty.NewStyle(colors.MarkdownImage, colors.BackgroundBase)
	MarkdownImageText       = utils_tty.NewStyle(colors.MarkdownImageText, colors.BackgroundBase)
	MarkdownStrike          = utils_tty.NewStyle(colors.TextMuted, colors.BackgroundBase).WithAttributes(utils_tty.AttrStrikethrough)

	TranscriptUser      = utils_tty.NewStyle(colors.TextFeedbackSuccessBase, colors.BackgroundBase).WithAttributes(utils_tty.AttrBold)
	TranscriptAssistant = utils_tty.NewStyle(colors.TextBase, colors.BackgroundBase).WithAttributes(utils_tty.AttrBold)
	TranscriptTool      = utils_tty.NewStyle(colors.TextMuted, colors.BackgroundBase).WithAttributes(utils_tty.AttrBold)
	TranscriptSystem    = utils_tty.NewStyle(colors.TextFeedbackWarningBase, colors.BackgroundBase).WithAttributes(utils_tty.AttrBold)
)
