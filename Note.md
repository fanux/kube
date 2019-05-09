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
Create a pod witch require 3 (node CPU real cores<3<after ratio cores) cores CPU, if pod run success, it worked.

> Without ratio, describe node:

kubectl describe node xxx
```
Capacity:
 cpu:                2
 ephemeral-storage:  41151808Ki
 hugepages-1Gi:      0
 hugepages-2Mi:      0
 memory:             3881688Ki
 pods:               110
Allocatable:
 cpu:                2
 ephemeral-storage:  37925506191
 hugepages-1Gi:      0
 hugepages-2Mi:      0
 memory:             3779288Ki
 pods:               110
```

> Create a pod require 3 cpu

```
[root@iZj6c4i2rvm4j0oek6et2zZ ~]# cat cpu.yaml 
apiVersion: v1
kind: Pod
metadata:
  name: cpu-demo
spec:
  containers:
  - name: cpu-demo-ctr
    image: nginx:latest
    resources:
      limits:
        cpu: "4"
      requests:
        cpu: "3"
```
kubectl create -f cpu.yaml, you can see it pending because of not enough cpu:
```
# kubectl get pod
NAME       READY   STATUS    RESTARTS   AGE
cpu-demo   0/1     Pending   0          6s

# kubectl describe pod cpu-demo
  Warning  FailedScheduling  62s (x2 over 62s)  default-scheduler  0/1 nodes are available: 1 Insufficient cpu.
```

> config cpu ratio

vim /var/lib/kubelet/config.yaml

```
cpuAllocationRatio: 3.5
memoryAllocationRatio: 1
```

restart kubelet, describe node again, CPU ups to 7
```
systemctl restart kubelet
kubectl describe node xxx
```
```
Capacity:
 cpu:                7
 ephemeral-storage:  41151808Ki
 hugepages-1Gi:      0
 hugepages-2Mi:      0
 memory:             3881688Ki
 pods:               110
Allocatable:
 cpu:                7
 ephemeral-storage:  37925506191
 hugepages-1Gi:      0
 hugepages-2Mi:      0
 memory:             3779288Ki
 pods:               110
```

now pod running:
```
[root@iZj6c4i2rvm4j0oek6et2zZ ~]# kubectl get pod
NAME       READY   STATUS    RESTARTS   AGE
cpu-demo   1/1     Running   0          4m49s
```

# Using Env config
memoryAllocationRatio uint64
cpuAllocationRatio float64

```
[Service]
Environment="cpuAllocationRatio=3.5"
Environment="memoryAllocationRatio=3"
```
