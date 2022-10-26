# External KEDA scaler for [DRONE CI](https://www.drone.io/)

An External Scaler for Drone Ci

## Influence Scheduler to remove the idle runners during downscale. 
- Annotates the idle drone runner pods with the [controller.kubernetes.io/pod-deletion-cost](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/#pod-deletion-cost) annotations to influence scheduler during downscaling. 
- Thus idle pods gets removed when HPA downscales. 
- This requires Kubnertes version  v1.22 or higher



## supported ScaledObject properties 
- `droneserver` - url for the drone server
- `dronetoken` - token for drone, this will be moved to [TriggerAuthentication](https://keda.sh/docs/1.4/concepts/authentication/#re-use-credentials-and-delegate-auth-with-triggerauthentication) soon [WIP]
- `minimumPendingBuilds` - miminum threshold for pending builds that can be in pending state that autoscaler can ignore

Triggers Supported :
- [External](https://keda.sh/docs/2.8/scalers/external/)
- [ExternlPush](https://keda.sh/docs/2.8/scalers/external-push/) - WIP 

## Example [see here](./deploy/sample-scaledobjref.yaml)



## TODOS

1. Trigger Authentication for drone token
2. Support Push Scaler, need to implment a GRPC service to listen for manually scaling drone CI 
3. Unit tests 
