package config

func GetConfig() map[string]string {
	return map[string]string{
		"download_url": "https://github.com/prappo/wordpress-plugin-boilerplate/archive/refs/heads/main.zip",
		"version":      "1.0.0",
		"author":       "prappo",
		"description":  "plugkit is a plugin for creating WordPress plugins",
		"license":      "MIT",
		"textdomain":   "plugkit",
		"domainPath":   "/languages",
	}
}
