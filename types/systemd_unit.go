package types

import (
	"strings"
)

const (
	JoinControlPlaneUnitTemplate = "[Unit]\nDescription=init k8s\nAfter=cpk8s.service\nRequires=docker.service\nConditionFileExists=!/etc/kubernetes/kubelet.conf\n[Service]\nType=oneshot\nUser=root\nEnvironment=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/opt/bin:/opt/bin/\nExecStart=/usr/local/bin/kubeadm join %s --config %s \n[Install]\nWantedBy=multi-user.target\n"
	JoinWorkerUnitTemplate       = "[Unit]\nDescription=init k8s\nAfter=cpk8s.service\nRequires=docker.service\nConditionFileExists=!/etc/kubernetes/kubelet.conf\n[Service]\nType=oneshot\nUser=root\nEnvironment=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/opt/bin:/opt/bin/\nExecStart=/usr/local/bin/kubeadm join %s --config %s \n[Install]\nWantedBy=multi-user.target\n"
	InitUnitTemplate             = "[Unit]\nDescription=init k8s\nAfter=cpk8s.service\nRequires=docker.service\nConditionFileExists=!/etc/kubernetes/kubelet.conf\n[Service]\nType=oneshot\nUser=root\nEnvironment=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/opt/bin:/opt/bin/\nExecStart=/usr/local/bin/kubeadm init %s --config %s \n[Install]\nWantedBy=multi-user.target\n"
	KubeadmConfigPath            = "/etc/kubernetes/kubeadm.yaml"
)

type Dropin struct {
	Name    string
	Content string
}

type ServiceUnit struct {
	Content string
	Dropins []Dropin
	Enabled bool
	Name    string
}

func GetCommandsDropins(preKubeadmCommand []string, postKubeadminCommand []string) []Dropin {
	if len(preKubeadmCommand) == 0 && len(postKubeadminCommand) == 0 {
		return []Dropin{}
	}
	builder := strings.Builder{}
	builder.WriteString("[Service]\n")
	for _, command := range preKubeadmCommand {
		builder.WriteString("ExecStartPre=" + command + "\n")
	}
	for _, command := range postKubeadminCommand {
		builder.WriteString("ExecStartPost=" + command + "\n")
	}
	return []Dropin{
		{
			Name:    "10-commands.conf",
			Content: builder.String(),
		},
	}
}
