package config

func ParseAfter(raw *Raw, err error) (*Config, error) {
	if err != nil {
		return nil, err
	}

	return Parse(raw)
}
