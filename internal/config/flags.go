package config

func (t *serverAddress) String() string {
	return string(*t)
}

func (t *serverAddress) Set(v string) error {
	err := t.validate(v)
	if err != nil {
		return err
	}
	*t = serverAddress(v)
	return nil
}

func (t *baseURL) String() string {
	return string(*t)
}

func (t *baseURL) Set(v string) error {
	err := t.validate(v)
	if err != nil {
		return err
	}
	*t = baseURL(v)
	return nil
}

func (t *fileStoragePath) String() string {
	return string(*t)
}

func (t *fileStoragePath) Set(v string) error {
	err := t.validate(v)
	if err != nil {
		return err
	}
	*t = fileStoragePath(v)
	return nil
}
