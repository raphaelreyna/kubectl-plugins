package main

import (
	"context"
	"errors"

	"github.com/raphaelreyna/kubectl-plugins/internal/k8scontext"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type options struct {
	namespace string
	sortBy    string // "cpu", "mem", "memory", "total"
	order     string // "asc", "desc"

	fs *pflag.FlagSet
}

func (o *options) setFlags(cmd *cobra.Command) {
	if o.fs != nil {
		cmd.Flags().AddFlagSet(o.fs)
		return
	}

	fs := pflag.NewFlagSet("node_top", pflag.ExitOnError)
	fs.StringVarP(&o.namespace, "namespace", "n", "", "namespace to use")
	fs.BoolP("all-namespaces", "A", false, "if present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	fs.StringVar(&o.sortBy, "sort-by", "", "If non-empty, sort nodes list using specified field. The field can be either 'cpu', 'mem', 'memory' or 'total'.")
	fs.StringVar(&o.order, "order", "desc", "Ascending or descending order. Values are 'desc' or 'asc'.")
	o.fs = fs
	cmd.Flags().AddFlagSet(fs)
}

func (o *options) validate(ctx context.Context) error {
	if o.sortBy != "" && o.sortBy != "cpu" && o.sortBy != "mem" && o.sortBy != "memory" && o.sortBy != "total" {
		return errors.New("invalid sort-by")
	}

	if o.order != "desc" && o.order != "asc" && o.order != "" {
		return errors.New("invalid order")
	}

	if all, _ := o.fs.GetBool("all-namespaces"); all {
		o.namespace = ""
	} else if o.namespace == "" {
		o.namespace = k8scontext.GetDefaultNamespace(ctx)
	}

	return nil
}
