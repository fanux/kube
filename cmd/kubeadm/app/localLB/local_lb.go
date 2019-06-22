package locallb

import (
	"fmt"

	"github.com/fanux/LVScare/service"
	"github.com/fanux/LVScare/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/staticpod"
)

//LVScare  is
var LVScare Config

//Config is local lb config
type Config struct {
	VirturlServer string // default is 127.0.0.1:6443
	Masters       []string
	Image         string   // default is fanux/lvscare:latest
	Command       []string // [lvscare care --vs 10.103.97.12:6443 --rs 127.0.0.1:8081 --rs 127.0.0.1:8082 --rs 127.0.0.1:8083 --health-path / --health-schem http]
}

func getSealyunLVScarePod() v1.Pod {
	v := make(map[string]v1.Volume)
	t := true
	pod := staticpod.ComponentPod(v1.Container{
		Name:            "kube-sealyun-lvscare",
		Image:           LVScare.Image,
		ImagePullPolicy: v1.PullIfNotPresent,
		Command:         LVScare.Command,
		SecurityContext: &v1.SecurityContext{Privileged: &t},
	}, v)
	pod.Spec.HostNetwork = true
	return pod
}

//LVScareStaticPodToDisk is
func LVScareStaticPodToDisk(manifests string) {
	staticpod.WriteStaticPodToDisk("kube-sealyun-lvscare", manifests, getSealyunLVScarePod())
}

//InitConfig is
func InitConfig(vs string) {
	LVScare.VirturlServer = vs
	LVScare.Command = []string{
		"/usr/bin/lvscare",
		"care",
		"--vs",
		LVScare.VirturlServer,
		"--health-path",
		"/healthz",
		"--health-schem",
		"https",
	}

	for _, m := range LVScare.Masters {
		LVScare.Command = append(LVScare.Command, "--rs", m)
	}

	fmt.Printf("lvscare command is: %s\n", LVScare.Command)
}

//CreateLVSFirstTime is
func CreateLVSFirstTime() {
	vs := LVScare.VirturlServer
	rs := LVScare.Masters

	lvs, err := service.BuildLvscare(vs, rs)
	if err != nil {
		fmt.Println(err)
	}

	//check virturl server
	service, _ := lvs.GetVirtualServer()
	if service == nil {
		lvs.CreateVirtualServer()
	}

	//check real server
	//lvs.CheckRealServers("/healthz", "https")

	for _, r := range rs {
		rip, rport := utils.SplitServer(r)
		if rip == "" || rport == "" {
			fmt.Println("real server ip and port is null")
		}
		lvs.AddRealServer(rip, rport)
	}

	fmt.Println("creat ipvs first time", vs, rs)
}

//CreateLocalLB is
func CreateLocalLB(vs string) {
	InitConfig(vs)
	CreateLVSFirstTime()
	//LVScareStaticPodToDisk(manifests)
}
