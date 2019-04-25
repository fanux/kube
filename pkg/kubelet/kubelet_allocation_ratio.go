package kubelet

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog"
)

func (kl *Kubelet) updateNodeRatio() {
	node, err := kl.GetNode()
	if err != nil {
		klog.Errorf(err.Error())
		return
	}

	cpuRatio := kl.kubeletConfiguration.CPUAllocationRatio
	memoryRatio := kl.kubeletConfiguration.MemoryAllocationRatio

	showNodeResource(node)
	node.Status.Capacity[v1.ResourceCPU] = ratioResource(cpuRatio, node.Status.Capacity[v1.ResourceCPU])
	node.Status.Allocatable[v1.ResourceCPU] = ratioResource(cpuRatio, node.Status.Allocatable[v1.ResourceCPU])
	node.Status.Capacity[v1.ResourceMemory] = ratioResource(memoryRatio, node.Status.Capacity[v1.ResourceMemory])
	node.Status.Allocatable[v1.ResourceMemory] = ratioResource(memoryRatio, node.Status.Allocatable[v1.ResourceMemory])
	showNodeResource(node)

	kl.syncNodeStatus()
}

func ratioResource(ratio float32, resource resource.Quantity) resource.Quantity {
	i, ok := resource.AsInt64()
	if !ok {
		klog.Errorf("ratio resource failed, resource as int64 failed\n")
		return resource
	}
	ratioRes := float32(i) * ratio

	resource.Set(int64(ratioRes))
	return resource
}

func showNodeResource(node *v1.Node) {
	cpuC, _ := node.Status.Capacity[v1.ResourceCPU].AsInt64()
	cpuA, _ := node.Status.Allocatable[v1.ResourceCPU].AsInt64()

	memC, _ := node.Status.Capacity[v1.ResourceMemory].AsInt64()
	memA, _ := node.Status.Allocatable[v1.ResourceMemory].AsInt64()
	klog.Infof("CPU capacity: %d, CPU allocatable: %d, memory capacity: %d, memory allocatable:%d\n", cpuC, cpuA, memC, memA)
}
