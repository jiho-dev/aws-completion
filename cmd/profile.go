package cmd

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/vaughan0/go-ini"
)

func listProfiles() []string {
	// Make sure the config file exists
	config := path.Join(awsDir, "config")

	if _, err := os.Stat(config); os.IsNotExist(err) {
		fmt.Printf("No credentials file found at: %s", config)
		os.Exit(1)
	}

	file, _ := ini.LoadFile(config)
	profiles := make([]string, 0)

	for key, _ := range file {
		if key == "default" {
			profiles = append(profiles, key)
		} else if strings.HasPrefix(key, "profile") {
			k := strings.Split(key, " ")
			if len(k) >= 2 {
				profiles = append(profiles, k[1])
			}
		}
	}

	sort.Strings(profiles)

	return profiles
}
