package drone

import (
	"context"

	"github.com/drone/drone-go/drone"
	"golang.org/x/oauth2"
)

func NewDroneMetdata(host string, token string) DroneInfo {
	droneMetadata := DroneInfo{}
	droneMetadata.DroneHost = host
	droneMetadata.droneToken = token

	config := new(oauth2.Config)
	http := config.Client(
		context.Background(),
		&oauth2.Token{
			AccessToken: droneMetadata.droneToken,
		},
	)
	droneMetadata.droneClient = drone.NewClient(droneMetadata.DroneHost, http)
	return droneMetadata
}

func (d DroneInfo) GetPendingBuildCount() (int, int, error) {

	pending := 0
	running := 0

	//deduplicate - 1 build -> N stages
	buildIds := make(map[int]string)
	stages, err := d.droneClient.Queue()
	if err != nil {
		return pending, running, err
	}
	for _, stage := range stages {
		if _, ok := buildIds[int(stage.BuildID)]; !ok {
			buildIds[int(stage.BuildID)] = stage.Machine
			switch stage.Status {
			case drone.StatusPending:
				pending++
			case drone.StatusRunning:
				running++
			}
		}
	}
	return pending, running, err
}

// check if Drone server is reachable
func (d DroneInfo) IsDroneActive() (bool, error) {
	_, err := d.droneClient.RepoList()
	if err != nil {
		return false, err
	}
	return true, nil
}
