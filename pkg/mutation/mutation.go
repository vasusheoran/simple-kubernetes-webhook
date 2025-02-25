package mutation

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/wI2L/jsondiff"
	corev1 "k8s.io/api/core/v1"
)

const (
	mwcKey = "istio.io/rev"
)

// Mutator is a container for mutation
type Mutator struct {
	Logger *logrus.Entry
}

// NewMutator returns an initialised instance of Mutator
func NewMutator(logger *logrus.Entry) *Mutator {
	return &Mutator{Logger: logger}
}

// podMutators is an interface used to group functions mutating pods
type podMutator interface {
	Mutate(*corev1.Pod) (*corev1.Pod, error)
	Name() string
}

// MutatePodPatch returns a json patch containing all the mutations needed for
// a given pod
func (m *Mutator) MutatePodPatch(pod *corev1.Pod) ([]byte, error) {
	var podName string
	if pod.Name != "" {
		podName = pod.Name
	} else {
		if pod.ObjectMeta.GenerateName != "" {
			podName = pod.ObjectMeta.GenerateName
		}
	}
	log := logrus.WithField("pod_name", podName)

	mutations := m.getMutations(log, pod)

	mpod := pod.DeepCopy()

	// apply all mutations
	for _, m := range mutations {
		var err error
		mpod, err = m.Mutate(mpod)
		if err != nil {
			return nil, err
		}
	}

	// generate json patch
	patch, err := jsondiff.Compare(pod, mpod)
	if err != nil {
		return nil, err
	}

	patchb, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	return patchb, nil
}

func (m *Mutator) getMutations(logger *logrus.Entry, pod *corev1.Pod) []podMutator {
	logger.Infof("fetching mutations for `%s`", pod.Name)

	for key, value := range pod.Labels {
		if key == mwcKey {
			logger.Infof("found mutation `%s=%s` for pod `%s`. Adding label `simple-webhook-injection/validation=%s`.", mwcKey, value, pod.Name, mwcKey)
			return []podMutator{
				injectValidation{Logger: logger},
			}
		}
	}

	logger.Infof("injecting default shard for `%s`", pod.Name)
	return []podMutator{
		injectIstioRev{Logger: logger},
	}

}
