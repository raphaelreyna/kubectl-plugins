package main

import "k8s.io/metrics/pkg/apis/metrics/v1beta1"

type line struct {
	namespace     string
	name          string
	cpu           int64
	cpuPercentage float64
	mem           int64
	memPercentage float64
}

func linesFromPodMetrics(nodeCPU, nodeMem float64, pmList *v1beta1.PodMetricsList, validPods map[string]struct{}) []*line {
	lines := make([]*line, 0)
	for _, pm := range pmList.Items {
		if _, ok := validPods[pm.Name]; !ok {
			continue
		}

		cpu := int64(0)
		mem := int64(0)
		for _, c := range pm.Containers {
			cpu += c.Usage.Cpu().MilliValue()
			mem += c.Usage.Memory().Value()
		}
		cpuf := float64(cpu)
		memf := float64(mem)

		lines = append(lines, &line{
			namespace:     pm.Namespace,
			name:          pm.Name,
			cpu:           cpu,
			cpuPercentage: (cpuf / float64(nodeCPU)) * 100.0,
			mem:           mem,
			memPercentage: (memf / nodeMem) * 100.0,
		})
	}

	return lines
}

type cpuSortableLines []*line

func (c cpuSortableLines) Len() int {
	return len(c)
}

func (c cpuSortableLines) Less(i, j int) bool {
	return c[i].cpu > c[j].cpu
}

func (c cpuSortableLines) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type memSortableLines []*line

func (c memSortableLines) Len() int {
	return len(c)
}

func (c memSortableLines) Less(i, j int) bool {
	return c[i].mem > c[j].mem
}

func (c memSortableLines) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type totalSortableLines []*line

func (c totalSortableLines) Len() int {
	return len(c)
}

func (c totalSortableLines) Less(i, j int) bool {
	return (c[i].cpu + c[i].mem) > (c[j].cpu + c[j].mem)
}

func (c totalSortableLines) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
