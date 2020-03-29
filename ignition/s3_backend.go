package ignition

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	ignTypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/google/uuid"
	"strings"
	"time"
)

var (
	k8sIgnitionUri = map[string]string{
		"v1.15.11": "ignition-config/k8s-v1.15.11.ign",
		"v1.16.8":  "ignition-config/k8s-v1.16.8.ign",
		"v1.17.4":  "ignition-config/k8s-v1.17.4.ign",
		"v1.18.0":  "ignition-config/k8s-v1.18.0.ign",
	}
)

const (
	KubernetesDefaultVersion      = "v1.17.4"
	IngitionSchemaVersion         = "2.2.0"
	ContainerLinuxBaseIgnitionUri = "ignition-config/containerlinux-base.ign"
)

func NewS3TemplateBackend(userdataDir string, userDataBucket string) (*S3TemplateBackend, error) {
	session, err := session.NewSession()
	if err != nil {
		ignitionLogger.Error(err, "failed to initialize s3 session")
		return nil, err
	}
	return &S3TemplateBackend{
		userdataDir:    userdataDir,
		userDataBucket: userDataBucket,
		session:        session,
	}, nil
}

type S3TemplateBackend struct {
	userdataDir    string
	userDataBucket string
	session        *session.Session
}

func (factory *S3TemplateBackend) getIngitionConfigTemplate(node *Node) (*ignTypes.Config, error) {
	templateConfigUri, ok := k8sIgnitionUri[node.Version]
	if !ok {
		err := errors.New("kubernetes version is not supported.")
		ignitionLogger.Error(err, "kubernetes version is not supported.")
		templateConfigUri = k8sIgnitionUri[KubernetesDefaultVersion]
	}

	out := factory.getIngitionBaseConfig()
	out.Ignition.Config = ignTypes.IgnitionConfig{
		Append: []ignTypes.ConfigReference{
			{
				Source: GetS3Url(factory.userDataBucket, ContainerLinuxBaseIgnitionUri),
			},
			{
				Source: GetS3Url(factory.userDataBucket, templateConfigUri),
			},
		},
	}
	return out, nil
}

func (factory *S3TemplateBackend) applyConfig(config *ignTypes.Config) (*ignTypes.Config, error) {
	userdata, err := json.Marshal(config)
	if err != nil {
		ignitionLogger.Error(err, "failed to marshal ignition file")
		return nil, err
	}

	uploader := s3manager.NewUploader(factory.session)
	filePath := strings.Join([]string{factory.userdataDir, uuid.New().String()}, "/")
	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:         bytes.NewReader(userdata),
		Bucket:       aws.String(factory.userDataBucket),
		Expires:      aws.Time(time.Now().Add(time.Hour * 168)),
		Key:          aws.String(filePath),
		StorageClass: aws.String(s3.StorageClassIntelligentTiering),
	})
	if err != nil {
		ignitionLogger.Error(err, "failed to upload ignition file to bucket")
		return nil, err
	}
	out := factory.getIngitionBaseConfig()
	out.Ignition.Config = ignTypes.IgnitionConfig{
		Replace: &ignTypes.ConfigReference{
			Source: GetS3Url(factory.userDataBucket, filePath),
		},
	}
	return out, nil
}

func (factory *S3TemplateBackend) getIngitionBaseConfig() *ignTypes.Config {
	return &ignTypes.Config{
		Ignition: ignTypes.Ignition{
			Version: IngitionSchemaVersion,
		},
	}
}
