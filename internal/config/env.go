package config

func (t *serverAddress) UnmarshalText(text []byte) error {
	v := string(text)
	if v == "" {
		return nil
	}
	err := t.validate(v)
	if err != nil {
		return err
	}
	*t = serverAddress(v)
	return nil
}

func (t *baseURL) UnmarshalText(text []byte) error {
	v := string(text)
	if v == "" {
		return nil
	}
	err := t.validate(v)
	if err != nil {
		return err
	}
	*t = baseURL(v)
	return nil
}

func (t *fileStoragePath) UnmarshalText(text []byte) error {
	v := string(text)
	if v == "" {
		return nil
	}
	err := t.validate(v)
	if err != nil {
		return err
	}
	*t = fileStoragePath(v)
	return nil
}
