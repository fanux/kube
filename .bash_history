cd /go/src/k8s.io/kubernetes/
 git clone https://github.com/fanux/LVScare /go/src/github.com/fanux/LVScare
cd /go/src/k8s.io/kubernetes/
KUBE_GIT_TREE_STATE="clean" KUBE_GIT_VERSION=v1.15.0 KUBE_BUILD_PLATFORMS=linux/amd64 make all WHAT=cmd/kubeadm GOFLAGS=-v
KUBE_GIT_TREE_STATE="clean" KUBE_GIT_VERSION=v1.15.0 KUBE_BUILD_PLATFORMS=linux/amd64 make all WHAT=cmd/kubeadm GOFLAGS=-v
exit
