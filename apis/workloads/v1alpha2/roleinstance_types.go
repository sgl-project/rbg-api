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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RoleInstanceSpec defines the desired state of RoleInstance
type RoleInstanceSpec struct {
	// Components is a list of components, each of which specifies a component and the number of replicas and template for RoleInstance that match the component.
	Components []RoleInstanceComponent `json:"components" patchStrategy:"merge" patchMergeKey:"name"`

	// RoleInstanceReadyPolicy specifies the policy for determining if the RoleInstance is ready.
	// Defaults to `AllPodReady`
	// +kubebuilder:default=AllPodReady
	ReadyPolicy RoleInstanceReadyPolicyType `json:"readyPolicy,omitempty"`

	// RestartPolicy defines the restart policy for all pods within the RoleInstance.
	RestartPolicy RoleInstanceRestartPolicyType `json:"restartPolicy,omitempty"`

	// ReadinessGates is an optional list of PodReadinessGates for the whole RoleInstance.
	ReadinessGates []RoleInstanceReadinessGate `json:"readinessGates,omitempty"`
}

// RoleInstanceReadinessGate contains the reference to a RoleInstance condition
type RoleInstanceReadinessGate struct {
	// ConditionType refers to a condition in the pod's condition list with matching type.
	ConditionType RoleInstanceConditionType `json:"conditionType"`
}

// RoleInstanceReadyPolicyType defines policy for RoleInstance readiness.
type RoleInstanceReadyPolicyType string

const (
	// RoleInstanceReadyOnAllPodReady means all Pods in the RoleInstance must be ready when RoleInstance Ready
	RoleInstanceReadyOnAllPodReady RoleInstanceReadyPolicyType = "AllPodReady"

	// RoleInstanceReadyPolicyTypeNone means do nothing for Pods
	RoleInstanceReadyPolicyTypeNone RoleInstanceReadyPolicyType = "None"
)

// RoleInstanceRestartPolicyType defines the restart policy for RoleInstance.
type RoleInstanceRestartPolicyType string

const (
	// NoneRoleInstanceRestartPolicy will follow the same behavior as the Pod.
	NoneRoleInstanceRestartPolicy RoleInstanceRestartPolicyType = "None"

	// RoleInstanceRestartPolicyRecreateOnPodRestart will recreate a role instance if its Pod restarted.
	RoleInstanceRestartPolicyRecreateOnPodRestart RoleInstanceRestartPolicyType = "RecreateRoleInstanceOnPodRestart"
)

// RoleInstanceComponent defines a component of a RoleInstance.
type RoleInstanceComponent struct {
	// Name is the type name of the component.
	Name string `json:"name"`

	// Size is the number of replicas for Pods that match the PodRule.
	Size *int32 `json:"size,omitempty"`

	// ServiceName is the name of the service that governs this RoleInstance Component.
	ServiceName string `json:"serviceName,omitempty"`

	// Template is the template for the component pods.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Template corev1.PodTemplateSpec `json:"template"`
}

// RoleInstanceStatus defines the observed state of RoleInstance
type RoleInstanceStatus struct {
	// ObservedGeneration is the most recent generation observed for this RoleInstance.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions track the condition of the RoleInstance
	Conditions []RoleInstanceCondition `json:"conditions,omitempty"`

	// ComponentStatuses is a list of RoleInstanceComponentStatus
	ComponentStatuses []RoleInstanceComponentStatus `json:"componentStatuses,omitempty"`

	// LabelSelector of a RoleInstance is a label query over Pods that should match the RoleInstance.
	LabelSelector string `json:"labelSelector,omitempty"`

	// CurrentRevision is a hash value that changes when the spec is changed.
	CurrentRevision string `json:"currentRevision,omitempty"`

	// UpdateRevision is a hash value that changes when the spec is changed.
	UpdateRevision string `json:"updateRevision,omitempty"`

	// CollisionCount is the count of hash collisions for the RoleInstanceSet.
	CollisionCount *int32 `json:"collisionCount,omitempty"`
}

// RoleInstanceComponentStatus defines the status of a component within a RoleInstance.
type RoleInstanceComponentStatus struct {
	// Name is the type name of the component.
	Name string `json:"name"`

	// Size is the number of Pod for RoleInstance that match the component.
	Size int32 `json:"size"`

	// ReadyReplicas is the number of ready Pod for RoleInstance that match the component.
	ReadyReplicas int32 `json:"readyReplicas"`

	// UpdatedReplicas is the number of updated Pod for RoleInstance that match the component.
	UpdatedReplicas int32 `json:"updatedReplicas"`

	// ScheduledReplicas is the number of scheduled Pod for RoleInstance that match the component.
	ScheduledReplicas int32 `json:"scheduledReplicas"`

	// AvailableReplicas is the number of available Pod for RoleInstance that match the component.
	AvailableReplicas int32 `json:"availableReplicas"`

	// UpdatedReadyReplicas is the number of updated and ready Pod for RoleInstance that match the component.
	UpdatedReadyReplicas int32 `json:"updatedReadyReplicas"`
}

// RoleInstanceConditionType is type for RoleInstance conditions.
type RoleInstanceConditionType string

const (
	RoleInstanceReady              RoleInstanceConditionType = "RoleInstanceReady"
	RoleInstanceInPlaceUpdateReady RoleInstanceConditionType = "RoleInstanceInPlaceUpdateReady"
	RoleInstanceCustomReady        RoleInstanceConditionType = "RoleInstanceCustomReady"
	RoleInstanceAllPodsReady       RoleInstanceConditionType = "RoleInstanceAllPodsReady"
	RoleInstanceFailedScale        RoleInstanceConditionType = "FailedScale"
	RoleInstanceFailedUpdate       RoleInstanceConditionType = "FailedUpdate"
)

// RoleInstanceCondition describes the state of a RoleInstance at a certain point.
type RoleInstanceCondition struct {
	// Type of RoleInstance condition.
	Type RoleInstanceConditionType `json:"type"`

	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`

	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`

	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=rins,path=roleinstances,scope=Namespaced
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='RoleInstanceReady')].status",description="Overall readiness status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// RoleInstance is the Schema for the roleinstances API
type RoleInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RoleInstanceSpec   `json:"spec,omitempty"`
	Status RoleInstanceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// RoleInstanceList contains a list of RoleInstance
type RoleInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RoleInstance `json:"items"`
}

// RoleInstanceTemplate is an embedded template for RoleInstanceSet.
type RoleInstanceTemplate struct {
	RoleInstanceSpec `json:",inline"`
}

func init() {
	SchemeBuilder.Register(&RoleInstance{}, &RoleInstanceList{})
}
