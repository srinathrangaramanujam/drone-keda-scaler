package k8s

import (
	"context"
	"strconv"
	"sync"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	config           *rest.Config
	dynamicClientSet dynamic.Interface
}

/*
for the given set of pods, set the annotation
*/

func (kc K8sClient) UpdateAnnotation(pods map[string]bool, namespace string, pod_del_cost int) error {
	apiGroupVersion, err := schema.ParseGroupVersion("v1")
	if err != nil {
		return err
	}

	podGVR := apiGroupVersion.WithResource("pods")
	var wg sync.WaitGroup
	for pod := range pods {
		wg.Add(1)
		func(pod string) {
			res, err := kc.dynamicClientSet.Resource(podGVR).Namespace(namespace).Get(context.Background(), pod, v1.GetOptions{})
			if err != nil {
			}

			currAnnotation := res.GetAnnotations()
			if _, ok := currAnnotation["controller.kubernetes.io/pod-deletion-cost"]; ok {
				return
			}
			currAnnotation["controller.kubernetes.io/pod-deletion-cost"] = strconv.Itoa(pod_del_cost)
			res.SetAnnotations(currAnnotation)
			_, err = kc.dynamicClientSet.Resource(podGVR).Namespace(namespace).Update(context.Background(), res, v1.UpdateOptions{})
			if err != nil {
			}

		}(pod)
	}
	wg.Wait()
	return nil
}

func GetK8sMetadata(k8s_configFile string) (K8sClient, error) {
	var config *rest.Config
	var err error

	if k8s_configFile != "" {
		config, err = clientcmd.BuildConfigFromFlags("", k8s_configFile)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		return K8sClient{}, err
	}
	config.QPS = 10
	config.Burst = 20
	dynamicClientSet, err := dynamic.NewForConfig(config)
	if err != nil {
		return K8sClient{}, err
	}
	return K8sClient{
		config:           config,
		dynamicClientSet: dynamicClientSet,
	}, nil
}
