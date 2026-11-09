package colors

import utils_tty "exocomp/utils/tty"

var (
	TextBase  = hex("#eeeeee")
	TextMuted = hex("#808080")

	TextActionPrimaryBase     = hex("#eeeeee")
	TextActionPrimaryDisabled = hex("#808080")
	TextActionPrimaryFocused  = hex("#0a0a0a")
	TextActionPrimarySelected = hex("#fab283")

	TextActionSecondaryBase    = hex("#808080")
	TextActionSecondaryHovered = hex("#eeeeee")

	TextActionDestructiveBase     = hex("#0a0a0a")
	TextActionDestructiveDisabled = hex("#808080")

	TextFormFieldBase     = hex("#eeeeee")
	TextFormFieldHovered  = hex("#fab283")
	TextFormFieldFocused  = hex("#fab283")
	TextFormFieldPressed  = hex("#fab283")
	TextFormFieldDisabled = hex("#808080")
	TextFormFieldSelected = hex("#fab283")

	TextFeedbackErrorBase   = hex("#e06c75")
	TextFeedbackWarningBase = hex("#f5a742")
	TextFeedbackSuccessBase = hex("#7fd88f")
	TextFeedbackInfoBase    = hex("#56b6c2")

	BackgroundBase       = hex("#0a0a0a")
	BackgroundRaisedBase = hex("#141414")
	BackgroundRaisedHigh = hex("#1e1e1e")
	BackgroundRaisedMax  = hex("#4c4c4c")

	BackgroundActionPrimaryBase     = utils_tty.NewColor(0.0, 0.0, 0.0, 0.0)
	BackgroundActionPrimaryHovered  = hex("#141414")
	BackgroundActionPrimaryFocused  = hex("#fab283")
	BackgroundActionPrimarySelected = utils_tty.NewColor(0.0, 0.0, 0.0, 0.0)
	BackgroundActionSecondaryBase   = utils_tty.NewColor(0.0, 0.0, 0.0, 0.0)
	BackgroundActionDestructiveBase = hex("#e06c75")

	BackgroundFormFieldBase = hex("#0a0a0a")

	BackgroundFeedbackErrorBase   = hex("#0a0a0a")
	BackgroundFeedbackWarningBase = hex("#0a0a0a")
	BackgroundFeedbackSuccessBase = hex("#0a0a0a")
	BackgroundFeedbackInfoBase    = hex("#0a0a0a")

	BorderBase    = hex("#484848")
	ScrollbarBase = hex("#606060")

	DiffTextAdded      = hex("#4fd6be")
	DiffTextRemoved    = hex("#c53b53")
	DiffTextContext    = hex("#828bb8")
	DiffTextHunkHeader = hex("#828bb8")

	DiffBackgroundAdded   = hex("#20303b")
	DiffBackgroundRemoved = hex("#37222c")
	DiffBackgroundContext = hex("#141414")

	DiffHighlightAdded   = hex("#b8db87")
	DiffHighlightRemoved = hex("#e26a75")

	DiffLineNumberText              = hex("#8f8f8f")
	DiffLineNumberBackgroundAdded   = hex("#1b2b34")
	DiffLineNumberBackgroundRemoved = hex("#2d1f26")

	SyntaxComment     = hex("#808080")
	SyntaxKeyword     = hex("#9d7cd8")
	SyntaxFunction    = hex("#fab283")
	SyntaxVariable    = hex("#e06c75")
	SyntaxString      = hex("#7fd88f")
	SyntaxNumber      = hex("#f5a742")
	SyntaxType        = hex("#e5c07b")
	SyntaxOperator    = hex("#56b6c2")
	SyntaxPunctuation = hex("#eeeeee")

	MarkdownText            = hex("#eeeeee")
	MarkdownHeading         = hex("#9d7cd8")
	MarkdownLink            = hex("#fab283")
	MarkdownLinkText        = hex("#56b6c2")
	MarkdownCode            = hex("#7fd88f")
	MarkdownBlockQuote      = hex("#e5c07b")
	MarkdownEmphasis        = hex("#e5c07b")
	MarkdownStrong          = hex("#f5a742")
	MarkdownHorizontalRule  = hex("#808080")
	MarkdownListItem        = hex("#fab283")
	MarkdownListEnumeration = hex("#56b6c2")
	MarkdownImage           = hex("#fab283")
	MarkdownImageText       = hex("#56b6c2")
	MarkdownCodeBlock       = hex("#eeeeee")

	DialogBackgroundBase                 = hex("#141414")
	DialogBackgroundActionPrimaryHovered = hex("#1e1e1e")
)
