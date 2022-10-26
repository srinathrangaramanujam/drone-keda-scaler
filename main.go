package main

import (
	"context"
	pb "drone-keda-scaler/externalscaler"
	"drone-keda-scaler/services/drone"
	"drone-keda-scaler/services/k8s"
	"fmt"
	"net"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type DroneScalarServer struct {
	pb.UnimplementedExternalScalerServer
}

/*

1. complete the scaling algo -- done
2. enable token auth
3. enable logging -- done
4. warm up time and warm dowm time,
5. add unit tests
6. code clean up - done
7. Add Environment config
8. add annotation and build up for descheduling - done
9. Documentation
*/

var log *logrus.Logger

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log = createLogger()

	// test()
	grpcServer := grpc.NewServer()
	listerner, err := net.Listen("tcp", fmt.Sprintf(":%v", getEnv("GRPC_LISTEN_PORT", "6000")))
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	pb.RegisterExternalScalerServer(grpcServer, &DroneScalarServer{})
	log.Info("listenting on :6000")
	if err := grpcServer.Serve(listerner); err != nil {
		log.Fatal(err)
	}
}

func createLogger() *logrus.Logger {
	var log = logrus.New()
	log.SetLevel(logrus.TraceLevel)
	log.SetReportCaller(true)
	formatter := &logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return funcName, fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
		},
	}
	log.SetFormatter(formatter)
	return log
}

func test() {

	log.Info("test")
	droneserver := os.Getenv("droneserver")
	dronetoken := os.Getenv("dronetoken")
	srckconfigFile := os.Getenv("configFile")

	droneC := drone.NewDroneMetdata(droneserver, dronetoken)
	droneC.GetPendingBuildCount()
	// workers, _ := droneC.GetWorkingRunners()

	kc, err := k8s.GetK8sMetadata(srckconfigFile)
	if err != nil {
	}

	tst := make(map[string]bool)
	tst["myrunner"] = true
	kc.UpdateAnnotation(tst, "drone", 10)

}

// check if the drone host is active and if the repo is provided check if that repo is available
func (s *DroneScalarServer) IsActive(ctx context.Context, scaleObjRef *pb.ScaledObjectRef) (*pb.IsActiveResponse, error) {
	droneSORef := droneScaleObjRef{scaleObjRef}
	droneInfo := droneSORef.getDroneInfo()

	droneActive, err := droneInfo.IsDroneActive()
	log.Infof("Active status: %v, error: %w", droneActive, err)

	if err != nil {
		log.Fatalf("cannot connect to Drone %v, Error: %v", droneInfo.DroneHost, err.Error())
		return &pb.IsActiveResponse{
			Result: false,
		}, err
	} else {
		return &pb.IsActiveResponse{
			Result: droneActive,
		}, nil
	}
}

func (s *DroneScalarServer) GetMetricSpec(ctx context.Context, scaleObjRef *pb.ScaledObjectRef) (*pb.GetMetricSpecResponse, error) {
	droneSORef := droneScaleObjRef{scaleObjRef}
	spec := pb.MetricSpec{
		MetricName: droneSORef.getMetricName(),
		TargetSize: droneSORef.getMinimumPendingBuilds(),
	}

	return &pb.GetMetricSpecResponse{
		MetricSpecs: []*pb.MetricSpec{&spec},
	}, nil
}

func (s *DroneScalarServer) GetMetrics(ctx context.Context, in *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	droneSORef := droneScaleObjRef{in.ScaledObjectRef}
	droneInfo := droneSORef.getDroneInfo()

	cnt, _, err := droneInfo.GetPendingBuildCount()
	if err != nil {
		log.Fatalf("cannot connect to drone Error: %v", err.Error())
		return nil, err
	}

	m := pb.MetricValue{
		MetricName:  droneSORef.getMetricName(),
		MetricValue: int64(cnt),
	}

	log.Infof("No of pending build: %v\n", m.MetricValue)
	return &pb.GetMetricsResponse{
		MetricValues: []*pb.MetricValue{&m},
	}, nil
}

// need to implement this
func (s *DroneScalarServer) StreamIsActive(scaleObjRef *pb.ScaledObjectRef, svc pb.ExternalScaler_StreamIsActiveServer) error {
	//need to decide this method, if ppl can call the this scaler to force scale
	for {
		select {
		case <-svc.Context().Done():
			return nil
		}
	}
}

func getEnv(key string, defvalue string) string {
	value := os.Getenv(key)

	if len(value) <= 0 {
		value = defvalue
	}

	return value
}
