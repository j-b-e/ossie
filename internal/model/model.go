package model

// Cloud represents an Openstack environment with env vars
type Cloud struct {
	Name string
	Env  map[string]string
}

type Clouds []Cloud

// Select a cloud by name
func (c Clouds) Select(name string) Cloud {
	for _, v := range c {
		if v.Name == name {
			return v
		}
	}
	return Cloud{}
}
