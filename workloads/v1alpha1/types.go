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

package v1alpha1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// AdapterPhase defines the phase of a RoleBasedGroupScalingAdapter.
type AdapterPhase string

const (
	AdapterPhaseNone     AdapterPhase = ""
	AdapterPhaseNotBound AdapterPhase = "NotBound"
	AdapterPhaseBound    AdapterPhase = "Bound"
)

// RolloutStrategyType defines the strategy type for rollout.
type RolloutStrategyType string

const (
	// RollingUpdateStrategyType indicates replicas will be updated one by one.
	RollingUpdateStrategyType RolloutStrategyType = "RollingUpdate"
)

// UpdateStrategyType defines strategies for Instances in-place update.
type UpdateStrategyType string

const (
	// RecreateUpdateStrategyType indicates always delete and recreate Instances.
	RecreateUpdateStrategyType UpdateStrategyType = "Recreate"

	// InPlaceIfPossibleUpdateStrategyType indicates in-place update when possible.
	InPlaceIfPossibleUpdateStrategyType UpdateStrategyType = "InPlaceIfPossible"
)

// RestartPolicyType defines the restart policy for RBG.
type RestartPolicyType string

const (
	// NoneRestartPolicy follows the same behavior as the StatefulSet/Deployment.
	NoneRestartPolicy RestartPolicyType = "None"

	// RecreateRBGOnPodRestart recreates all pods in the rbg if any individual pod is recreated.
	RecreateRBGOnPodRestart RestartPolicyType = "RecreateRBGOnPodRestart"

	// RecreateRoleInstanceOnPodRestart recreates an instance of a role on pod failure.
	RecreateRoleInstanceOnPodRestart RestartPolicyType = "RecreateRoleInstanceOnPodRestart"
)

// ProgressionType defines how to wait for pods before proceeding to next batch.
type ProgressionType string

const (
	// OrderScheduled means wait for all pods in current batch to be scheduled.
	OrderScheduled ProgressionType = "OrderScheduled"

	// OrderReady means wait for all pods in current batch to be ready.
	OrderReady ProgressionType = "OrderReady"
)

// InPlaceUpdateStrategy defines the strategies for in-place update.
type InPlaceUpdateStrategy struct {
	// GracePeriodSeconds is the timespan between set Pod status to not-ready and update images in Pod spec.
	GracePeriodSeconds int32 `json:"gracePeriodSeconds,omitempty"`
}

// RolloutStrategy defines the strategy that the rbg controller will use to perform replica updates of role.
type RolloutStrategy struct {
	// Type defines the rollout strategy.
	// +kubebuilder:validation:Enum={RollingUpdate}
	// +kubebuilder:default=RollingUpdate
	Type RolloutStrategyType `json:"type"`

	// RollingUpdate defines the parameters to be used when type is RollingUpdateStrategyType.
	// +optional
	RollingUpdate *RollingUpdate `json:"rollingUpdate,omitempty"`
}

