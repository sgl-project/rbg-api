/*
Copyright 2026 The RBG Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/rbgs/api/workloads/constants"
)

const (
	DefaultRoleInstanceSetMaxUnavailable = "10%"
)

// RoleInstanceSetSpec defines the desired state of RoleInstanceSet
type RoleInstanceSetSpec struct {
	// Replicas is the desired number of replicas of the given Template.
	Replicas *int32 `json:"replicas,omitempty"`

	// Selector is a label query over role instances that should match the replica count.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// RoleInstanceTemplate describes the data a role instance should have when created from a template
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	RoleInstanceTemplate RoleInstanceTemplate `json:"roleInstanceTemplate"`

	// PodManagementPolicy controls how pods are created during the initial scale-up.
	// +optional
	// +kubebuilder:default=Parallel
	PodManagementPolicy constants.PodManagementPolicyType `json:"podManagementPolicy,omitempty"`

	// ScaleStrategy indicates the ScaleStrategy that will be employed to create and delete RoleInstances.
	// +optional
	ScaleStrategy RoleInstanceSetScaleStrategy `json:"scaleStrategy,omitempty"`

	// UpdateStrategy indicates the UpdateStrategy that will be employed to update RoleInstances.
	// +optional
	UpdateStrategy RoleInstanceSetUpdateStrategy `json:"updateStrategy,omitempty"`

	// RevisionHistoryLimit is the maximum number of revisions that will be maintained.
	RevisionHistoryLimit *int32 `json:"revisionHistoryLimit,omitempty"`

	// MinReadySeconds is the minimum number of seconds for which a newly created RoleInstance should be ready.
	// +optional
	MinReadySeconds int32 `json:"minReadySeconds,omitempty"`

	// Lifecycle defines the lifecycle hooks for RoleInstances.
	Lifecycle *RoleInstanceSetLifecycle `json:"lifecycle,omitempty"`
}

// RoleInstanceSetLifecycle contains the hooks for RoleInstance lifecycle.
type RoleInstanceSetLifecycle struct {
	// PreDelete is the hook before RoleInstance to be deleted.
	PreDelete *RoleInstanceSetLifecycleHook `json:"preDelete,omitempty"`
	// InPlaceUpdate is the hook before RoleInstance to update and after RoleInstance has been updated.
	InPlaceUpdate *RoleInstanceSetLifecycleHook `json:"inPlaceUpdate,omitempty"`
}

// RoleInstanceSetLifecycleHook defines handlers for a lifecycle event.
type RoleInstanceSetLifecycleHook struct {
	LabelsHandler     map[string]string `json:"labelsHandler,omitempty"`
	FinalizersHandler []string          `json:"finalizersHandler,omitempty"`
	// MarkNotReady = true means RoleInstance will be set to 'NotReady' at preparingDelete/preparingUpdate state.
	// Default to false.
	MarkNotReady bool `json:"markPodNotReady,omitempty"`
}

// RoleInstanceSetScaleStrategy defines strategies for RoleInstances scale.
type RoleInstanceSetScaleStrategy struct {
	// RoleInstanceToDelete is the names of RoleInstance should be deleted.
	RoleInstanceToDelete []string `json:"roleInstanceToDelete,omitempty"`

	// MaxUnavailable is the maximum number of RoleInstances that can be unavailable for scaled RoleInstances.
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`
}

// RoleInstanceSetUpdateStrategy defines strategies for RoleInstances update.
type RoleInstanceSetUpdateStrategy struct {
	// Type indicates the type of the RoleInstanceSetUpdateStrategy. Default is InPlaceIfPossible.
	Type UpdateStrategyType `json:"type,omitempty"`

	// Partition is the desired number of RoleInstances in old revisions.
	Partition *intstr.IntOrString `json:"partition,omitempty"`

	// MaxUnavailable is the maximum number of RoleInstances that can be unavailable during update or scale.
	// Defaults to 10%.
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`

	// MaxSurge is the maximum number of RoleInstances that can be scheduled above the desired replicas.
	MaxSurge *intstr.IntOrString `json:"maxSurge,omitempty"`

	// Paused indicates that the RoleInstanceSet is paused. Default value is false.
	Paused bool `json:"paused,omitempty"`

	// InPlaceUpdateStrategy contains strategies for in-place update.
	InPlaceUpdateStrategy *RoleInstanceSetInPlaceUpdateStrategy `json:"inPlaceUpdateStrategy,omitempty"`
}

// RoleInstanceSetInPlaceUpdateStrategy defines the strategies for in-place update.
type RoleInstanceSetInPlaceUpdateStrategy struct {
	// GracePeriodSeconds is the timespan between set Pod status to not-ready and update images in Pod spec.
	GracePeriodSeconds int32 `json:"gracePeriodSeconds,omitempty"`
}

// RoleInstanceSetStatus defines the observed state of RoleInstanceSet
type RoleInstanceSetStatus struct {
	// ObservedGeneration is the most recent generation observed for this RoleInstanceSet.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Replicas is the number of RoleInstances created by the RoleInstanceSet controller.
	Replicas int32 `json:"replicas"`

	// ReadyReplicas is the number of RoleInstances that have a Ready Condition.
	ReadyReplicas int32 `json:"readyReplicas"`

	// AvailableReplicas is the number of RoleInstances that have been ready for at least minReadySeconds.
	AvailableReplicas int32 `json:"availableReplicas"`

	// CurrentReplicas is the number of RoleInstances from the current revision.
	CurrentReplicas int32 `json:"currentReplicas"`

	// UpdatedReplicas is the number of RoleInstances from the update revision.
	UpdatedReplicas int32 `json:"updatedReplicas"`

	// UpdatedReadyReplicas is the number of updated RoleInstances that have a Ready Condition.
	UpdatedReadyReplicas int32 `json:"updatedReadyReplicas"`

	// ExpectedUpdatedReplicas is the number of RoleInstances that should be updated.
	ExpectedUpdatedReplicas int32 `json:"expectedUpdatedReplicas,omitempty"`

	// UpdateRevision, if not empty, indicates the latest revision of the RoleInstanceSet.
	UpdateRevision string `json:"updateRevision,omitempty"`

	// CurrentRevision, if not empty, indicates the current revision version.
	CurrentRevision string `json:"currentRevision,omitempty"`

	// CollisionCount is the count of hash collisions for the RoleInstanceSet.
	CollisionCount *int32 `json:"collisionCount,omitempty"`

	// Conditions represents the latest available observations of a RoleInstanceSet's current state.
	Conditions []RoleInstanceSetCondition `json:"conditions,omitempty"`

	// LabelSelector is label selectors for query over RoleInstances used by HPA.
	LabelSelector string `json:"labelSelector,omitempty"`
}

// RoleInstanceSetConditionType is type for RoleInstanceSet conditions.
type RoleInstanceSetConditionType string

const (
	RoleInstanceSetConditionFailedScale  RoleInstanceSetConditionType = "FailedScale"
	RoleInstanceSetConditionFailedUpdate RoleInstanceSetConditionType = "FailedUpdate"
)

// RoleInstanceSetCondition describes the state of a RoleInstanceSet at a certain point.
type RoleInstanceSetCondition struct {
	// Type of RoleInstanceSet condition.
	Type RoleInstanceSetConditionType `json:"type"`

	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`

	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`

	// A human-readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

// +genclient
// +genclient:method=GetScale,verb=get,subresource=scale,result=k8s.io/api/autoscaling/v1.Scale
// +genclient:method=UpdateScale,verb=update,subresource=scale,input=k8s.io/api/autoscaling/v1.Scale,result=k8s.io/api/autoscaling/v1.Scale
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=ris,path=roleinstancesets,scope=Namespaced
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.labelSelector
// +kubebuilder:printcolumn:name="DESIRED",type="integer",JSONPath=".spec.replicas"
// +kubebuilder:printcolumn:name="UPDATED",type="integer",JSONPath=".status.updatedReplicas"
// +kubebuilder:printcolumn:name="READY",type="integer",JSONPath=".status.readyReplicas"
// +kubebuilder:printcolumn:name="TOTAL",type="integer",JSONPath=".status.replicas"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// RoleInstanceSet is the Schema for the RoleInstanceSet API
type RoleInstanceSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RoleInstanceSetSpec   `json:"spec,omitempty"`
	Status RoleInstanceSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RoleInstanceSetList contains a list of RoleInstanceSet
type RoleInstanceSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RoleInstanceSet `json:"items"`
}


