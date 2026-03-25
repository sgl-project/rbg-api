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

package constants

// ========== Annotations ==========

// Group level annotations
const (
	// GroupExclusiveTopologyKey declares the topology domain (e.g. kubernetes.io/hostname)
	// for 1:1 exclusive scheduling.
	GroupExclusiveTopologyKey = RBGPrefix + "group-exclusive-topology"

	// DisableExclusiveKeyAnnotationKey can be set to "true" on a Pod template
	// to skip exclusive-topology affinity injection for that pod.
	DisableExclusiveKeyAnnotationKey = RBGPrefix + "role-disable-exclusive"

	// GangSchedulingAnnotationKey enables gang scheduling for a RoleBasedGroup when set to "true".
	GangSchedulingAnnotationKey = RBGPrefix + "group-gang-scheduling"

	// GangSchedulingScheduleTimeoutSecondsKey specifies the schedule timeout seconds for
	// scheduler-plugins based gang scheduling. Defaults to 60 seconds if not set.
	GangSchedulingScheduleTimeoutSecondsKey = RBGPrefix + "group-gang-scheduling-timeout"

	// GangSchedulingVolcanoPriorityClassKey specifies the PriorityClassName for volcano gang scheduling.
	GangSchedulingVolcanoPriorityClassKey = RBGPrefix + "group-gang-scheduling-volcano-priority"

	// GangSchedulingVolcanoQueueKey specifies the Queue for volcano gang scheduling.
	GangSchedulingVolcanoQueueKey = RBGPrefix + "group-gang-scheduling-volcano-queue"
)

// Role level annotations
const (
	// RoleSizeAnnotationKey identifies the role replica size
	RoleSizeAnnotationKey = RBGPrefix + "role-size"

	// RoleDisableExclusiveKey can be set to "true" on a Role template
	// to skip exclusive-topology affinity injection for that role.
	RoleDisableExclusiveKey = RBGPrefix + "role-disable-exclusive"
)

// RoleInstance level annotations
const (
	// RoleInstancePatternKey identifies the RoleInstance organization pattern (Stateful/Stateless)
	RoleInstancePatternKey = RBGPrefix + "role-instance-pattern"

	// RoleInstanceGangSchedulingAnnotationKey enables gang-scheduling aware behavior at the
	// RoleInstance level when set to "true".
	RoleInstanceGangSchedulingAnnotationKey = RBGPrefix + "role-instance-gang-scheduling"

	// DiscoveryConfigModeAnnotationKey identifies discovery config handling mode.
	DiscoveryConfigModeAnnotationKey = RBGPrefix + "discovery-config-mode"
)

// Lifecycle management annotations
const (
	// RoleInstanceSetLifecycleStateKey identifies the lifecycle state of a RoleInstance (stored in labels)
	RoleInstanceSetLifecycleStateKey = RBGPrefix + "role-instance-lifecycle-state"

	// RoleInstanceSetLifecycleTimestampKey identifies the lifecycle state change timestamp of a RoleInstance
	RoleInstanceSetLifecycleTimestampKey = RBGPrefix + "role-instance-lifecycle-timestamp"
)

// InPlace update annotations
const (
	// InPlaceUpdateStateKey identifies the in-place update state
	InPlaceUpdateStateKey = RBGPrefix + "inplace-update-state"

	// InPlaceUpdateGraceKey identifies the in-place update grace period configuration
	InPlaceUpdateGraceKey = RBGPrefix + "inplace-update-grace"

	// RuntimeContainerMetaKey is a key in pod annotations for runtime container states.
	RuntimeContainerMetaKey = "workloads.x-k8s.io/runtime-containers-meta"
)