// RollingUpdate defines the parameters to be used for RollingUpdateStrategyType.
type RollingUpdate struct {
	// Type indicates the type of the update strategy.
	Type UpdateStrategyType `json:"type,omitempty"`

	// Partition indicates the ordinal at which the role should be partitioned for updates.
	// +optional
	// +kubebuilder:default=0
	Partition *intstr.IntOrString `json:"partition,omitempty"`

	// MaxUnavailable is the maximum number of replicas that can be unavailable during update.
	// +kubebuilder:validation:XIntOrString
	// +kubebuilder:default=1
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`

	// MaxSurge is the maximum number of replicas that can be scheduled above the original number.
	// +kubebuilder:validation:XIntOrString
	// +kubebuilder:default=0
	MaxSurge *intstr.IntOrString `json:"maxSurge,omitempty"`

	// Paused indicates that the update is paused.
	Paused bool `json:"paused,omitempty"`

	// InPlaceUpdateStrategy contains strategies for in-place update.
	InPlaceUpdateStrategy *InPlaceUpdateStrategy `json:"inPlaceUpdateStrategy,omitempty"`
}

// RoleTemplate defines a reusable Pod template that can be referenced by roles.
type RoleTemplate struct {
	// Name is the unique identifier for this template.
	Name string `json:"name"`

	// Template defines the Pod template specification.
	Template corev1.PodTemplateSpec `json:"template"`
}

// TemplateRef references a RoleTemplate defined in spec.roleTemplates.
type TemplateRef struct {
	// Name of the RoleTemplate to reference.
	Name string `json:"name"`
}

// TemplateSource defines either an inline template or a reference to a RoleTemplate.
// +kubebuilder:validation:XValidation:rule="!(has(self.template) && has(self.templateRef))",message="template and templateRef are mutually exclusive"
type TemplateSource struct {
	// Template defines the Pod template specification inline.
	// +optional
	Template *corev1.PodTemplateSpec `json:"template,omitempty"`

	// TemplateRef references a RoleTemplate from spec.roleTemplates.
	// +optional
	TemplateRef *TemplateRef `json:"templateRef,omitempty"`
}

// PodGroupPolicy represents a PodGroup configuration for gang-scheduling.
type PodGroupPolicy struct {
	PodGroupPolicySource `json:",inline"`
}

// PodGroupPolicySource represents supported plugins for gang-scheduling.
type PodGroupPolicySource struct {
	// KubeScheduling plugin for gang-scheduling.
	KubeScheduling *KubeSchedulingPodGroupPolicySource `json:"kubeScheduling,omitempty"`

	// VolcanoScheduling plugin for gang-scheduling.
	VolcanoScheduling *VolcanoSchedulingPodGroupPolicySource `json:"volcanoScheduling,omitempty"`
}

// KubeSchedulingPodGroupPolicySource represents configuration for Kubernetes scheduling plugin.
type KubeSchedulingPodGroupPolicySource struct {
	// Time threshold to schedule PodGroup for gang-scheduling.
	// +kubebuilder:default=60
	ScheduleTimeoutSeconds *int32 `json:"scheduleTimeoutSeconds,omitempty"`
}

// VolcanoSchedulingPodGroupPolicySource represents configuration for volcano podgroup scheduling plugin.
type VolcanoSchedulingPodGroupPolicySource struct {
	// PriorityClassName indicates the PodGroup's priority.
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`

	// Queue defines the queue to allocate resource for PodGroup.
	// +optional
	Queue string `json:"queue,omitempty"`
}

// CoordinationScaling defines the scaling coordination strategy for progressive deployment.
type CoordinationScaling struct {
	// MaxSkew defines the maximum allowed difference in deployment progress between roles.
	// +kubebuilder:validation:Pattern=`^([0-9]|[1-9][0-9]|100)%$`
	// +optional
	MaxSkew *string `json:"maxSkew,omitempty"`

	// Progression defines the progression strategy for scaling.
	// +kubebuilder:validation:Enum=OrderScheduled;OrderReady
	// +kubebuilder:default=OrderScheduled
	// +optional
	Progression *ProgressionType `json:"progression,omitempty"`
}

// CoordinationRollingUpdate describes the rolling update coordination strategy.
type CoordinationRollingUpdate struct {
	// MaxSkew defines the max skew requirement about updated replicas between the roles.
	// +kubebuilder:validation:Pattern=`^([0-9]|[1-9][0-9]|100)%$`
	MaxSkew *string `json:"maxSkew,omitempty"`

	// Partition indicates the replicas at which the role should be partitioned for rolling update.
	// +kubebuilder:validation:Pattern=`^([0-9]|[1-9][0-9]|100)%$`
	Partition *string `json:"partition,omitempty"`

	// MaxUnavailable defines the updating step during rolling.
	// +kubebuilder:validation:Pattern=`^([0-9]|[1-9][0-9]|100)%$`
	MaxUnavailable *string `json:"maxUnavailable,omitempty"`
}

