// Package assets embeds static binary assets (fonts, etc.) into the backend binary.
package assets

import _ "embed"

//go:embed fonts/DejaVuSans.ttf
var DejaVuSansRegular []byte

//go:embed fonts/DejaVuSans-Bold.ttf
var DejaVuSansBold []byte
