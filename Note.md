# This is a super kubeadm, support master HA with LVS loadbalance!
```
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: v1.14.0
controlPlaneEndpoint: "apiserver.cluster.local:6443" # apiserver DNS name
apiServer:
        certSANs:
        - 127.0.0.1
        - apiserver.cluster.local
        - 172.20.241.205
        - 172.20.241.206
        - 172.20.241.207
        - 172.20.241.208
        - 10.103.97.1          # virturl ip
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: "ipvs"
ipvs:
        excludeCIDRs: 
        - "10.103.97.1/32" # if you don't add this kube-proxy will clean your ipvs rule(kube-proxy still remove it)
```
## On master0 10.103.97.100
```
echo "10.103.97.100 apiserver.cluster.local" >> /etc/hosts
kubeadm init --config=kubeadm-config.yaml --experimental-upload-certs  
mkdir ~/.kube && cp /etc/kubernetes/admin.conf ~/.kube/config
kubectl apply -f https://docs.projectcalico.org/v3.6/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/1.7/calico.yaml
```

## On master1 10.103.97.101
```
echo "10.103.97.100 apiserver.cluster.local" >> /etc/hosts
kubeadm join 10.103.97.100:6443 --token 9vr73a.a8uxyaju799qwdjv \
    --discovery-token-ca-cert-hash sha256:7c2e69131a36ae2a042a339b33381c6d0d43887e2de83720eff5359e26aec866 \
    --experimental-control-plane \
    --certificate-key f8902e114ef118304e561c3ecd4d0b543adc226b7a07f675f56564185ffe0c07 

sed "s/10.103.97.100/10.103.97.101/g" -i /etc/hosts  # if you don't change this, kubelet and kube-proxy will also using master0 apiserver, when master0 down, everything dead!
```

## On master2 10.103.97.102
```
echo "10.103.97.100 apiserver.cluster.local" >> /etc/hosts
kubeadm join 10.103.97.100:6443 --token 9vr73a.a8uxyaju799qwdjv \
    --discovery-token-ca-cert-hash sha256:7c2e69131a36ae2a042a339b33381c6d0d43887e2de83720eff5359e26aec866 \
    --experimental-control-plane \
    --certificate-key f8902e114ef118304e561c3ecd4d0b543adc226b7a07f675f56564185ffe0c07  

sed "s/10.103.97.100/10.103.97.101/g" -i /etc/hosts
```

## On your nodes
Join your nodes with local LVS LB 
```
echo "10.103.97.1 apiserver.cluster.local" >> /etc/hosts   # using vip
kubeadm join 10.103.97.1:6443 --token 9vr73a.a8uxyaju799qwdjv \
    --master 10.103.97.100:6443 \
    --master 10.103.97.101:6443 \
    --master 10.103.97.102:6443 \
    --discovery-token-ca-cert-hash sha256:7c2e69131a36ae2a042a339b33381c6d0d43887e2de83720eff5359e26aec866
```
Life is much easier!   

# Architecture
```
  +----------+                       +---------------+  virturl server: 127.0.0.1:6443
  | mater0   |<----------------------| ipvs nodes    |    real servers:
  +----------+                      |+---------------+            10.103.97.100:6443
                                    |                             10.103.97.101:6443
  +----------+                      |                             10.103.97.102:6443
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
