package ignition

import (
	"github.com/minsheng-fintech-corp-ltd/cluster-api-bootstrap-provider-kubeadm-ignition/types"
	"sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha3"
)

const (
	DefaultFileMode = 0644
	DefaultDirMode  = 0755
)

type Node struct {
	Files    []v1alpha3.File
	Services []types.ServiceUnit
	Version  string
}
