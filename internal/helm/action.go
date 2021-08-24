package helm

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

// HelmList is a wrapper around the helm list API
type HelmList interface {
	RunCommand() ([]*release.Release, error)
	ActionNewList(*action.Configuration) *action.List
}

// HelmDelete is a wrapper around the helm delete api.
type HelmDelete interface {
	RunCommand(string) (*release.UninstallReleaseResponse, error)
	ActionNewUninstall(*action.Configuration) *action.Uninstall
}

type List struct {
	a *action.List
}

// Delete struct to uninstall a given release
type Delete struct {
	d *action.Uninstall
}

func New() *List {
	return &List{}
}

func NewDelete() *Delete {
	return &Delete{}
}

func (d *Delete) ActionNewUninstall(a *action.Configuration) *action.Uninstall {
	d.d = action.NewUninstall(a)
	return d.d
}

func (l *List) ActionNewList(a *action.Configuration) *action.List {
	l.a = action.NewList(a)
	return l.a
}

func (l *List) RunCommand() ([]*release.Release, error) {
	return l.a.Run()
}

func (d *Delete) RunCommand(rel string) (*release.UninstallReleaseResponse, error) {
	return d.d.Run(rel)
}
