package ignition

import (
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
