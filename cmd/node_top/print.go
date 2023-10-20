package main

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

func (c *rootCmd) printLines(out io.Writer, lines []*line) {
	switch c.options.sortBy {
	case "cpu":
		sort.Sort(cpuSortableLines(lines))
	case "mem":
		fallthrough
	case "memory":
		sort.Sort(memSortableLines(lines))
	case "total":
		sort.Sort(totalSortableLines(lines))
	}

	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	tw.Write([]byte("NAME\tCPU(cores)\tCPU(percentage)\tMEMORY(bytes)\tMEMORY(percentage)\n"))
	defer tw.Flush()
	printLine := func(l *line) {
		fmt.Fprintf(tw, "%s\t%dm\t%.2f%%\t%dMi\t%.2f%%\n",
			l.name,
			l.cpu,
			l.cpuPercentage,
			l.mem/(1024*1024),
			l.memPercentage)
	}
	if c.options.order == "asc" {
		for i := len(lines) - 1; i >= 0; i-- {
			if c.options.namespace == "" {
				lines[i].name = lines[i].namespace + "/" + lines[i].name
			}
			printLine(lines[i])
		}
	} else {
		for _, l := range lines {
			if c.options.namespace == "" {
				l.name = l.namespace + "/" + l.name
			}
			printLine(l)
		}
	}
}
