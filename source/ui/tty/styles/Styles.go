package styles

import "exocomp/ui/tty/colors"
import utils_tty "exocomp/utils/tty"

var (
	Background   = utils_tty.NewStyle(colors.MainBackground, colors.MainBackground)
	Text         = utils_tty.NewStyle(colors.MainColor, colors.MainBackground)
	TextMuted    = utils_tty.NewStyle(colors.TextMuted, colors.MainBackground)
	Border       = utils_tty.NewStyle(colors.BorderBase, colors.MainBackground)
	BorderRaised = utils_tty.NewStyle(colors.BorderBase, colors.AsideBackground)
	Selection    = utils_tty.NewStyle(colors.AccentText, colors.Accent)

	SidebarPanel      = utils_tty.NewStyle(colors.AsideColor, colors.AsideBackground)
	SidebarTitle      = utils_tty.NewStyle(colors.AsideColor, colors.AsideBackground).WithAttributes(utils_tty.AttrBold)
	SidebarSection    = utils_tty.NewStyle(colors.AsideColor, colors.AsideBackground).WithAttributes(utils_tty.AttrBold)
	SidebarText       = utils_tty.NewStyle(colors.AsideColor, colors.AsideBackground)
	SidebarMuted      = utils_tty.NewStyle(colors.TextMuted, colors.AsideBackground)
	SidebarSelected   = utils_tty.NewStyle(colors.Accent, colors.AsideBackground).WithAttributes(utils_tty.AttrBold)
	SidebarPanelError = utils_tty.NewStyle(colors.FeedbackError, colors.AsideBackground)

	ScrollbarThumb = utils_tty.NewStyle(colors.ScrollbarThumb, colors.AsideBackground)

	PromptPanel       = utils_tty.NewStyle(colors.FooterColor, colors.FooterBackground)
	PromptText        = utils_tty.NewStyle(colors.FooterColor, colors.FooterBackground)
	PromptMuted       = utils_tty.NewStyle(colors.TextMuted, colors.FooterBackground)
	PromptBorder      = utils_tty.NewStyle(colors.BorderBase, colors.FooterBackground)
	PromptBorderFocus = utils_tty.NewStyle(colors.FeedbackSuccess, colors.FooterBackground)
	PromptCap         = utils_tty.NewStyle(colors.FooterBackground, colors.MainBackground)

	FooterKey  = utils_tty.NewStyle(colors.MainColor, colors.MainBackground)
	FooterHint = utils_tty.NewStyle(colors.TextMuted, colors.MainBackground)

	Dialog       = utils_tty.NewStyle(colors.ArticleColor, colors.ArticleBackground)
	DialogBorder = utils_tty.NewStyle(colors.ArticleBorder, colors.ArticleBackground)
	DialogTitle  = utils_tty.NewStyle(colors.FeedbackInfo, colors.ArticleBackground).WithAttributes(utils_tty.AttrBold)
	DialogInput  = utils_tty.NewStyle(colors.MainColor, colors.BackgroundRaisedHigh)
	DialogMuted  = utils_tty.NewStyle(colors.TextMuted, colors.ArticleBackground)

	Palette = utils_tty.NewStyle(colors.MainColor, colors.BackgroundRaisedHigh)

	MarkdownText            = utils_tty.NewStyle(colors.MarkdownText, colors.MainBackground)
	MarkdownHeading         = utils_tty.NewStyle(colors.MarkdownHeading, colors.MainBackground).WithAttributes(utils_tty.AttrBold)
	MarkdownStrong          = utils_tty.NewStyle(colors.MarkdownStrong, colors.MainBackground).WithAttributes(utils_tty.AttrBold)
	MarkdownEmphasis        = utils_tty.NewStyle(colors.MarkdownEmphasis, colors.MainBackground).WithAttributes(utils_tty.AttrItalic)
	MarkdownCode            = utils_tty.NewStyle(colors.MarkdownCode, colors.BackgroundRaisedHigh)
	MarkdownCodeBlock       = utils_tty.NewStyle(colors.MarkdownCodeBlock, colors.BackgroundRaised)
	MarkdownLink            = utils_tty.NewStyle(colors.MarkdownLinkText, colors.MainBackground).WithAttributes(utils_tty.AttrUnderline)
	MarkdownLinkURL         = utils_tty.NewStyle(colors.TextMuted, colors.MainBackground)
	MarkdownBlockQuote      = utils_tty.NewStyle(colors.MarkdownBlockQuote, colors.MainBackground).WithAttributes(utils_tty.AttrItalic)
	MarkdownListItem        = utils_tty.NewStyle(colors.MarkdownListItem, colors.MainBackground)
	MarkdownListEnumeration = utils_tty.NewStyle(colors.MarkdownListEnumeration, colors.MainBackground)
	MarkdownHorizontalRule  = utils_tty.NewStyle(colors.MarkdownHorizontalRule, colors.MainBackground)
	MarkdownImage           = utils_tty.NewStyle(colors.MarkdownImage, colors.MainBackground)
	MarkdownImageText       = utils_tty.NewStyle(colors.MarkdownImageText, colors.MainBackground)
	MarkdownStrike          = utils_tty.NewStyle(colors.TextMuted, colors.MainBackground).WithAttributes(utils_tty.AttrStrikethrough)

	TranscriptUser      = utils_tty.NewStyle(colors.FeedbackSuccess, colors.MainBackground).WithAttributes(utils_tty.AttrBold)
	TranscriptAssistant = utils_tty.NewStyle(colors.MainColor, colors.MainBackground).WithAttributes(utils_tty.AttrBold)
	TranscriptTool      = utils_tty.NewStyle(colors.TextMuted, colors.MainBackground).WithAttributes(utils_tty.AttrBold)
	TranscriptSystem    = utils_tty.NewStyle(colors.FeedbackWarning, colors.MainBackground).WithAttributes(utils_tty.AttrBold)

	ReasoningHeader = utils_tty.NewStyle(colors.FeedbackThinking, colors.MainBackground)
	ReasoningText   = utils_tty.NewStyle(colors.TextMuted, colors.MainBackground)
	SystemHeader    = utils_tty.NewStyle(colors.FeedbackInfo, colors.MainBackground)
	ToolHeader      = utils_tty.NewStyle(colors.Accent, colors.MainBackground)
	CollapsedHint   = utils_tty.NewStyle(colors.TextMuted, colors.MainBackground)

	BorderSystem    = utils_tty.NewStyle(colors.FeedbackWarning, colors.MainBackground)
	BorderAssistant = utils_tty.NewStyle(colors.FeedbackInfo, colors.MainBackground)
	BorderTool      = utils_tty.NewStyle(colors.FeedbackError, colors.MainBackground)
	BorderUser      = utils_tty.NewStyle(colors.FeedbackSuccess, colors.MainBackground)
	BorderThinking  = utils_tty.NewStyle(colors.FeedbackThinking, colors.MainBackground)
)
