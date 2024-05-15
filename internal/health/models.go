package health

import "time"

type Status struct {
	Status      string `json:"status"`
	StartupTime string `json:"startupTime"`
	IpAddress   string `json:"ipAddress"`
	AppVersion  string `json:"appVersion"`
}

func StatusUp(startupTime time.Time, ipAddress string, appVersion string) *Status {
	return status("UP", startupTime, ipAddress, appVersion)
}

func StatusDown(startupTime time.Time, ipAddress string, appVersion string) *Status {
	return status("DOWN", startupTime, ipAddress, appVersion)
}

func status(status string, startupTime time.Time, ipAddress string, appVersion string) *Status {
	return &Status{
		Status:      status,
		StartupTime: startupTime.Format(time.RFC3339),
		IpAddress:   ipAddress,
		AppVersion:  appVersion,
	}
}
