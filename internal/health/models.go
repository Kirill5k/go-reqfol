package health

import "time"

type Status struct {
	Status      string `json:"status"`
	StartupTime string `json:"startupTime"`
	IpAddress   string `json:"ipAddress"`
	AppVersion  string `json:"appVersion"`
}

func StatusUp(startupTime time.Time, ipAddress string, appVersion string) *Status {
	return &Status{
		Status:      "UP",
		StartupTime: startupTime.Format(time.RFC3339),
		IpAddress:   ipAddress,
		AppVersion:  appVersion,
	}
}
