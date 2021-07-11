package delete

import (
	"os"

	internalhelm "github.com/edify42/helm-janitor/internal/helm"
	log "github.com/sirupsen/logrus"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

// RunV2 is the main exported method to delete a release
func RunV2(sr InputRun) {
	mycfg := sr.Config()
	cfg := sr.Makeawscfg()
	cluster := sr.Getekscluster(cfg)
	if !mycfg.Env.DebugFlag {
		log.Info("should clean " + cluster.CAFile)
		defer os.Remove(cluster.CAFile)
	} else {
		log.Info("no clean")
	}

	actionConfig := new(action.Configuration)
	rel := release.Release{
		Name:      mycfg.Release,
		Namespace: mycfg.Namespace,
	}
	del := internalhelm.NewDelete()
	err := sr.Deleterelease(cluster, actionConfig, &rel, del)
	if err != nil {
		log.Error(err)
	}
}

// ValidateArgs should check the argument (release)
func ValidateArgs(a []string) {
}
