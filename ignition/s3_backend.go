package ignition

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	ignTypes "github.com/coreos/ignition/v2/config/v3_0/types"
	"github.com/google/uuid"
)

const (
	IngitionSchemaVersion         = "2.2.0"
	ContainerLinuxBaseIgnitionUri = "ignition-config/containerlinux-base.ign"
)

func NewS3TemplateBackend(userdataDir, templateDir string, userDataBucket string) (*S3TemplateBackend, error) {
	session, err := session.NewSession()
	if err != nil {
		ignitionLogger.Error(err, "failed to initialize s3 session")
		return nil, err
	}
	return &S3TemplateBackend{
		userdataDir:    userdataDir,
		templateDir:    templateDir,
		userDataBucket: userDataBucket,
		session:        session,
	}, nil
}

type S3TemplateBackend struct {
	userdataDir    string
	templateDir    string
	userDataBucket string
	session        *session.Session
}

func (factory *S3TemplateBackend) getIngitionConfigTemplate(node *Node) (*ignTypes.Config, error) {
	svc := s3.New(factory.session)
	_, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(factory.userDataBucket),
		Key:    aws.String(factory.getK8sIgnitionFileName(node.Version)),
	})
	if err != nil {
		ignitionLogger.Error(err, "kubernetes version is not supported.")
		return nil, err
	}
	templateConfigUri := factory.getK8sIgnitionFileName(node.Version)

	out := factory.getIngitionBaseConfig()
	out.Ignition.Config = ignTypes.IgnitionConfig{
		Merge: []ignTypes.ConfigReference{
			{
				Source: StringToPtr(GetS3Url(factory.userDataBucket, ContainerLinuxBaseIgnitionUri)),
			},
			{
				Source: StringToPtr(GetS3Url(factory.userDataBucket, templateConfigUri)),
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

	file, err := os.Create("index.html")
	if err != nil {
		ignitionLogger.Error(err, "can not create file")
	}

	defer file.Close()
	downloader := s3manager.NewDownloader(factory.session)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(factory.userDataBucket),
		Key:    aws.String("index.html"),
	})
	if err != nil {
		ignitionLogger.Error(err, "failed to get index.html")
	}

	uploader := s3manager.NewUploader(factory.session)
	filePath := strings.Join([]string{factory.userdataDir, uuid.New().String()}, "/")
	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:         bytes.NewReader(userdata),
		Bucket:       aws.String(factory.userDataBucket),
		Expires:      aws.Time(time.Now().Add(time.Hour * 24)),
		Key:          aws.String(filePath),
		StorageClass: aws.String(s3.StorageClassIntelligentTiering),
	})
	if err != nil {
		ignitionLogger.Error(err, "failed to upload ignition file to bucket")
		return nil, err
	}
	out := factory.getIngitionBaseConfig()
	out.Ignition.Config = ignTypes.IgnitionConfig{
		Replace: ignTypes.ConfigReference{
			Source: StringToPtr(GetS3Url(factory.userDataBucket, filePath)),
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
func (factory *S3TemplateBackend) getK8sIgnitionFileName(k8sVersion string) string {
	return "//" + factory.templateDir + "//" + "k8s-" + k8sVersion + ".ign"
}
