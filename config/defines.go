package config

type Config struct {
	CrestronIP string
	Port       string
	Password   string
	Devices    []GenericDevice
}

type GenericDevice struct {
	Name        string
	Room        string
	ID          int
	Up          int
	Down        int
	Type        string
	CheckBool   bool
	CTRLOFF     []int
	CTRLON      []int
	StateString string
}
