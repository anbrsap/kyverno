/*
Copyright 2020 The Kubernetes authors.

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
	"reflect"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov2beta1 "github.com/kyverno/kyverno/api/kyverno/v2beta1"
	"github.com/robfig/cron"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:shortName=cleanpol,categories=kyverno
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Schedule",type=string,JSONPath=".spec.schedule"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// CleanupPolicy defines a rule for resource cleanup.
type CleanupPolicy struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec declares policy behaviors.
	Spec CleanupPolicySpec `json:"spec"`

	// Status contains policy runtime data.
	// +optional
	Status CleanupPolicyStatus `json:"status,omitempty"`
}

// GetSpec returns the policy spec
func (p *CleanupPolicy) GetSpec() *CleanupPolicySpec {
	return &p.Spec
}

// GetStatus returns the policy status
func (p *CleanupPolicy) GetStatus() *CleanupPolicyStatus {
	return &p.Status
}

// Validate implements programmatic validation
func (p *CleanupPolicy) Validate(clusterResources sets.String) (errs field.ErrorList) {
	errs = append(errs, kyvernov1.ValidatePolicyName(field.NewPath("metadata").Child("name"), p.Name)...)
	errs = append(errs, p.Spec.Validate(field.NewPath("spec"), clusterResources, true)...)
	return errs
}

// GetKind returns the resource kind
func (p *CleanupPolicy) GetKind() string {
	return p.Kind
}

// GetAPIVersion returns the resource kind
func (p *CleanupPolicy) GetAPIVersion() string {
	return p.APIVersion
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CleanupPolicyList is a list of ClusterPolicy instances.
type CleanupPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []CleanupPolicy `json:"items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:scope=Cluster,shortName=ccleanpol,categories=kyverno
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Schedule",type=string,JSONPath=".spec.schedule"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ClusterCleanupPolicy defines rule for resource cleanup.
type ClusterCleanupPolicy struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec declares policy behaviors.
	Spec CleanupPolicySpec `json:"spec"`

	// Status contains policy runtime data.
	// +optional
	Status CleanupPolicyStatus `json:"status,omitempty"`
}

// GetSpec returns the policy spec
func (p *ClusterCleanupPolicy) GetSpec() *CleanupPolicySpec {
	return &p.Spec
}

// GetStatus returns the policy status
func (p *ClusterCleanupPolicy) GetStatus() *CleanupPolicyStatus {
	return &p.Status
}

// GetKind returns the resource kind
func (p *ClusterCleanupPolicy) GetKind() string {
	return p.Kind
}

// GetAPIVersion returns the resource kind
func (p *ClusterCleanupPolicy) GetAPIVersion() string {
	return p.APIVersion
}

// Validate implements programmatic validation
func (p *ClusterCleanupPolicy) Validate(clusterResources sets.String) (errs field.ErrorList) {
	errs = append(errs, kyvernov1.ValidatePolicyName(field.NewPath("metadata").Child("name"), p.Name)...)
	errs = append(errs, p.Spec.Validate(field.NewPath("spec"), clusterResources, false)...)
	return errs
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterCleanupPolicyList is a list of ClusterCleanupPolicy instances.
type ClusterCleanupPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ClusterCleanupPolicy `json:"items"`
}

// CleanupPolicySpec stores specifications for selecting resources that the user needs to delete
// and schedule when the matching resources needs deleted.
type CleanupPolicySpec struct {
	// MatchResources defines when cleanuppolicy should be applied. The match
	// criteria can include resource information (e.g. kind, name, namespace, labels)
	// and admission review request information like the user name or role.
	// At least one kind is required.
	MatchResources kyvernov2beta1.MatchResources `json:"match,omitempty"`

	// ExcludeResources defines when cleanuppolicy should not be applied. The exclude
	// criteria can include resource information (e.g. kind, name, namespace, labels)
	// and admission review request information like the name or role.
	// +optional
	ExcludeResources *kyvernov2beta1.MatchResources `json:"exclude,omitempty"`

	// The schedule in Cron format
	Schedule string `json:"schedule"`

	// Conditions defines conditions used to select resources which user needs to delete
	// +optional
	Conditions *kyvernov2beta1.AnyAllConditions `json:"conditions,omitempty"`
}

// CleanupPolicyStatus stores the status of the policy.
type CleanupPolicyStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

// Validate implements programmatic validation
func (p *CleanupPolicySpec) Validate(path *field.Path, clusterResources sets.String, namespaced bool) (errs field.ErrorList) {
	errs = append(errs, ValidateSchedule(path.Child("schedule"), p.Schedule)...)
	errs = append(errs, p.MatchResources.Validate(path.Child("match"), namespaced, clusterResources)...)
	if p.ExcludeResources != nil {
		errs = append(errs, p.ExcludeResources.Validate(path.Child("exclude"), namespaced, clusterResources)...)
	}
	errs = append(errs, p.ValidateMatchExcludeConflict(path)...)
	return errs
}

// ValidateSchedule validates whether the schedule specified is in proper cron format or not.
func ValidateSchedule(path *field.Path, schedule string) (errs field.ErrorList) {
	if _, err := cron.ParseStandard(schedule); err != nil {
		errs = append(errs, field.Invalid(path, schedule, "schedule spec in the cleanupPolicy is not in proper cron format"))
	}
	return errs
}

// ValidateMatchExcludeConflict checks if the resultant of match and exclude block is not an empty set
func (spec *CleanupPolicySpec) ValidateMatchExcludeConflict(path *field.Path) (errs field.ErrorList) {
	if spec.ExcludeResources == nil || len(spec.ExcludeResources.All) > 0 || len(spec.MatchResources.All) > 0 {
		return errs
	}
	// if both have any then no resource should be common
	if len(spec.MatchResources.Any) > 0 && len(spec.ExcludeResources.Any) > 0 {
		for _, rmr := range spec.MatchResources.Any {
			for _, rer := range spec.ExcludeResources.Any {
				if reflect.DeepEqual(rmr, rer) {
					return append(errs, field.Invalid(path, spec, "CleanupPolicy is matching an empty set"))
				}
			}
		}
		return errs
	}
	if reflect.DeepEqual(spec.ExcludeResources, kyvernov2beta1.MatchResources{}) {
		return errs
	}
	return append(errs, field.Invalid(path, spec, "CleanupPolicy is matching an empty set"))
}
