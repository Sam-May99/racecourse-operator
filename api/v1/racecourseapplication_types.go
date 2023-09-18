/*
Copyright 2023.

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RacecourseApplicationSpec defines the desired state of RacecourseApplication
type RacecourseApplicationSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	//+kubebuilder:validation:Minimum=0

	// The number of instances of the racecourse application to create on the cluster
	Replicas *int64 `json:"replicas,omitempty"`

	// Container information for the racecourses if provided
	Containers []corev1.Container `json:"containers,omitempty"`
}

// RacecourseApplicationStatus defines the observed state of RacecourseApplication
type RacecourseApplicationStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
	Deployed bool `json:"deployed,omitempty"`

	// The names of the deployments spun up by the reconciler
	Racecourses []string `json:"racecourses,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RacecourseApplication is the Schema for the racecourseapplications API
type RacecourseApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RacecourseApplicationSpec   `json:"spec,omitempty"`
	Status RacecourseApplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RacecourseApplicationList contains a list of RacecourseApplication
type RacecourseApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RacecourseApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RacecourseApplication{}, &RacecourseApplicationList{})
}
