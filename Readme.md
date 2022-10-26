# External KEDA scaler for DRONE

An External Scaler for Drone Ci


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
