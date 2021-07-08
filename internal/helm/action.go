package helm

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

type HelmList interface {
	RunCommand() ([]*release.Release, error)
	NewList(*action.Configuration) *action.List
}

type List struct {
	a *action.List
}

func New() *List {
	return &List{}
}

func (l *List) NewList(a *action.Configuration) *action.List {
	l.a = action.NewList(a)
	return l.a
}

func (l *List) RunCommand() ([]*release.Release, error) {
	return l.a.Run()
}
