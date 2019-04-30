# This is for CPU allocation ratio
Actual CPU utilization is much lower then pod CPU request quota.

# Using kubelet config file, 
default is: `/var/lib/kubelet/config.yaml`
```
cpuAllocationRatio: 3.5   # float32
memoryAllocationRatio: 2  # must uint64
```

# Usage 
Download the kubelet bin file.

```
cp kubelet /usr/bin
```

```
kubectl describe node xx
Capacity:
 cpu:                2
 ephemeral-storage:  20510288Ki
 hugepages-2Mi:      0
 memory:             3881656Ki
 pods:               110
Allocatable:
 cpu:                2
 ephemeral-storage:  18902281390
 hugepages-2Mi:      0
 memory:             3779256Ki
 pods:               110
```

Edit kubelet config file, change the allocation ratio, and restart kubelet.

Then describe node again

# Test
Create a pod witch require 6 (node CPU real cores<6<after ratio cores) cores CPU, if pod run success, it worked.