// CoordinationStrategy describes the coordination strategies.
type CoordinationStrategy struct {
	// RollingUpdate defines the coordination strategies about rolling update.
	RollingUpdate *CoordinationRollingUpdate `json:"rollingUpdate,omitempty"`

	Scaling *CoordinationScaling `json:"scaling,omitempty"`
}

// Coordination describes the requirements of coordination strategies for roles.
type Coordination struct {
	// Name of the coordination.
	Name string `json:"name"`

	// Roles that should be constrained by this coordination.
	Roles []string `json:"roles"`

	// Strategy describes the coordination strategies.
	Strategy *CoordinationStrategy `json:"strategy,omitempty"`
}

// WorkloadSpec describes the workload type for a role.
type WorkloadSpec struct {
	// +optional
	// +kubebuilder:default="apps/v1"
	APIVersion string `json:"apiVersion"`

	// +optional
	// +kubebuilder:default="StatefulSet"
	Kind string `json:"kind"`
}

func (w *WorkloadSpec) String() string {
	return fmt.Sprintf("%s/%s", w.APIVersion, w.Kind)
}

// EngineRuntime defines the engine runtime for a role.
type EngineRuntime struct {
	// ProfileName specifies the name of the engine runtime profile.
	ProfileName string `json:"profileName"`

	// InjectContainers specifies the containers to be injected with the engine runtime.
	// +optional
	InjectContainers []string `json:"injectContainers,omitempty"`

	// Containers specifies the engine runtime containers to be overridden.
	Containers []corev1.Container `json:"containers,omitempty"`
}

// LeaderWorkerTemplate defines the leader-worker set template configuration.
type LeaderWorkerTemplate struct {
	// Size is the total number of pods in each group.
	// +optional
	// +kubebuilder:default=1
	Size *int32 `json:"size,omitempty"`

	// PatchLeaderTemplate indicates patching LeaderTemplate.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	PatchLeaderTemplate *runtime.RawExtension `json:"patchLeaderTemplate,omitempty"`

	// PatchWorkerTemplate indicates patching WorkerTemplate.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	PatchWorkerTemplate *runtime.RawExtension `json:"patchWorkerTemplate,omitempty"`
}

// InstanceComponent describes a component within an instance.
type InstanceComponent struct {
	// Name is the type name of the component.
	Name string `json:"name"`

	// Size is the number of replicas for Pods that match the component.
	Size *int32 `json:"size,omitempty"`

	// ServiceName is the name of the service that governs this Instance Component.
	ServiceName string `json:"serviceName,omitempty"`

	// Template is the template for the component pods.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Template corev1.PodTemplateSpec `json:"template"`
}

// ScalingAdapter indicates whether the ScalingAdapter is enabled for the Role.
type ScalingAdapter struct {
	// Enable indicates whether the ScalingAdapter is enabled.
	// +optional
	// +kubebuilder:default=false
	Enable bool `json:"enable,omitempty"`
}

