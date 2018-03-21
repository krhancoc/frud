package config

type Configuration struct {
	Context  *Context       `json:"context"`
	Database *Database      `json:"database"`
	Manager  *ManagerConfig `json:"manager"`
}

func (conf *Configuration) Validate() error {
	err := conf.Manager.Validate()
	if err != nil {
		return err
	}
	return nil
}
