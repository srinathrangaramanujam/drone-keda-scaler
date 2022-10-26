package drone

import "github.com/drone/drone-go/drone"

type DroneInfo struct {
	DroneHost   string
	droneToken  string
	droneClient drone.Client
}

type CIRunInfo struct {
	PendingBuilds int
	RunningBuilds int
	BuildWorkers  []BuildWorker
}

type BuildWorker struct {
	BuildID string
	Machine string
}