// RoleSpec defines the specification for a role in the group.
type RoleSpec struct {
	// Name is a unique identifier for the role.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Labels to be applied to workloads of this role.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations to be applied to workloads of this role.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Replicas is the desired number of replicas for this role.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=1
	Replicas *int32 `json:"replicas"`

	// RolloutStrategy defines the strategy that will be applied to update replicas.
	// +optional
	RolloutStrategy *RolloutStrategy `json:"rolloutStrategy,omitempty"`

	// RestartPolicy defines the restart policy when pod failures happen.
	// +kubebuilder:validation:Enum={None,RecreateRBGOnPodRestart,RecreateRoleInstanceOnPodRestart}
	// +optional
	RestartPolicy RestartPolicyType `json:"restartPolicy,omitempty"`

	// Dependencies of the role.
	// +optional
	Dependencies []string `json:"dependencies,omitempty"`

	// Workload type specification.
	// +kubebuilder:default={apiVersion:"apps/v1", kind:"StatefulSet"}
	// +optional
	Workload WorkloadSpec `json:"workload,omitempty"`

	// TemplateSource defines the Pod template source, either inline or via reference.
	// +optional
	TemplateSource `json:",inline"`

	// TemplatePatch specifies modifications to apply to the referenced template.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	TemplatePatch runtime.RawExtension `json:"templatePatch,omitempty"`

	// LeaderWorkerSet template configuration.
	// +optional
	LeaderWorkerSet *LeaderWorkerTemplate `json:"leaderWorkerSet,omitempty"`

	// Components describe the components that will be created.
	// +optional
	Components []InstanceComponent `json:"components,omitempty"`

	// ServicePorts to expose for this role.
	// +optional
	ServicePorts []corev1.ServicePort `json:"servicePorts,omitempty"`

	// EngineRuntimes defines the engine runtimes to inject.
	// +optional
	EngineRuntimes []EngineRuntime `json:"engineRuntimes,omitempty"`

	// ScalingAdapter configuration for this role.
	// +optional
	ScalingAdapter *ScalingAdapter `json:"scalingAdapter,omitempty"`

	// MinReadySeconds is the minimum seconds for a newly created pod to be ready.
	// +optional
	// +kubebuilder:default=0
	MinReadySeconds int32 `json:"minReadySeconds,omitempty" protobuf:"varint,9,opt,name=minReadySeconds"`
}

