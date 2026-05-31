package summarize

// Meta carries presentation metadata from prepare/present hooks to summarizers.
//
// Named fields are optional convenience slots for common patterns. Extra is an
// open bag for app-specific context. Summarizers should ignore fields they do
// not use.
type Meta struct {
	Limit    int            // List cap applied during present.
	SinceSeq int            // Optional paging cursor.
	RunID    string         // Optional entity id context.
	Extra    map[string]any // App-specific present→summarize context.
}

// Put stores a value in Extra, allocating the map when needed.
func (m *Meta) Put(key string, value any) {
	if m == nil {
		return
	}
	if m.Extra == nil {
		m.Extra = make(map[string]any)
	}
	m.Extra[key] = value
}

// Get returns a value from Extra.
func (m Meta) Get(key string) any {
	if m.Extra == nil {
		return nil
	}
	return m.Extra[key]
}

// String returns a string value from Extra.
func (m Meta) String(key string) string {
	v, _ := m.Get(key).(string)
	return v
}

// Int returns an int value from Extra.
func (m Meta) Int(key string) int {
	switch v := m.Get(key).(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}
