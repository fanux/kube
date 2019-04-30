package kubelet

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog"
)

func (kl *Kubelet) updateNodeRatio() {
	// For debug
	node, err := kl.GetNode()
	if err != nil {
		klog.Errorf(err.Error())
		return
	}
	showNodeResource(node)

	cpuRatio := kl.kubeletConfiguration.CPUAllocationRatio
	memoryRatio := kl.kubeletConfiguration.MemoryAllocationRatio

	cpu := int(float32(kl.machineInfo.NumCores) * cpuRatio)
	if cpu > kl.machineInfo.NumCores {
		kl.machineInfo.NumCores = cpu
	}

	mem := kl.machineInfo.MemoryCapacity * memoryRatio
	if mem > kl.machineInfo.MemoryCapacity {
		kl.machineInfo.MemoryCapacity = mem
	}

	//For debug
	showNodeResource(node)
}

/*
func (kl *Kubelet) updateNodeRatio() {
	// TO set machine info instead node info!
	kl.machineInfo =
	NumCores int `json:"num_cores"`

	// Maximum clock speed for the cores, in KHz.
	CpuFrequency uint64 `json:"cpu_frequency_khz"`

	// The amount of memory (in bytes) in this machine
	MemoryCapacity uint64 `json:"memory_capacity"`



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
*/

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
	cpuC, _ := node.Status.Capacity[v1.ResourceCPU]
	cpuCint, _ := cpuC.AsInt64()
	cpuA, _ := node.Status.Allocatable[v1.ResourceCPU]
	cpuAint, _ := cpuA.AsInt64()

	memC, _ := node.Status.Capacity[v1.ResourceMemory]
	memCint, _ := memC.AsInt64()
	memA, _ := node.Status.Allocatable[v1.ResourceMemory]
	memAint, _ := memA.AsInt64()
	klog.Infof("CPU capacity: %d, CPU allocatable: %d, memory capacity: %d, memory allocatable:%d\n", cpuCint, cpuAint, memCint, memAint)
}
