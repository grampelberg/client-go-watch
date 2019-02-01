package main

import (
	"context"
	"flag"
	"time"

	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

const resyncPeriod = 10 * time.Minute

func getConfig() clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Set("v", "10")
	flag.Parse()

	lvl, _ := log.ParseLevel("debug")
	log.SetLevel(lvl)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	cfg, err := getConfig().ClientConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	p := client.CoreV1().Pods("emojivoto")
	plw := &cache.ListWatch{
		ListFunc: func(opt metav1.ListOptions) (runtime.Object, error) {
			return p.List(opt)
		},
		WatchFunc: func(opt metav1.ListOptions) (watch.Interface, error) {
			return p.Watch(opt)
		},
	}

	inf := cache.NewSharedInformer(plw, &apiv1.Pod{}, resyncPeriod)
	inf.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(o interface{}) {
			log.Info("======================================================  ADD")
		},
		DeleteFunc: func(o interface{}) {
			log.Info("======================================================  DELETE")
		},
		UpdateFunc: func(_, o interface{}) {
			log.Info("======================================================  UPDATE")
		},
	})

	ctx, _ := context.WithCancel(context.Background())

	inf.Run(ctx.Done())

	<-ctx.Done()
}
