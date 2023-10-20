package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/raphaelreyna/kubectl-plugins/internal/k8scontext"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func getNodeCapacity(ctx context.Context, nodeName string) (cpu, mem float64, err error) {
	cs := k8scontext.GetClientSet(ctx)
	if cs == nil {
		return 0, 0, errors.New("failed to get kubernetes clientset")
	}

	nodeInfo, err := cs.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get node info: %w", err)
	}

	return float64(nodeInfo.Status.Capacity.Cpu().MilliValue()),
		float64(nodeInfo.Status.Capacity.Memory().Value()),
		nil
}

func getPodsInNode(ctx context.Context, namespace, nodeName string) (map[string]struct{}, error) {
	cs := k8scontext.GetClientSet(ctx)
	if cs == nil {
		return nil, errors.New("failed to get kubernetes clientset")
	}
	if nodeName == "" {
		return nil, errors.New("node name cannot be empty")
	}

	podsList, err := cs.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	if len(podsList.Items) == 0 {
		return nil, fmt.Errorf("no pods found in node %s", nodeName)
	}

	podMap := make(map[string]struct{}, len(podsList.Items))
	for _, podInfo := range podsList.Items {
		podMap[podInfo.Name] = struct{}{}
	}

	return podMap, nil
}

func getPodMetrics(ctx context.Context, namespace string) (*v1beta1.PodMetricsList, error) {
	mc := k8scontext.GetMetricsClient(ctx)
	if mc == nil {
		return nil, errors.New("failed to get kubernetes metrics client")
	}

	pmList, err := mc.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod metrics: %w", err)
	}
	if len(pmList.Items) == 0 {
		return nil, fmt.Errorf("no pod metrics found")
	}

	return pmList, nil
}
