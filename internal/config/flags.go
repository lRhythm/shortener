package config

// String - имплементация flag.Value для типа serverAddress.
func (t *serverAddress) String() string {
	return string(*t)
}

// Set - имплементация flag.Value для типа serverAddress.
func (t *serverAddress) Set(v string) error {
	err := t.validate(v)
	if err != nil {
		return err
	}
	*t = serverAddress(v)
	return nil
}

// String - имплементация flag.Value для типа baseURL.
func (t *baseURL) String() string {
	return string(*t)
}

// Set - имплементация flag.Value для типа baseURL.
func (t *baseURL) Set(v string) error {
	err := t.validate(v)
	if err != nil {
		return err
	}
	*t = baseURL(v)
	return nil
}

// String - имплементация flag.Value для типа fileStoragePath.
func (t *fileStoragePath) String() string {
	return string(*t)
}

// Set - имплементация flag.Value для типа fileStoragePath.
func (t *fileStoragePath) Set(v string) error {
	err := t.validate(v)
	if err != nil {
		return err
	}
	*t = fileStoragePath(v)
	return nil
}

// String - имплементация flag.Value для типа serverAddress.
func (t *databaseDSN) String() string {
	return string(*t)
}

// Set - имплементация flag.Value для типа serverAddress.
func (t *databaseDSN) Set(v string) error {
	*t = databaseDSN(v)
	return nil
}
