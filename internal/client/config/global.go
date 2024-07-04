package config

type globProfileItem map[string]any

func newGlobProfileItem(path string) globProfileItem {
	return map[string]any{
		"path": path,
	}
}
