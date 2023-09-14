package runtimedef

func InitRuntime() error {
	if err := initExecuteable(); err != nil {
		return err
	}
	if err := initYHCHome(); err != nil {
		return err
	}
	return nil
}
