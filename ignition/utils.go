package ignition

import (
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	ignitionLogger = ctrl.Log.WithName("ignition")
)

func intToPtr(i int) *int {
	return &i
}

func boolToPtr(b bool) *bool {
	return &b
}

func StringToPtr(s string) *string {
	return &s
}

func GetS3Url(bucket, uri string) string {
	url := &url.URL{
		Scheme: "s3",
		Host:   bucket,
		Path:   uri,
	}
	return url.String()
}
