package main

import (
	"context"

	"github.com/spf13/cobra"
)

func executeContext(ctx context.Context) error {
	root := rootCmd{}

	root.Use = "kubectl node_top"
	root.Args = cobra.ExactArgs(1)
	root.RunE = root.run

	root.options.setFlags(&root.Command)

	return root.ExecuteContext(ctx)
}

type rootCmd struct {
	options options
	cobra.Command
}

func (c *rootCmd) run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	err := c.options.validate(ctx)
	if err != nil {
		return err
	}
	nodeName := args[0]
	o := c.options

	pmList, err := getPodMetrics(ctx, o.namespace)
	if err != nil {
		return err
	}

	validPods, err := getPodsInNode(ctx, o.namespace, nodeName)
	if err != nil {
		return err
	}

	nodeCPU, nodeMem, err := getNodeCapacity(ctx, nodeName)
	if err != nil {
		return err
	}

	lines := linesFromPodMetrics(nodeCPU, nodeMem, pmList, validPods)
	c.printLines(cmd.OutOrStdout(), lines)

	return nil
}
