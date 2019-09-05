package rollout

// import (
// 	"testing"

// 	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/util/intstr"
// 	"k8s.io/utils/pointer"
// )

// func experimentFromSteps(r *v1alpha1.Rollout) *v1alpha1.Experiment {
// 	return &v1alpha1.Experiment{}
// }

// func TestRolloutCreateExperiment(t *testing.T) {
// 	f := newFixture(t)
// 	defer f.Close()

// 	steps := []v1alpha1.CanaryStep{{
// 		Experiment: &v1alpha1.RolloutCanaryExperimentStep{
// 			Templates: []v1alpha1.RolloutExperimentTemplate{{
// 				Name:     "stable-template",
// 				SpecRef:  v1alpha1.StableSpecRef,
// 				Replicas: int32(1),
// 			}},
// 		},
// 	}}

// 	r1 := newCanaryRollout("foo", 1, nil, steps, pointer.Int32Ptr(0), intstr.FromInt(0), intstr.FromInt(1))
// 	r2 := bumpVersion(r1)
// 	ex := experimentFromSteps(r2)

// 	rs1 := newReplicaSetWithStatus(r1, 1, 1)
// 	rs2 := newReplicaSetWithStatus(r2, 0, 0)
// 	f.kubeobjects = append(f.kubeobjects, rs1, rs2)
// 	f.replicaSetLister = append(f.replicaSetLister, rs1, rs2)
// 	rs1PodHash := rs1.Labels[v1alpha1.DefaultRolloutUniqueLabelKey]

// 	r2 = updateCanaryRolloutStatus(r2, rs1PodHash, 1, 1, 1, false)

// 	f.rolloutLister = append(f.rolloutLister, r2)
// 	f.objects = append(f.objects, r2)

// 	f.expectCreateExperimentAction(ex)
// 	f.expectPatchRolloutAction(r1)

// 	f.run(getKey(r2, t))
// }

// func TestRolloutExperimentProcessingDoNothing(t *testing.T) {
// 	f := newFixture(t)
// 	defer f.Close()

// 	steps := []v1alpha1.CanaryStep{{
// 		Experiment: &v1alpha1.RolloutCanaryExperimentStep{},
// 	}}

// 	r1 := newCanaryRollout("foo", 1, nil, steps, pointer.Int32Ptr(0), intstr.FromInt(0), intstr.FromInt(1))
// 	r2 := bumpVersion(r1)
// 	ex := experimentFromSteps(r2)

// 	rs1 := newReplicaSetWithStatus(r1, 1, 1)
// 	rs2 := newReplicaSetWithStatus(r2, 0, 0)
// 	f.kubeobjects = append(f.kubeobjects, rs1, rs2)
// 	f.replicaSetLister = append(f.replicaSetLister, rs1, rs2)
// 	rs2PodHash := rs2.Labels[v1alpha1.DefaultRolloutUniqueLabelKey]

// 	r2 = updateCanaryRolloutStatus(r2, rs2PodHash, 1, 1, 1, false)

// 	f.rolloutLister = append(f.rolloutLister, r2)
// 	f.experimentLister = append(f.experimentLister, ex)
// 	f.objects = append(f.objects, r2, ex)

// 	f.expectPatchRolloutAction(r1)
// 	f.run(getKey(r2, t))
// }

// func TestRolloutFailedExperimentEnterDegraded(t *testing.T) {
// 	f := newFixture(t)
// 	defer f.Close()

// 	steps := []v1alpha1.CanaryStep{{
// 		Experiment: &v1alpha1.RolloutCanaryExperimentStep{},
// 	}}

// 	r1 := newCanaryRollout("foo", 1, nil, steps, pointer.Int32Ptr(0), intstr.FromInt(0), intstr.FromInt(1))
// 	r2 := bumpVersion(r1)
// 	ex := experimentFromSteps(r2)

