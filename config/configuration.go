package config

// Configuration is the overarching struct for our configuration object
type Configuration struct {
	Context  *Context       `json:"context"`
	Database *Database      `json:"database"`
	Manager  *ManagerConfig `json:"manager"`
}

// Validate will validate the configuration object, it does this by using the validation methods of
// the fields below it.  Drip down validation.
func (conf *Configuration) validate() error {
	err := conf.Manager.validate()
	if err != nil {
		return err
	}
	return nil
}
