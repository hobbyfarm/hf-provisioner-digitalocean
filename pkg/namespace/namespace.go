package namespace

import (
	"io"
	"os"
)

func Resolve() string {
	// first try based on env
	if ns := os.Getenv("HF_NAMESPACE"); ns != "" {
		return ns
	}

	f, err := os.Open("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err == nil {
		data, err := io.ReadAll(f)
		if err == nil {
			return string(data)
		}
	}

	return "default"
}
