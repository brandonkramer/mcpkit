package tool

// NormalizeField normalizes one string field in tool args.
func NormalizeField(p *string, norm func(string) (string, error)) error {
	v, err := norm(*p)
	if err != nil {
		return err
	}
	*p = v
	return nil
}

// NormalizeID normalizes one id field in typed tool args.
func NormalizeID[T any](args T, norm func(string) (string, error), idOf func(*T) *string) (T, error) {
	if err := NormalizeField(idOf(&args), norm); err != nil {
		return args, err
	}
	return args, nil
}
