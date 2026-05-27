package probes

import (
	"fmt"
	"os"
)

var VERSION = "SNAPSHOT"

type Fake struct {
	Name string
}

func NewFake(name string, args []string) (*Fake, error) {

	return &Fake{
		Name: name,
	}, nil
}

func (p Fake) Exec() Result {

	src := os.Getenv("KUBERNETES_NODE_NAME")
	if len(src) == 0 {
		src = "localhost"
	}

	return Result{
		ProbeID:    p.Name,
		MetricName: fmt.Sprintf("meta"),
		MetricLabels: map[string]string{
			"node":    src,
			"version": VERSION,
		},
		ReturnCode: Success,
		Msg:        "Self-test",
	}
}
