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

// SimpleController implements Controller interface
type SimpleController struct {
	cfg      *Config
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
	handler  Handler
}

// NewSimpleController create an instance of Controller
func NewSimpleController(cfg *Config, crd CRD, handler Handler) Controller {
	if cfg == nil {
		cfg = &Config{}
	}
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
		crd.GetListerWatcher(),
		crd.GetObject(),
		0,
		resourceEventHandlerFuncs,
		cache.Indexers{})

	// Create the SimpleController Instance
	return &SimpleController{
		cfg:      cfg,
		informer: informer,
		indexer:  indexer,
		queue:    queue,
		handler:  handler,
	}
}

// Run will list and watch the resources and then process them.
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

func (c *SimpleController) runWorker() {
	for c.processNextItem() {
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
