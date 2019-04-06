# This is a super kubeadm, support master HA with LVS loadbalance!
```
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: v1.14.0
controlPlaneEndpoint: "apiserver.cluster.local:6443"
apiServer:
        certSANs:
        - 127.0.0.1
        - apiserver.cluster.local
        - 192.168.0.200
        - 192.168.0.201
        - 192.168.0.202
        - 192.168.0.203
        - 192.168.0.2          # virturl ip
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: "ipvs"
```
## On master0 192.168.0.200
```
echo "192.168.0.200 apiserver.cluster.local" >> /etc/hosts
kubeadm init --config=kubeadm-config.yaml --experimental-upload-certs  
mkdir ~/.kube && cp /etc/kubernetes/admin.conf ~/.kube/config
kubectl apply -f https://docs.projectcalico.org/v3.6/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/1.7/calico.yaml
```

## On master1 192.168.0.201
```
echo "192.168.0.200 apiserver.cluster.local" >> /etc/hosts
kubeadm join 192.168.0.200:6443 --token 9vr73a.a8uxyaju799qwdjv \
    --discovery-token-ca-cert-hash sha256:7c2e69131a36ae2a042a339b33381c6d0d43887e2de83720eff5359e26aec866 \
    --experimental-control-plane \
    --certificate-key f8902e114ef118304e561c3ecd4d0b543adc226b7a07f675f56564185ffe0c07 

sed "s/192.168.0.200/192.168.0.201/g" -i /etc/hosts
```

## On master2 192.168.0.202
```
echo "192.168.0.200 apiserver.cluster.local" >> /etc/hosts
kubeadm join 192.168.0.200:6443 --token 9vr73a.a8uxyaju799qwdjv \
    --discovery-token-ca-cert-hash sha256:7c2e69131a36ae2a042a339b33381c6d0d43887e2de83720eff5359e26aec866 \
    --experimental-control-plane \
    --certificate-key f8902e114ef118304e561c3ecd4d0b543adc226b7a07f675f56564185ffe0c07  

sed "s/192.168.0.200/192.168.0.201/g" -i /etc/hosts
```

## On your nodes
Join your nodes with local LVS LB 
```
echo "192.168.0.2 apiserver.cluster.local" >> /etc/hosts   # using vip
kubeadm join 192.168.0.2:6443 --token 9vr73a.a8uxyaju799qwdjv \
    --master 192.168.0.200:6443 \
    --master 192.168.0.201:6443 \
    --master 192.168.0.202:6443 \
    --discovery-token-ca-cert-hash sha256:7c2e69131a36ae2a042a339b33381c6d0d43887e2de83720eff5359e26aec866
```
Life is much easier!   

# Architecture
```
  +----------+                       +---------------+  virturl server: 127.0.0.1:6443
  | mater0   |<----------------------| ipvs nodes    |    real servers:
  +----------+                      |+---------------+            192.168.0.200:6443
                                    |                             192.168.0.201:6443
  +----------+                      |                             192.168.0.202:6443
  | mater1   |<---------------------+
  +----------+                      |
                                    |
  +----------+                      |
  | mater2   |<---------------------+
  +----------+
```

Every node config a ipvs for masters LB.

Then run a lvscare as a staic pod to check realserver is aviliable. `/etc/kubernetes/manifests/sealyun-lvscare.yaml`

# [LVScare](https://github.com/sealyun/LVScare)
A lvs for kubernetes masters
