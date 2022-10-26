```go
// another style
func (d DroneInfo) getPendingBuildCount2(useV2 bool, repoInfo *RepoInfo) (int, error) {
	totalPendingBuilds := 0
	var err error
	if useV2 {
		repoBuilds, err := d.droneClient.IncompleteV2()
		if err == nil {
			for _, repo := range repoBuilds {
				if (repoInfo == nil || repoInfo.EqualsRepoBuildStage(repo)) && repo.BuildStarted == int64(0) {
					totalPendingBuilds++
				}
			}
		}
	} else {
		repos, err := d.droneClient.Incomplete()
		if err == nil {
			repoMap := make(map[string]bool)
			var w sync.WaitGroup
			var m sync.Mutex
			repos = getRepos(repoInfo, repos)
			for _, repo := range repos {
				name := fmt.Sprintf("%v/%v", repo.Namespace, repo.Name)
				if _, ok := repoMap[name]; ok {
					continue
				}
				repoMap[name] = true
				w.Add(1)
				go func(repo *drone.Repo) {
					defer w.Done()
					builds, _ := d.droneClient.BuildList(repo.Namespace, repo.Name, drone.ListOptions{Size: 10})
					for _, build := range builds {
						if build.Started == 0 {
							m.Lock()
							totalPendingBuilds++
							m.Unlock()
						}
					}
				}(repo)
			}
			w.Wait()
		}
	}
	return totalPendingBuilds, err
}



func (d DroneInfo) GetWorkingRunners() (map[string]string, error) {
	runningWorkers := make(map[string]string)
	repoBuilds, err := d.droneClient.Incomplete()
	if err != nil {
		return runningWorkers, err
	}

	var wg sync.WaitGroup
	var m sync.Mutex
	start := time.Now()

	for _, repo := range repoBuilds {

		builds, err := d.droneClient.BuildList(repo.Namespace, repo.Name, drone.ListOptions{Size: 10})
		if err != nil {
			return runningWorkers, err
		}
		for _, build := range builds {
			wg.Add(1)
			if build.Started == 0 {
				continue
			}
			b, err := d.droneClient.Build(repo.Namespace, repo.Name, int(build.ID))
			if err != nil {
				return runningWorkers, err
			}
			func(b *drone.Build) {
				defer wg.Done()
				for _, stage := range b.Stages {
					if stage.Started == 0 {
						continue
					}
					m.Lock()
					if _, ok := runningWorkers[stage.Machine]; !ok {
						runningWorkers[stage.Machine] = stage.Machine
					}
					m.Unlock()

				}
			}(b)
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("Binomial took %s", elapsed)

	return runningWorkers, err
}
// import "fmt"

// func (d DroneInfo) GetIncompleteReposV2() (map[string]RepoInfo, int, error) {
// 	repoMap := make(map[string]RepoInfo)
// 	totalPendingBuilds := 0
// 	repoBuilds, err := d.droneClient.IncompleteV2()
// 	if err == nil {
// 		for _, repo := range repoBuilds {
// 			if repo.BuildStarted != int64(0) {
// 				continue
// 			}
// 			totalPendingBuilds++
// 			name := fmt.Sprintf("%v/%v", repo.RepoNamespace, repo.RepoName)
// 			if rInfo, ok := repoMap[name]; ok {
// 				rInfo.PendingBuildCount++
// 				repoMap[name] = rInfo
// 				continue
// 			} else {
// 				repoMap[name] = RepoInfo{RepoName: repo.RepoName, RepoNamespace: repo.RepoNamespace, PendingBuildCount: 1}
// 			}

// 		}
// 	}

// 	return repoMap, totalPendingBuilds, err
// }

// func (d DroneInfo) GetIncompleteReposV2() (map[string]RepoInfo, int, error) {
// 	repoMap := make(map[string]RepoInfo)
// 	totalPendingBuilds := 0
// 	repoBuilds, err := d.droneClient.IncompleteV2()
// 	if err == nil {
// 		for _, repo := range repoBuilds {
// 			if repo.BuildStarted != int64(0) {
// 				continue
// 			}
// 			totalPendingBuilds++
// 			name := fmt.Sprintf("%v/%v", repo.RepoNamespace, repo.RepoName)
// 			if rInfo, ok := repoMap[name]; ok {
// 				rInfo.PendingBuildCount++
// 				repoMap[name] = rInfo
// 				continue
// 			} else {
// 				repoMap[name] = RepoInfo{RepoName: repo.RepoName, RepoNamespace: repo.RepoNamespace, PendingBuildCount: 1}
// 			}

// 		}
// 	}

// 	return repoMap, totalPendingBuilds, err
// }

// //Incomplete() returns lits of repos that have stages thatare in running or pending status. thus a single repo can be returned twic if it have two builds running
// func (d DroneInfo) GetIncompleteRepos() (map[string]RepoInfo, int, error) {
// 	totalPendingBuilds := 0
// 	repoMap := make(map[string]RepoInfo)
// 	repos, err := d.droneClient.Incomplete()
// 	if err == nil {
// 		for _, repo := range repos {
// 			//if repo exists in the map then the repo pending build will already by counted in
// 			name := fmt.Sprintf("%v/%v", repo.Namespace, repo.Name)
// 			if _, ok := repoMap[name]; ok {
// 				continue
// 			}

// 			builds, err := d.droneClient.BuildList(repo.Namespace, repo.Name, drone.ListOptions{Size: 10})
// 			fmt.Print(err)

// 			for _, build := range builds {
// 				fmt.Printf("build id %v\n", build.ID)
// 				if build.Started == 0 {
// 					totalPendingBuilds++
// 				}
// 			}

// 			name := fmt.Sprintf("%v/%v", repo.Namespace, repo.Name)
// 			if _, ok := repoMap[name]; ok {
// 				continue
// 			}
// 			totalPendingBuilds++

// 			func(repo *drone.Repo) {
// 				builds, err := d.droneClient.BuildList(repo.Namespace, repo.Name, drone.ListOptions{Size: 10})
// 				fmt.Print(err)
// 				for _, build := range builds {
// 					if build.Started == 0 {
// 						continue
// 					}
// 				}
// 			}(repo)

// 		}

// 	}
// 	return repoMap1, totalPendingBuilds, err
// }


```
