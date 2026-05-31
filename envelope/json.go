package envelope

//
// ────────────────────────────────────────
// json encoding.
//

import "encoding/json"

// MarshalIndent encodes an envelope as indented JSON.
func MarshalIndent(env *Envelope) ([]byte, error) {
	return json.MarshalIndent(env, "", "  ")
}

// MarshalText encodes an envelope as indented JSON text.
func MarshalText(env *Envelope) (string, error) {
	data, err := MarshalIndent(env)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