// RoleBasedGroupSpec defines the desired state of RoleBasedGroup.
type RoleBasedGroupSpec struct {
	// Roles defines the list of roles in this group.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	Roles []RoleSpec `json:"roles" patchStrategy:"merge" patchMergeKey:"name"`

	// RoleTemplates defines reusable Pod templates that can be referenced by roles.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	RoleTemplates []RoleTemplate `json:"roleTemplates,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// PodGroupPolicy configures gang-scheduling via supported plugins.
	PodGroupPolicy *PodGroupPolicy `json:"podGroupPolicy,omitempty"`

	// CoordinationRequirements describes the requirements of coordination strategies.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	CoordinationRequirements []Coordination `json:"coordination,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

// RoleStatus shows the current state of a specific role.
type RoleStatus struct {
	// Name of the role.
	Name string `json:"name"`

	// ReadyReplicas is the number of ready replicas.
	ReadyReplicas int32 `json:"readyReplicas"`

	// Replicas is the total number of desired replicas.
	Replicas int32 `json:"replicas"`

	// UpdatedReplicas is the total number of updated replicas.
	UpdatedReplicas int32 `json:"updatedReplicas"`
}

// RoleBasedGroupStatus defines the observed state of RoleBasedGroup.
type RoleBasedGroupStatus struct {
	// The generation observed by the controller.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions track the condition of the RBG.
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// RoleStatuses contains the status of individual roles.
	RoleStatuses []RoleStatus `json:"roleStatuses"`
}

// RoleBasedGroupConditionType defines condition types for RoleBasedGroup.
type RoleBasedGroupConditionType string

const (
	// RoleBasedGroupReady means the rbg is available.
	RoleBasedGroupReady RoleBasedGroupConditionType = "Ready"

	// RoleBasedGroupProgressing means rbg is progressing.
	RoleBasedGroupProgressing RoleBasedGroupConditionType = "Progressing"

	// RoleBasedGroupRollingUpdateInProgress means rbg is performing a rolling update.
	RoleBasedGroupRollingUpdateInProgress RoleBasedGroupConditionType = "RollingUpdateInProgress"

	// RoleBasedGroupRestartInProgress means rbg is restarting.
	RoleBasedGroupRestartInProgress RoleBasedGroupConditionType = "RestartInProgress"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:shortName={rbg}

// RoleBasedGroup is the Schema for the rolebasedgroups API.
type RoleBasedGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RoleBasedGroupSpec   `json:"spec,omitempty"`
	Status RoleBasedGroupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RoleBasedGroupList contains a list of RoleBasedGroup.
type RoleBasedGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RoleBasedGroup `json:"items"`
}



// RoleBasedGroupSetSpec defines the desired state of RoleBasedGroupSet.
type RoleBasedGroupSetSpec struct {
	// Replicas is the number of RoleBasedGroup that will be created.
	// +kubebuilder:default=1
	Replicas *int32 `json:"replicas,omitempty"`

	// Template describes the RoleBasedGroup that will be created.
	Template RoleBasedGroupSpec `json:"template"`
}

// RoleBasedGroupSetConditionType defines condition types for RoleBasedGroupSet.
type RoleBasedGroupSetConditionType string

const (
	RoleBasedGroupSetReady RoleBasedGroupSetConditionType = "Ready"
)

// RoleBasedGroupSetStatus defines the observed state of RoleBasedGroupSet.
type RoleBasedGroupSetStatus struct {
	// The generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// Replicas is the current number of RoleBasedGroups.
	// +optional
	Replicas int32 `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`

	// ReadyReplicas is the number of ready RoleBasedGroups.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas" protobuf:"varint,3,opt,name=readyReplicas"`

	// Conditions track the condition of the RoleBasedGroupSet.
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas
// +kubebuilder:printcolumn:name="DESIRED",type="string",JSONPath=".status.replicas",description="desired replicas"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.readyReplicas",description="ready replicas"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:shortName={rbgs}

// RoleBasedGroupSet is the Schema for the rolebasedgroupsets API.
type RoleBasedGroupSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RoleBasedGroupSetSpec   `json:"spec,omitempty"`
	Status RoleBasedGroupSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RoleBasedGroupSetList contains a list of RoleBasedGroupSet.
type RoleBasedGroupSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RoleBasedGroupSet `json:"items"`
}



// AdapterScaleTargetRef holds the name and role for scaling target reference.
type AdapterScaleTargetRef struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

// RoleBasedGroupScalingAdapterSpec defines the desired state of RoleBasedGroupScalingAdapter.
type RoleBasedGroupScalingAdapterSpec struct {
	// Replicas is the number of RoleBasedGroupRole that will be scaled.
	Replicas *int32 `json:"replicas,omitempty"`

	// ScaleTargetRef is a reference to the target resource that should be scaled.
	ScaleTargetRef *AdapterScaleTargetRef `json:"scaleTargetRef"`
}

// RoleBasedGroupScalingAdapterStatus shows the current state of a RoleBasedGroupScalingAdapter.
type RoleBasedGroupScalingAdapterStatus struct {
	// Phase indicates the current phase of the RoleBasedGroupScalingAdapter.
	Phase AdapterPhase `json:"phase,omitempty"`

	// Replicas is the current effective number of target RoleBasedGroupRole.
	Replicas *int32 `json:"replicas,omitempty"`

	// Selector is a label query used to filter and identify a set of resources for metrics collection.
	Selector string `json:"selector,omitempty"`

	// LastScaleTime is the last time the RoleBasedGroupScalingAdapter scaled the number of pods.
	LastScaleTime *metav1.Time `json:"lastScaleTime,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector
// +kubebuilder:resource:shortName={rbgsa}
// +kubebuilder:printcolumn:name="PHASE",type="string",JSONPath=".status.phase",description="The current phase of the adapter"
// +kubebuilder:printcolumn:name="REPLICAS",type="integer",JSONPath=".status.replicas",description="The current number of replicas"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp",description="Time since creation."

// RoleBasedGroupScalingAdapter is the Schema for the rolebasedgroupscalingadapters API.
type RoleBasedGroupScalingAdapter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RoleBasedGroupScalingAdapterSpec   `json:"spec,omitempty"`
	Status RoleBasedGroupScalingAdapterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RoleBasedGroupScalingAdapterList contains a list of RoleBasedGroupScalingAdapter.
type RoleBasedGroupScalingAdapterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RoleBasedGroupScalingAdapter `json:"items"`
}


