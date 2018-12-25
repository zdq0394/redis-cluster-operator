package controller

import (
	"context"
	"fmt"
	"time"

	objectRuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// Controller is the object that will implement the different kinds of controllers that will be running
// on the application.
type Controller interface {
	// Run runs the controller, it receives a channel that when receiving a signal it will stop the controller,
	// Run will block until it's stopped.
	Run(stopper <-chan struct{}) error
}

// CRD is the custom resource definition.
type CRD interface {
	GetListerWatcher() cache.ListerWatcher
	GetObject() objectRuntime.Object
	Initialize() error
}

// Handler knows how to handle the received resources from a kubernetes cluster.
type Handler interface {
	Add(context.Context, objectRuntime.Object) error
	Delete(context.Context, string) error
}

// SimpleController implements Controller interface
type SimpleController struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
	handler  Handler
	cfg      *Config
}

// NewSimpleController create an instance of Controller
func NewSimpleController(cfg *Config, watchedCRD CRD, handler Handler) Controller {
	cfg.setDefaults()
	// queue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	resourceEventHandlerFuncs := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}

	// indexer and informer
	indexer, informer := cache.NewIndexerInformer(
		watchedCRD.GetListerWatcher(),
		watchedCRD.GetObject(),
		0,
		resourceEventHandlerFuncs,
		cache.Indexers{})

	return &SimpleController{
		informer: informer,
		indexer:  indexer,
		queue:    queue,
		handler:  handler,
	}
}

func (c *SimpleController) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.processWatchedResource(context.Background(), key.(string))
	c.handleErr(err, key)
	return true
}

func (c *SimpleController) processWatchedResource(ctx context.Context, key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		return err
	}
	if !exists {
		return c.handler.Delete(ctx, key)
	}
	return c.handler.Add(ctx, obj.(objectRuntime.Object))
}

func (c *SimpleController) runWorker() {
	for c.processNextItem() {
	}
}

func (c *SimpleController) handleErr(err error, key interface{}) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < 5 {
		fmt.Printf("Error syncing pod %v: %v", key, err)
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	runtime.HandleError(err)
	fmt.Printf("Dropping pod %q out of the queue: %v", key, err)
}

// Run will list and watch the resource and process it.
func (c *SimpleController) Run(stopper <-chan struct{}) error {
	defer runtime.HandleCrash()

	defer c.queue.ShutDown()
	fmt.Println(c.cfg.Name, " Starts...")

	go c.informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return nil
	}

	for i := 0; i < c.cfg.ConcurrentWorkers; i++ {
		go wait.Until(c.runWorker, time.Second, stopper)
	}

	<-stopper
	fmt.Printf("Stopping controller")
	return nil
}
