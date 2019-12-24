package logging

import "os"

func getHostname() (hostname string) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return
}
