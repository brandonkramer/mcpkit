package present

//
// ────────────────────────────────────────
// text previews.
//

import "github.com/brandonkramer/mcpkit/envelope"

// TextPreview holds a bounded text preview for tool result data.
type TextPreview struct {
	Preview   string
	Truncated bool
}

// RedactText returns a preview and the text to expose in tool result data.
// When truncated, cleared is empty so callers can omit the full field from payloads.
func RedactText(text string, limit int) (preview TextPreview, cleared string) {
	preview.Preview, preview.Truncated = envelope.PreviewField(text, limit)
	if preview.Truncated {
		return preview, ""
	}
	return preview, text
}
