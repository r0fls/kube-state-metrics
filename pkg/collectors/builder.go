/*
Copyright 2018 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// TODO: rename collector
package collectors

import (
	"strings"

	// 	apps "k8s.io/api/apps/v1beta1"
	// 	autoscaling "k8s.io/api/autoscaling/v2beta1"
	// 	batchv1 "k8s.io/api/batch/v1"
	// 	batchv1beta1 "k8s.io/api/batch/v1beta1"
	// 	extensions "k8s.io/api/extensions/v1beta1"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"k8s.io/api/core/v1"
	// 	"k8s.io/api/policy/v1beta1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kube-state-metrics/pkg/metrics"
	metricsstore "k8s.io/kube-state-metrics/pkg/metrics_store"
	"k8s.io/kube-state-metrics/pkg/options"
)

// Builder helps to build collectors. It follows the builder pattern
// (https://en.wikipedia.org/wiki/Builder_pattern).
type Builder struct {
	kubeClient        clientset.Interface
	namespaces        options.NamespaceList
	opts              *options.Options
	ctx               context.Context
	enabledCollectors options.CollectorSet
	metricWhitelist   map[string]struct{}
	metricBlacklist   map[string]struct{}
}

// NewBuilder returns a new builder.
func NewBuilder(
	ctx context.Context,
	opts *options.Options,
) *Builder {
	return &Builder{
		opts: opts,
		ctx:  ctx,
	}
}

// WithEnabledCollectors sets the enabledCollectors property of a Builder.
func (b *Builder) WithEnabledCollectors(c options.CollectorSet) {
	b.enabledCollectors = c
}

// WithNamespaces sets the namespaces property of a Builder.
func (b *Builder) WithNamespaces(n options.NamespaceList) {
	b.namespaces = n
}

// WithKubeClient sets the kubeClient property of a Builder.
func (b *Builder) WithKubeClient(c clientset.Interface) {
	b.kubeClient = c
}

// WithMetricWhitelist configures the whitelisted metrics to be exposed by the
// collectors build by the Builder
func (b *Builder) WithMetricWhitelist(l map[string]struct{}) {
	b.metricWhitelist = l
}

// WithMetricBlacklist configures the blacklisted metrics to be exposed by the
// collectors build by the Builder
func (b *Builder) WithMetricBlacklist(l map[string]struct{}) {
	b.metricBlacklist = l
}

// Build initializes and registers all enabled collectors.
func (b *Builder) Build() []*Collector {

	collectors := []*Collector{}
	activeCollectorNames := []string{}

	for c := range b.enabledCollectors {
		constructor, ok := availableCollectors[c]
		if ok {
			collector := constructor(b)
			activeCollectorNames = append(activeCollectorNames, c)
			collectors = append(collectors, collector)
		}
		// TODO: What if not ok?
	}

	glog.Infof("Active collectors: %s", strings.Join(activeCollectorNames, ","))

	return collectors
}

var availableCollectors = map[string]func(f *Builder) *Collector{
	// 	"configmaps":               func(b *Builder) *Collector { return b.buildConfigMapCollector() },
	// 	"cronjobs":                 func(b *Builder) *Collector { return b.buildCronJobCollector() },
	// 	"daemonsets":               func(b *Builder) *Collector { return b.buildDaemonSetCollector() },
	// 	"deployments":              func(b *Builder) *Collector { return b.buildDeploymentCollector() },
	// 	"endpoints":                func(b *Builder) *Collector { return b.buildEndpointsCollector() },
	// 	"horizontalpodautoscalers": func(b *Builder) *Collector { return b.buildHPACollector() },
	// 	"jobs":                   func(b *Builder) *Collector { return b.buildJobCollector() },
	// 	"limitranges":            func(b *Builder) *Collector { return b.buildLimitRangeCollector() },
	// 	"namespaces":             func(b *Builder) *Collector { return b.buildNamespaceCollector() },
	// 	"nodes":                  func(b *Builder) *Collector { return b.buildNodeCollector() },
	// 	"persistentvolumeclaims": func(b *Builder) *Collector { return b.buildPersistentVolumeClaimCollector() },
	// 	"persistentvolumes":      func(b *Builder) *Collector { return b.buildPersistentVolumeCollector() },
	// 	"poddisruptionbudgets":   func(b *Builder) *Collector { return b.buildPodDisruptionBudgetCollector() },
	// 	"pods":                   func(b *Builder) *Collector { return b.buildPodCollector() },
	// 	"replicasets":            func(b *Builder) *Collector { return b.buildReplicaSetCollector() },
	// 	"replicationcontrollers": func(b *Builder) *Collector { return b.buildReplicationControllerCollector() },
	// 	"resourcequotas":         func(b *Builder) *Collector { return b.buildResourceQuotaCollector() },
	// 	"secrets":                func(b *Builder) *Collector { return b.buildSecretCollector() },
	"services": func(b *Builder) *Collector { return b.buildServiceCollector() },
	//	"statefulsets":           func(b *Builder) *Collector { return b.buildStatefulSetCollector() },
}

// func (b *Builder) buildPodCollector() *Collector {
// 	genFunc := func(obj interface{}) []*metrics.Metric {
// 		return generatePodMetrics(b.opts.DisablePodNonGenericResourceMetrics, obj)
// 	}
// 	store := metricsstore.NewMetricsStore(genFunc)
// 	reflectorPerNamespace(b.ctx, b.kubeClient, &v1.Pod{}, store, b.namespaces, createPodListWatch)
//
// 	return NewCollector(store)
// }
//
// func (b *Builder) buildCronJobCollector() *Collector {
// 	store := metricsstore.NewMetricsStore(generateCronJobMetrics)
// 	reflectorPerNamespace(b.ctx, b.kubeClient, &batchv1beta1.CronJob{}, store, b.namespaces, createCronJobListWatch)
//
// 	return NewCollector(store)
// }
//
// func (b *Builder) buildConfigMapCollector() *Collector {
// 	store := metricsstore.NewMetricsStore(generateConfigMapMetrics)
// 	reflectorPerNamespace(b.ctx, b.kubeClient, &v1.ConfigMap{}, store, b.namespaces, createConfigMapListWatch)
//
// 	return NewCollector(store)
// }
//
// func (b *Builder) buildDaemonSetCollector() *Collector {
// 	store := metricsstore.NewMetricsStore(generateDaemonSetMetrics)
// 	reflectorPerNamespace(b.ctx, b.kubeClient, &extensions.DaemonSet{}, store, b.namespaces, createDaemonSetListWatch)
//
// 	return NewCollector(store)
// }
//
// func (b *Builder) buildDeploymentCollector() *Collector {
// 	store := metricsstore.NewMetricsStore(generateDeploymentMetrics)
// 	reflectorPerNamespace(b.ctx, b.kubeClient, &extensions.Deployment{}, store, b.namespaces, createDeploymentListWatch)
//
// 	return NewCollector(store)
// }
//
// func (b *Builder) buildEndpointsCollector() *Collector {
// 	store := metricsstore.NewMetricsStore(generateEndpointsMetrics)
// 	reflectorPerNamespace(b.ctx, b.kubeClient, &v1.Endpoints{}, store, b.namespaces, createEndpointsListWatch)
//
// 	return NewCollector(store)
// }
//
// func (b *Builder) buildHPACollector() *Collector {
// 	store := metricsstore.NewMetricsStore(generateHPAMetrics)
// 	reflectorPerNamespace(b.ctx, b.kubeClient, &autoscaling.HorizontalPodAutoscaler{}, store, b.namespaces, createHPAListWatch)
//
// 	return NewCollector(store)
// }
//
// func (b *Builder) buildJobCollector() *Collector {
// 	store := metricsstore.NewMetricsStore(generateJobMetrics)
// 	reflectorPerNamespace(b.ctx, b.kubeClient, &batchv1.Job{}, store, b.namespaces, createJobListWatch)
//
// 	return NewCollector(store)
// }
//
// func (b *Builder) buildLimitRangeCollector() *Collector {
// 	store := metricsstore.NewMetricsStore(generateLimitRangeMetrics)
// 	reflectorPerNamespace(b.ctx, b.kubeClient, &v1.LimitRange{}, store, b.namespaces, createLimitRangeListWatch)
//
// 	return NewCollector(store)
// }
func (b *Builder) buildServiceCollector() *Collector {
	filteredMetricFamilies := filterMetricFamilies(b.metricWhitelist, b.metricBlacklist, serviceMetricFamilies)

	store := metricsstore.NewMetricsStore(
		composeMetricGenFuncs(filteredMetricFamilies),
	)
	reflectorPerNamespace(b.ctx, b.kubeClient, &v1.Service{}, store, b.namespaces, createServiceListWatch)

	return NewCollector(store)
}

// func (b *Builder) buildNodeCollector() *Collector {
// 	genFunc := func(obj interface{}) []*metrics.Metric {
// 		return generateNodeMetrics(b.opts.DisableNodeNonGenericResourceMetrics, obj)
// 	}
//
// 	return newCollector(store)
// }

// composeMetricGenFuncs takes a slice of metric families and returns a function
// that composes their metric generation functions into a single one.
func composeMetricGenFuncs(families []metrics.MetricFamily) func(obj interface{}) []*metrics.Metric {
	funcs := []func(obj interface{}) []*metrics.Metric{}

	for _, f := range families {
		funcs = append(funcs, f.GenerateFunc)
	}

	return func(obj interface{}) []*metrics.Metric {
		metrics := []*metrics.Metric{}

		for _, f := range funcs {
			metrics = append(metrics, f(obj)...)
		}

		return metrics
	}
}

// filterMetricFamilies takes a white- and a blacklist and a slice of metric
// families and returns a filtered slice.
func filterMetricFamilies(white, black map[string]struct{}, families []metrics.MetricFamily) []metrics.MetricFamily {
	if len(white) != 0 && len(black) != 0 {
		panic("Whitelist and blacklist are both set. They are mutually exclusive, only one of them can be set.")
	}

	filtered := []metrics.MetricFamily{}

	if len(white) != 0 {
		for _, f := range families {
			if _, whitelisted := white[f.Name]; whitelisted {
				filtered = append(filtered, f)
			}
		}

		return filtered
	}

	for _, f := range families {
		if _, blacklisted := black[f.Name]; !blacklisted {
			filtered = append(filtered, f)
		}
	}

	return filtered
}

//
// func (b *Builder) buildStatefulSetCollector() *Collector {
// 	store := metricsstore.NewMetricsStore(generateStatefulSetMetrics)
// 	reflectorPerNamespace(b.ctx, b.kubeClient, &apps.StatefulSet{}, store, b.namespaces, createStatefulSetListWatch)
//
// 	return newCollector(store)
// }

// reflectorPerNamespace creates a Kubernetes client-go reflector with the given
// listWatchFunc for each given namespace and registers it with the given store.
func reflectorPerNamespace(
	ctx context.Context,
	kubeClient clientset.Interface,
	expectedType interface{},
	store cache.Store,
	namespaces []string,
	listWatchFunc func(kubeClient clientset.Interface, ns string) cache.ListWatch,
) {
	for _, ns := range namespaces {
		lw := listWatchFunc(kubeClient, ns)
		reflector := cache.NewReflector(&lw, expectedType, store, 0)
		go reflector.Run(ctx.Done())
	}
}
