package config

type Folder struct {
	Path        string   `json:"path" yaml:"path"`
	Destination string   `json:"destination" yaml:"destination"`
	Schedule    string   `json:"schedule" yaml:"schedule"`
	Archiver    string   `json:"archiver" yaml:"archiver"`
	Storages    []string `json:"storages" yaml:"storages"`
	Notifiers   []string `json:"notifiers" yaml:"notifiers"`
}

type Storage struct {
	Type     string         `json:"type" yaml:"type"`
	Settings map[string]any `json:"settings" yaml:"settings"`
}

type Notifier struct {
	Type     string         `json:"type" yaml:"type"`
	Settings map[string]any `json:"settings" yaml:"settings"`
}

type Settings struct {
	Folders   []Folder   `json:"folders" yaml:"folders"`
	Storages  []Storage  `json:"storages" yaml:"storages"`
	Notifiers []Notifier `json:"notifiers" yaml:"notifiers"`
}
