package main

import (
	pb "drone-keda-scaler/externalscaler"
	"drone-keda-scaler/services/drone"
	"strconv"
)

const (
	metricName           = "metricName"
	minimumPendingBuilds = "minimumPendingBuilds"
	droneserver          = "droneserver"
	dronetoken           = ""
)

type droneScaleObjRef struct {
	*pb.ScaledObjectRef
}

func (s droneScaleObjRef) getMetricName() string {
	metricName, ok := s.ScalerMetadata[metricName]
	if ok {
		return metricName
	} else {
		return s.GetName()
	}
}

func (s droneScaleObjRef) getMinimumPendingBuilds() int64 {
	res, ok := s.ScalerMetadata[minimumPendingBuilds]
	if !ok {
		return int64(5)
	}
	val, err := strconv.Atoi(res)
	if err != nil {
		val = 5
	}
	return int64(val)
}

func (s droneScaleObjRef) getDroneInfo() drone.DroneInfo {
	return drone.NewDroneMetdata(s.ScalerMetadata[droneserver], s.ScalerMetadata[dronetoken])
}
