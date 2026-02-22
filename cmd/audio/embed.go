package audio

import (
	"embed"
)

//go:embed azan1.mp3
var FS embed.FS

// AzanFile is the embedded filename.
const AzanFile = "azan1.mp3"