// 	rs1 := newReplicaSetWithStatus(r1, 1, 1)
// 	rs2 := newReplicaSetWithStatus(r2, 0, 0)
// 	f.kubeobjects = append(f.kubeobjects, rs1, rs2)
// 	f.replicaSetLister = append(f.replicaSetLister, rs1, rs2)
// 	rs2PodHash := rs2.Labels[v1alpha1.DefaultRolloutUniqueLabelKey]

// 	r2 = updateCanaryRolloutStatus(r2, rs2PodHash, 1, 1, 1, false)

// 	f.rolloutLister = append(f.rolloutLister, r2)
// 	f.experimentLister = append(f.experimentLister, ex)
// 	f.objects = append(f.objects, r2, ex)

// 	f.expectPatchRolloutAction(r1)
// 	f.run(getKey(r2, t))
// }

// func TestRolloutExperimentScaleDownExtraExperiment(t *testing.T) {
// 	f := newFixture(t)
// 	defer f.Close()

// 	steps := []v1alpha1.CanaryStep{{
// 		Experiment: &v1alpha1.RolloutCanaryExperimentStep{},
// 	}}

// 	ex := &v1alpha1.Experiment{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:            "rando",
// 			OwnerReferences: []metav1.OwnerReference{},
// 		},
// 	}

// 	r1 := newCanaryRollout("foo", 1, nil, steps, pointer.Int32Ptr(0), intstr.FromInt(0), intstr.FromInt(1))
// 	r2 := bumpVersion(r1)
// 	// ex := experimentFromSteps(r2)

// 	rs1 := newReplicaSetWithStatus(r1, 1, 1)
// 	rs2 := newReplicaSetWithStatus(r2, 0, 0)
// 	f.kubeobjects = append(f.kubeobjects, rs1, rs2)
// 	f.replicaSetLister = append(f.replicaSetLister, rs1, rs2)
// 	rs2PodHash := rs2.Labels[v1alpha1.DefaultRolloutUniqueLabelKey]

// 	r2 = updateCanaryRolloutStatus(r2, rs2PodHash, 1, 1, 1, false)

// 	f.rolloutLister = append(f.rolloutLister, r2)
// 	f.experimentLister = append(f.experimentLister, ex)
// 	f.objects = append(f.objects, r2, ex)

// 	f.expectPatchRolloutAction(r1)
// 	f.run(getKey(r2, t))
// }

// func TestRolloutExperimentFinishedIncrementStep(t *testing.T) {
// 	f := newFixture(t)
// 	defer f.Close()

// 	steps := []v1alpha1.CanaryStep{{
// 		Experiment: &v1alpha1.RolloutCanaryExperimentStep{
// 			Templates: []v1alpha1.RolloutExperimentTemplate{{
// 				Name:     "stable-template",
// 				SpecRef:  v1alpha1.StableSpecRef,
// 				Replicas: int32(1),
// 			}},
// 		},
// 	}}

// 	r1 := newCanaryRollout("foo", 1, nil, steps, pointer.Int32Ptr(0), intstr.FromInt(0), intstr.FromInt(1))
// 	r2 := bumpVersion(r1)
// 	ex := experimentFromSteps(r2)

// 	rs1 := newReplicaSetWithStatus(r1, 1, 1)
// 	rs2 := newReplicaSetWithStatus(r2, 0, 0)
// 	f.kubeobjects = append(f.kubeobjects, rs1, rs2)
// 	f.replicaSetLister = append(f.replicaSetLister, rs1, rs2)
// 	rs2PodHash := rs2.Labels[v1alpha1.DefaultRolloutUniqueLabelKey]

// 	r2 = updateCanaryRolloutStatus(r2, rs2PodHash, 1, 1, 1, false)

// 	f.rolloutLister = append(f.rolloutLister, r2)
// 	f.experimentLister = append(f.experimentLister, ex)
// 	f.objects = append(f.objects, r2, ex)

// 	f.expectPatchRolloutAction(r1)
// 	f.run(getKey(r2, t))
// }
