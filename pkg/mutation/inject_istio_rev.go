package mutation

import (
	"github.com/sirupsen/logrus"
	"github.com/slackhq/simple-kubernetes-webhook/pkg/features"
	corev1 "k8s.io/api/core/v1"
)

// injectEnv is a container for the mutation injecting environment vars
type injectIstioRev struct {
	Logger logrus.FieldLogger
}

// injectEnv implements the podMutator interface
var _ podMutator = (*injectIstioRev)(nil)

// Name returns the struct name
func (sh injectIstioRev) Name() string {
	return "inject_istio_rev"
}

// Mutate returns a new mutated pod according to set env rules
func (sh injectIstioRev) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {

	const (
		mwcSkipKey = "skip-shard"
	)
	mpod := pod.DeepCopy()

	// pod.Labels[mwcKey] == ""
	for key, _ := range pod.Labels {
		if key == mwcSkipKey {
			sh.Logger.WithField(mwcSkipKey, false).
				Printf("skipping mutation as `%s` label found", key)
			return mpod, nil
		}
	}

	sh.Logger = sh.Logger.WithField("mutation", sh.Name())

	if pod.Labels == nil || pod.Labels[mwcKey] == "" {
		sh.Logger.WithField(mwcKey, false).
			Printf("no %s label found, applying default shard", mwcKey)

		labels := mpod.Labels
		labels[mwcKey] = features.InjectedLabelValue
		mpod.SetLabels(labels)

		return mpod, nil
	}

	return mpod, nil
}
