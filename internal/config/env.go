package config

// UnmarshalText - decode собственного типа serverAddress для вызова пакетом caarlos0/env.
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

// UnmarshalText - decode собственного типа baseURL для вызова пакетом caarlos0/env.
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

// UnmarshalText - decode собственного типа fileStoragePath для вызова пакетом caarlos0/env.
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
