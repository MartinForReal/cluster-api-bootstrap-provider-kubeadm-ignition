package ignition

import (
	"encoding/json"
	ignTypes "github.com/coreos/ignition/config/v3_0/types"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
)

const (
	IngitionFedoraSchemaVersion = "3.0.0"
	ContainerLinuxBaseIgnitionUriHttp = "containerlinux-base.ign"
)

type HttpTemplateBackend struct {
	templateDir    string
	userDataAddr   string
	uploadPath     string
	downloadPath   string
}

func NewHttpTemplateBackend( templateDir,uploadPath ,downloadPath,userDataAddr string) (*HttpTemplateBackend, error) {
	return &HttpTemplateBackend{
		templateDir:    templateDir,
		userDataAddr:   userDataAddr,
		uploadPath:     uploadPath,
		downloadPath:   downloadPath,

	}, nil
}

func (factory *HttpTemplateBackend) getIngitionConfigTemplate(node *Node) (*ignTypes.Config, error) {
	templateConfigUri := factory.getK8sIgnitionFileName(node.Version)

	out := factory.getIngitionBaseConfig()
	out.Ignition.Config = ignTypes.IgnitionConfig{
		Merge: []ignTypes.ConfigReference{
			{
				Source: GetHttpUrl(factory.userDataAddr, factory.templateDir+ContainerLinuxBaseIgnitionUriHttp),
			},
			{
				Source: GetHttpUrl(factory.userDataAddr, factory.templateDir+templateConfigUri),
			},
		},
	}
	//fmt.Println(GetHttpUrl(factory.userDataAddr, ContainerLinuxBaseIgnitionUri))
	return out, nil
}

func (factory *HttpTemplateBackend) getIngitionBaseConfig() *ignTypes.Config {
	return &ignTypes.Config{
		Ignition: ignTypes.Ignition{
			Version: IngitionFedoraSchemaVersion,
		},
	}
}

func (factory *HttpTemplateBackend) getK8sIgnitionFileName(k8sVersion string) string {
	return "k8s-" + k8sVersion + ".ign"
}

func (factory *HttpTemplateBackend) applyConfig(config *ignTypes.Config) (*ignTypes.Config, error) {
	userdata, err := json.Marshal(config)
	if err != nil {
		ignitionLogger.Error(err, "failed to marshal ignition file")
		return nil, err
	}

	fileName := uuid.New().String()+".ign"

	err = ioutil.WriteFile(fileName, userdata, 0666)
	if err != nil {
		ignitionLogger.Error(err,"failed to save ignition file")
	}

	err = UploadFile(factory.userDataAddr,factory.uploadPath,fileName)
	if err !=nil{
		ignitionLogger.Error(err,"Upload file failed")
	}

	if err := os.Remove(fileName); err != nil{
		ignitionLogger.Error(err,"Delete file error")
	}

	out := factory.getIngitionBaseConfig()
	out.Ignition.Config = ignTypes.IgnitionConfig{
		Replace: ignTypes.ConfigReference{
			Source: GetHttpUrl(factory.userDataAddr,factory.downloadPath+fileName),
		},
	}
	return  out,  nil
}