package provider

import (
	"context"
	"encoding/json"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Services functions - once the service data is taken from the configMap, these functions will interact with the data

func (s *simpleServices) addService(newSvc services) {
	s.Services = append(s.Services, newSvc)
}

func (s *simpleServices) findService(UID string) *services {
	for x := range s.Services {
		if s.Services[x].UID == UID {
			return &s.Services[x]
		}
	}
	return nil
}

func (s *simpleServices) delServiceFromUID(UID string) *simpleServices {
	// New Services list
	updatedServices := &simpleServices{}
	// Add all [BUT] the removed service
	for x := range s.Services {
		if s.Services[x].UID != UID {
			updatedServices.Services = append(updatedServices.Services, s.Services[x])
		}
	}
	// Return the updated service list (without the mentioned service)
	return updatedServices
}

// ConfigMap functions - these wrap all interactions with the kubernetes configmaps

func (k *simpleLoadBalancerManager) GetServices(cm *v1.ConfigMap) (svcs *simpleServices, err error) {
	// Attempt to retrieve the config map
	b := cm.Data[SimpleServicesKey]
	// Unmarshall raw data into struct
	err = json.Unmarshal([]byte(b), &svcs)
	return
}

func (k *simpleLoadBalancerManager) GetConfigMap(ctx context.Context, cm, nm string) (*v1.ConfigMap, error) {
	// Attempt to retrieve the config map
	return k.kubeClient.CoreV1().ConfigMaps(nm).Get(ctx, k.cloudConfigMap, metav1.GetOptions{})
}

func (k *simpleLoadBalancerManager) CreateConfigMap(ctx context.Context, cm, nm string) (*v1.ConfigMap, error) {
	// Create new configuration map in the correct namespace
	newConfigMap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      k.cloudConfigMap,
			Namespace: nm,
		},
	}
	// Return results of configMap create
	return k.kubeClient.CoreV1().ConfigMaps(nm).Create(ctx, &newConfigMap, metav1.CreateOptions{})
}

func (k *simpleLoadBalancerManager) UpdateConfigMap(ctx context.Context, cm *v1.ConfigMap, s *simpleServices) (*v1.ConfigMap, error) {
	// Create new configuration map in the correct namespace

	// If the cm.Data / cm.Annotations haven't been initialised
	if cm.Data == nil {
		cm.Data = map[string]string{}
	}
	if cm.Annotations == nil {
		cm.Annotations = map[string]string{}
		cm.Annotations["provider"] = ProviderName
	}

	// Set ConfigMap data
	b, _ := json.Marshal(s)
	cm.Data[SimpleServicesKey] = string(b)

	// Return results of configMap create
	return k.kubeClient.CoreV1().ConfigMaps(cm.Namespace).Update(ctx, cm, metav1.UpdateOptions{})
}
