package ui

type TextArgument uint8

const (
	Reset     TextArgument = 0
	Bold      TextArgument = 1
	Dim       TextArgument = 2
	Italic    TextArgument = 3
	Underline TextArgument = 4
	Blink     TextArgument = 5

	ForegroundColorBlack   TextArgument = 30
	ForegroundColorRed     TextArgument = 31
	ForegroundColorGreen   TextArgument = 32
	ForegroundColorYellow  TextArgument = 33
	ForegroundColorBlue    TextArgument = 34
	ForegroundColorMagenta TextArgument = 35
	ForegroundColorCyan    TextArgument = 36
	ForegroundColorWhite   TextArgument = 37
)

type TextFormatting struct {
	Args         []TextArgument // Args represents the ANSI escape codes for text formatting.
	LeftPadding  int            // LeftPadding represents the number of spaces to pad on the left.
	RightPadding int            // RightPadding represents the number of spaces to pad on the right.
}
