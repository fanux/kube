# This is a super kubeadm, support master HA with LVS loadbalance!
## On master0 192.168.0.200
```
kubeadm init --config=kubeadm-config.yaml --experimental-upload-certs  
```

## On master1 192.168.0.201
```
kubeadm join 192.168.0.200:6443 --token 9vr73a.a8uxyaju799qwdjv \
    --discovery-token-ca-cert-hash sha256:7c2e69131a36ae2a042a339b33381c6d0d43887e2de83720eff5359e26aec866 \
    --experimental-control-plane \
    --certificate-key f8902e114ef118304e561c3ecd4d0b543adc226b7a07f675f56564185ffe0c07 
```

## On master2 192.168.0.202
```
kubeadm join 192.168.0.200:6443 --token 9vr73a.a8uxyaju799qwdjv \
    --discovery-token-ca-cert-hash sha256:7c2e69131a36ae2a042a339b33381c6d0d43887e2de83720eff5359e26aec866 \
    --experimental-control-plane \
    --certificate-key f8902e114ef118304e561c3ecd4d0b543adc226b7a07f675f56564185ffe0c07  
```

## On your nodes
Join your nodes with local LVS LB 
```
kubeadm join 127.0.0.1:6443 --token 9vr73a.a8uxyaju799qwdjv \
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
