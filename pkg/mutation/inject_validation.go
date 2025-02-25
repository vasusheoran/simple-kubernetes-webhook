package mutation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// injectEnv is a container for the mutation injecting environment vars
type injectValidation struct {
	Logger logrus.FieldLogger
}

// injectEnv implements the podMutator interface
var _ podMutator = (*injectValidation)(nil)

// Name returns the struct name
func (v injectValidation) Name() string {
	return "inject_validation"
}

// Mutate returns a new mutated pod according to set env rules
func (v injectValidation) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
	const (
		labelKey   = "simple-webhook-injection/validation"
		labelValue = "done"
	)

	v.Logger = v.Logger.WithField("mutation", v.Name())
	mpod := pod.DeepCopy()

	if pod.Labels == nil || pod.Labels[labelKey] == "" {
		v.Logger.WithField(labelKey, false).
			Printf("no %s label found, applying default shard", labelKey)

		labels := mpod.Labels
		labels[mwcKey] = labelValue
		mpod.SetLabels(labels)

		return mpod, nil
	}

	return mpod, nil
}
