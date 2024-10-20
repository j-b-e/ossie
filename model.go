package main

type Cloud struct {
	Name string
	Env  map[string]string
}

func GetCloud(name string, clouds []Cloud) Cloud {
	for _, v := range clouds {
		if v.Name == name {
			return v
		}
	}
	return Cloud{}
}
