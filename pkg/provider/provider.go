package provider

import (
	"fmt"
	"io"
	"path/filepath"

	"os"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	cloudprovider "k8s.io/cloud-provider"
)

// OutSideCluster allows the controller to be started using a local kubeConfig for testing
var OutSideCluster bool

const (
	//ProviderName is the name of the cloud provider
	ProviderName = "simple"

	//SimpleCloudConfig is the default name of the load balancer config Map
	SimpleCloudConfig = "simple"

	//SimpleClientConfig is the default name of the load balancer config Map
	SimpleClientConfig = "simple"

	//SimpleServicesKey is the key in the ConfigMap that has the services configuration
	SimpleServicesKey = "simple-services"
)

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName, newSimpleCloudProvider)
}

// SimpleCloudProvider - contains all of the interfaces for the cloud provider
type SimpleCloudProvider struct {
	lb cloudprovider.LoadBalancer
}

var _ cloudprovider.Interface = &SimpleCloudProvider{}

func newSimpleCloudProvider(io.Reader) (cloudprovider.Interface, error) {
	ns := os.Getenv("SIMPLE_NAMESPACE")
	cm := os.Getenv("SIMPLE_CONFIG_MAP")

	if cm == "" {
		cm = SimpleCloudConfig
	}

	if ns == "" {
		ns = "default"
	}

	var cl *kubernetes.Clientset
	if !OutSideCluster {
		// This will attempt to load the configuration when running within a POD
		cfg, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("error creating kubernetes client config: %s", err.Error())
		}
		cl, err = kubernetes.NewForConfig(cfg)

		if err != nil {
			return nil, fmt.Errorf("error creating kubernetes client: %s", err.Error())
		}
		// use the current context in kubeconfig
	} else {
		config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
		if err != nil {
			panic(err.Error())
		}
		cl, err = kubernetes.NewForConfig(config)

		if err != nil {
			return nil, fmt.Errorf("error creating kubernetes client: %s", err.Error())
		}
	}
	return &SimpleCloudProvider{
		lb: newLoadBalancer(cl, ns, cm),
	}, nil
}

// Initialize - starts the clound-provider controller
func (p *SimpleCloudProvider) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
	clientset := clientBuilder.ClientOrDie("do-shared-informers")
	sharedInformer := informers.NewSharedInformerFactory(clientset, 0)

	//res := NewResourcesController(c.resources, sharedInformer.Core().V1().Services(), clientset)

	sharedInformer.Start(nil)
	sharedInformer.WaitForCacheSync(nil)
	//go res.Run(stop)
	//go c.serveDebug(stop)
}

// LoadBalancer returns a loadbalancer interface. Also returns true if the interface is supported, false otherwise.
func (p *SimpleCloudProvider) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return p.lb, true
}

// ProviderName returns the cloud provider ID.
func (p *SimpleCloudProvider) ProviderName() string {
	return ProviderName
}
