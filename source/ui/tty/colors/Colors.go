package colors

import utils_tty "exocomp/utils/tty"

var (
	// Shared palette. Each exported name maps to the same CSS custom property
	// in app/public/app/theme.css, e.g. MainColor -> --main-color.
	MainColor      = hex("#eeeeee")
	MainBackground = hex("#0a0a0a")

	TextMuted = hex("#808080")

	BackgroundRaised     = hex("#141414")
	BackgroundRaisedHigh = hex("#1e1e1e")
	BackgroundRaisedMax  = hex("#4c4c4c")
	BackgroundComposer   = hex("#0f0f0f")

	BorderBase     = hex("#484848")
	ScrollbarThumb = hex("#606060")

	Accent     = hex("#fab283")
	AccentText = hex("#0a0a0a")

	FeedbackError    = hex("#e06c75")
	FeedbackWarning  = hex("#f5a742")
	FeedbackSuccess  = hex("#7fd88f")
	FeedbackInfo     = hex("#56b6c2")
	FeedbackThinking = hex("#9d7cd8")

	ActionPrimaryBackgroundBase     = utils_tty.NewColor(0.0, 0.0, 0.0, 0.0)
	ActionPrimaryBackgroundHovered  = BackgroundRaised
	ActionPrimaryBackgroundFocused  = Accent
	ActionPrimaryBackgroundSelected = utils_tty.NewColor(0.0, 0.0, 0.0, 0.0)
	ActionSecondaryBackgroundBase   = utils_tty.NewColor(0.0, 0.0, 0.0, 0.0)
	ActionDestructiveBackgroundBase = FeedbackError

	TextActionSecondaryBase       = TextMuted
	TextActionSecondaryHovered    = MainColor
	TextActionDestructiveBase     = AccentText
	TextActionDestructiveDisabled = TextMuted

	// Gooey-compatible role tokens. These names are identical to the CSS
	// variables used by the Web UI.
	HeaderColor      = MainColor
	HeaderBackground = MainBackground

	AsideColor      = MainColor
	AsideBackground = BackgroundRaised

	ArticleColor      = MainColor
	ArticleBackground = BackgroundRaised
	ArticleBorder     = BorderBase

	FooterColor      = MainColor
	FooterBackground = BackgroundComposer

	ButtonDefaultColor       = MainColor
	ButtonDefaultBackground  = BackgroundRaised
	ButtonDefaultBorder      = BorderBase
	ButtonFocusColor         = Accent
	ButtonFocusBackground    = BackgroundRaisedHigh
	ButtonDisabledColor      = TextMuted
	ButtonDisabledBackground = MainBackground

	InputDefaultColor       = MainColor
	InputDefaultBackground  = MainBackground
	InputDefaultBorder      = BorderBase
	InputFocusColor         = Accent
	InputFocusBackground    = MainBackground
	InputHoverColor         = Accent
	InputPressedColor       = Accent
	InputSelectedColor      = Accent
	InputInvalidColor       = FeedbackError
	InputInvalidBackground  = MainBackground
	InputActiveBackground   = Accent
	InputDisabledColor      = TextMuted
	InputDisabledBackground = MainBackground

	PalettePrimaryColor        = AccentText
	PalettePrimaryBackground   = Accent
	PaletteSecondaryColor      = AccentText
	PaletteSecondaryBackground = FeedbackError

	Palette1Color  = Accent
	Palette2Color  = FeedbackThinking
	Palette3Color  = FeedbackInfo
	Palette4Color  = FeedbackSuccess
	Palette5Color  = FeedbackWarning
	Palette6Color  = FeedbackError
	Palette7Color  = hex("#e5c07b")
	Palette8Color  = hex("#828bb8")
	Palette9Color  = hex("#4fd6be")
	Palette10Color = hex("#c53b53")

	// Diff colors.
	DiffTextAdded      = hex("#4fd6be")
	DiffTextRemoved    = hex("#c53b53")
	DiffTextContext    = hex("#828bb8")
	DiffTextHunkHeader = hex("#828bb8")

	DiffBackgroundAdded   = hex("#20303b")
	DiffBackgroundRemoved = hex("#37222c")
	DiffBackgroundContext = BackgroundRaised

	DiffHighlightAdded   = hex("#b8db87")
	DiffHighlightRemoved = hex("#e26a75")

	DiffLineNumberText              = hex("#8f8f8f")
	DiffLineNumberBackgroundAdded   = hex("#1b2b34")
	DiffLineNumberBackgroundRemoved = hex("#2d1f26")

	// Syntax colors.
	SyntaxComment     = TextMuted
	SyntaxKeyword     = FeedbackThinking
	SyntaxFunction    = Accent
	SyntaxVariable    = FeedbackError
	SyntaxString      = FeedbackSuccess
	SyntaxNumber      = FeedbackWarning
	SyntaxType        = hex("#e5c07b")
	SyntaxOperator    = FeedbackInfo
	SyntaxPunctuation = MainColor

	// Markdown colors.
	MarkdownText            = MainColor
	MarkdownHeading         = FeedbackThinking
	MarkdownLink            = Accent
	MarkdownLinkText        = FeedbackInfo
	MarkdownCode            = FeedbackSuccess
	MarkdownBlockQuote      = hex("#e5c07b")
	MarkdownEmphasis        = hex("#e5c07b")
	MarkdownStrong          = FeedbackWarning
	MarkdownHorizontalRule  = TextMuted
	MarkdownListItem        = Accent
	MarkdownListEnumeration = FeedbackInfo
	MarkdownImage           = Accent
	MarkdownImageText       = FeedbackInfo
	MarkdownCodeBlock       = MainColor

	// Dialog colors.
	DialogBackgroundBase                 = ArticleBackground
	DialogBackgroundActionPrimaryHovered = BackgroundRaisedHigh
)
