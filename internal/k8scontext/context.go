package k8scontext

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func getKubeConfigPath() (string, error) {
	var kubeconfigPath *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfigPath = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfigPath = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	if *kubeconfigPath == "" {
		return "", fmt.Errorf("KUBECONFIG environment variable is not set")
	}

	return *kubeconfigPath, nil
}

type clientSetKey struct{}

type metricsClientKey struct{}

type defaultNSKey struct{}

func With(ctx context.Context) (context.Context, error) {
	configPath, err := getKubeConfigPath()
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.LoadFromFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubernetes config: %w", err)
	}

	kcontext := config.Contexts[config.CurrentContext]
	if kcontext == nil {
		return nil, fmt.Errorf("failed to get current context")
	}

	defaultNS := kcontext.Namespace
	if defaultNS == "" {
		defaultNS = "default"
	}
	ctx = context.WithValue(ctx, defaultNSKey{}, defaultNS)

	restConfig, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubernetes config: %w", err)
	}

	cs, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}
	ctx = context.WithValue(ctx, clientSetKey{}, cs)

	mc, err := metrics.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %w", err)
	}
	ctx = context.WithValue(ctx, metricsClientKey{}, mc)

	return ctx, nil
}

func GetClientSet(ctx context.Context) *kubernetes.Clientset {
	return ctx.Value(clientSetKey{}).(*kubernetes.Clientset)
}

func GetMetricsClient(ctx context.Context) *metrics.Clientset {
	return ctx.Value(metricsClientKey{}).(*metrics.Clientset)
}

func GetDefaultNamespace(ctx context.Context) string {
	return ctx.Value(defaultNSKey{}).(string)
}
