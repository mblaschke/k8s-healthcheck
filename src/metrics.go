package main

import (
	"time"
	"log"
	"net/http"
	"k8s.io/api/core/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	kubeNodeInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kube_node_info",
			Help: "Node informations",
		},
		[]string{"node"},
	)

	kubeNodeAge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kube_node_age_seconds",
			Help: "Node age in minutes",
		},
		[]string{"node"},
	)

	kubeNodeStatusCondition = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kube_node_status_condition",
			Help: "Node status conditions",
		},
		[]string{"node", "condition", "status"},
	)


	kubeNodeSpecUnschedulable = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kube_node_spec_unschedulable",
			Help: "Node spec unschedulable",
		},
		[]string{"node"},
	)
)


func initMetrics() {
	prometheus.MustRegister(kubeNodeInfo)
	prometheus.MustRegister(kubeNodeAge)
	prometheus.MustRegister(kubeNodeStatusCondition)
	prometheus.MustRegister(kubeNodeSpecUnschedulable)


	go func() {
		for {
			probeCollect()
			time.Sleep(time.Duration(opts.ScrapeTime) * time.Second)
		}
	}()
}

func startHttpServer() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func probeCollect() {
	nodeList, err := k8sService.NodeList()
	if err != nil {
		panic(err)
	}

	for _, node := range nodeList.Items {
		kubeNodeInfo.With(prometheus.Labels{"node": node.Name}).Set(1)
		kubeNodeAge.With(prometheus.Labels{"node": node.Name}).Set(time.Since(node.CreationTimestamp.Time).Minutes())
		kubeNodeSpecUnschedulable.With(prometheus.Labels{"node": node.Name}).Set(boolFloat64(node.Spec.Unschedulable))

		for _, condition := range node.Status.Conditions {
			kubeNodeStatusCondition.With(prometheus.Labels{"node": node.Name, "condition": string(condition.Type), "status": "true"}).Set(boolFloat64(condition.Status == v1.ConditionTrue))
			kubeNodeStatusCondition.With(prometheus.Labels{"node": node.Name, "condition": string(condition.Type), "status": "false"}).Set(boolFloat64(condition.Status == v1.ConditionFalse))
			kubeNodeStatusCondition.With(prometheus.Labels{"node": node.Name, "condition": string(condition.Type), "status": "unknown"}).Set(boolFloat64(condition.Status == v1.ConditionUnknown))
		}
	}
}

