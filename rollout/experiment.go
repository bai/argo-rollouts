package rollout

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	patchtypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/utils/pointer"

	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	"github.com/argoproj/argo-rollouts/utils/defaults"
	experimentutil "github.com/argoproj/argo-rollouts/utils/experiment"
	replicasetutil "github.com/argoproj/argo-rollouts/utils/replicaset"
)

const (
	cancelExperimentPatch = `{
		"status": {
			"running": false
		}
	}`
)

// getExperimentFromTemplate takes the canary experiment step and converts it to an experiment
func getExperimentFromTemplate(r *v1alpha1.Rollout, stableRS, newRS *appsv1.ReplicaSet) (*v1alpha1.Experiment, error) {
	step := replicasetutil.GetCurrentExperimentStep(r)
	if step == nil {
		return nil, nil
	}
	rolloutTemplateHash := controller.ComputeHash(&r.Spec.Template, r.Status.CollisionCount)
	name := fmt.Sprintf("%s-%s", r.Name, rolloutTemplateHash)

	experiment := &v1alpha1.Experiment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       r.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(r, controllerKind)},
		},
		Spec: v1alpha1.ExperimentSpec{
			Duration:                &step.Duration,
			ProgressDeadlineSeconds: pointer.Int32Ptr(defaults.GetProgressDeadlineSecondsOrDefault(r)),
		},
	}
	for i := range step.Templates {
		templateStep := step.Templates[i]
		template := v1alpha1.TemplateSpec{
			Name:     templateStep.Name,
			Replicas: &templateStep.Replicas,
		}
		templateRS := &appsv1.ReplicaSet{}
		switch templateStep.SpecRef {
		case v1alpha1.CanarySpecRef:
			templateRS = newRS
		case v1alpha1.StableSpecRef:
			templateRS = stableRS
		default:
			return nil, fmt.Errorf("Invalid template step SpecRef: must be canary or stable")
		}
		template.Template = templateRS.Spec.Template
		template.MinReadySeconds = templateRS.Spec.MinReadySeconds

		template.Selector = templateRS.Spec.Selector.DeepCopy()
		for key := range templateStep.Metadata.Labels {
			template.Template.ObjectMeta.Labels[key] = templateStep.Metadata.Labels[key]
		}
		for key := range templateStep.Metadata.Annotations {
			template.Template.ObjectMeta.Annotations[key] = templateStep.Metadata.Annotations[key]
		}
		experiment.Spec.Templates = append(experiment.Spec.Templates, template)
	}

	return experiment, nil
}

// getExperimentsForRollout get all experiments owned by the Rollout
// changing steps in the Rollout Spec would cause multiple experiments to exist which is why it returns an array
func (c *RolloutController) getExperimentsForRollout(rollout *v1alpha1.Rollout) ([]*v1alpha1.Experiment, error) {
	experiments, err := c.experimentsLister.Experiments(rollout.Namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}
	//TODO(dthomson) consider adoption
	ownedByRollout := make([]*v1alpha1.Experiment, 0)
	for i := range experiments {
		e := experiments[i]
		controllerRef := metav1.GetControllerOf(e)
		if controllerRef != nil && controllerRef.UID == rollout.UID {
			ownedByRollout = append(ownedByRollout, e)
		}
	}
	return ownedByRollout, nil
}

func (c *RolloutController) reconcileExperiments(rollout *v1alpha1.Rollout, stableRS, newRS *appsv1.ReplicaSet, currentEx *v1alpha1.Experiment, otherExs []*v1alpha1.Experiment) (bool, error) {
	// Check if there are unexpected experiments (don't match
	// Scale them down (Delete them?)
	// Check for expected Experiment
	// check if it should be running (at step in Rollout canary steps)
	//Scale down otherwise (delete?)
	// Check status of experiment
	// Progressing cool!
	// Running cool
	// Degraded add condition to rollout
	// Finished increment step
	// Update rollout status

	for i := range otherExs {
		otherEx := otherExs[i]
		if otherEx.Status.Running != nil && *otherEx.Status.Running {
			_, err := c.argoprojclientset.ArgoprojV1alpha1().Experiments(otherEx.Namespace).Patch(otherEx.Name, patchtypes.MergePatchType, []byte(cancelExperimentPatch))
			if err != nil {
				return false, err
			}
		}
	}

	step, _ := replicasetutil.GetCurrentCanaryStep(rollout)
	if step == nil || step.Experiment == nil {
		if currentEx != nil && currentEx.Status.Running != nil && *currentEx.Status.Running {
			_, err := c.argoprojclientset.ArgoprojV1alpha1().Experiments(currentEx.Namespace).Patch(currentEx.Name, patchtypes.MergePatchType, []byte(cancelExperimentPatch))
			if err != nil {
				return false, err
			}
		}
		return false, nil
	}
	if currentEx == nil {
		newEx, err := getExperimentFromTemplate(rollout, stableRS, newRS)
		if err != nil {
			return false, err
		}
		currentEx, err = c.argoprojclientset.ArgoprojV1alpha1().Experiments(newEx.Namespace).Create(newEx)
		if err != nil {
			return false, err
		}
	}
	if experimentutil.HasFinished(currentEx) {
		return true, nil
	}
	return false, nil
}
