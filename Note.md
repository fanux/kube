# kubelet with lxcfs
A better kubelet for java applications.

# install lxcfs 
```
wget https://copr-be.cloud.fedoraproject.org/results/ganto/lxd/epel-7-x86_64/00486278-lxcfs/lxcfs-2.0.5-3.el7.centos.x86_64.rpm
yum install lxcfs-2.0.5-3.el7.centos.x86_64.rpm 
systemctl enable lxcfs
systemctl start lxcfs
```

# Test 
```
mkdir /root/static
touch  /root/static/pod.yaml
```

```
[root@izj6c996fvahbz1keqg2iqz ~]# cat static/pod.yaml 
apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
  labels:
    app: myapp
spec:
  containers:
  - name: myapp-container
    image: ubuntu:18.04
    command: ['sh', '-c', 'echo Hello Kubernetes! && sleep 3600']
    resources:
      requests:
        memory: "1Gi"
        cpu: "500m"
      limits:
        memory: "1Gi"
        cpu: "500m"
```

```
./kubelet --feature-gates KubeletSupportLxcfs=true --pod-manifest-path /root/static --cgroup-driver systemd  --runtime-cgroups=/systemd/system.slice --kubelet-cgroups=/systemd/system.slice
```

```
docker exec -it xxxx bash

root@myapp-pod-izj6c996fvahbz1keqg2iqz:/# free -h
              total        used        free      shared  buff/cache   available
Mem:           1.0G        928K        1.0G        720K          0B        1.0G
Swap:            0B          0B          0B
```
As you can see, the memory is 1.0G，the host memory is 4G
