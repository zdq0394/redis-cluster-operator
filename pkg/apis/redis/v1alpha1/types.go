/*
Copyright 2017 The Kubernetes Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedisCluster is a specification for a RedisCluster resource
type RedisCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RedisClusterSpec   `json:"spec"`
	Status RedisClusterStatus `json:"status"`
}

// RedisClusterSpec is the spec for a RedisCluster resource
type RedisClusterSpec struct {
	Redis RedisSetting `json:"redis"`
}

// CPUAndMem defines how many cpu and ram the container will request/limit
type CPUAndMem struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// RedisResources sets the limits and requests for a container
type RedisResources struct {
	Requests CPUAndMem `json:"requests,omitempty"`
	Limits   CPUAndMem `json:"limits,omitempty"`
}

type RedisSetting struct {
	Replicas  int32          `json:"replicas,omitempty"`
	Resources RedisResources `json:"resources,omitempty"`
	Image     string         `json:"image,omitempty"`
}

// RedisClusterStatus is the status for a RedisCluster resource
type RedisClusterStatus struct {
	AvailableReplicas int32           `json:"availableReplicas"`
	Instances         []RedisInstance `json:"instances"`
}

// RedisInstance is the instance running.
type RedisInstance struct {
	Role          string
	IP            string
	Port          string
	ClusterNodeID string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedisClusterList is a list of RedisCluster resources
type RedisClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []RedisCluster `json:"items"`
}
