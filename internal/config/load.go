package config

func Load(_ string) (Config, error) {
	return Default(), nil
}
