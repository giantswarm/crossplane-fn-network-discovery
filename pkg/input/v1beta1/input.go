// Package v1beta1 contains the input type for this Function
// +kubebuilder:object:generate=true
// +groupName=networkdiscovery.fn.giantswarm.io
// +versionName=v1beta1
package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This isn't a custom resource, in the sense that we never install its CRD.
// It is a KRM-like object, so we generate a CRD to describe its schema.

// Input can be used to provide input to this Function.
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:categories=crossplane
type Input struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Defines the spec for this input
	Spec *Spec `json:"spec,omitempty"`
}

type RemoteVpc struct {
	// The VPC name
	Name string `json:"name"`

	// GroupBy is an AWS tag name that is used to group subnet and route table
	// results into logical "sets" of data
	GroupBy string `json:"groupBy"`

	// The VPC region
	Region string `json:"region"`

	// The VPC provider config
	ProviderConfig string `json:"providerConfig"`
}

// Spec - Defines the spec given to this input type, providing the required,
// and optional elements that may be defined
type Spec struct {
	// EnabledRef A path to a field on the claim that determines if this function
	// is enabled in the current composition allowing for conditional execution
	// of the function in complex compositions
	//
	// +optional
	EnabledRef string `json:"enabledRef,omitempty"`

	// GroupByRef A path to the field on the claim that determines the grouping
	// of the subnets and route tables in the VPC
	//
	// +optional
	GroupByRef string `json:"groupByRef,omitempty"`

	// PatchTo specified the path to apply the VPC map
	//
	// +required
	PatchTo string `json:"patchTo"`

	// ProviderConfig A path to the provider config in the Claim
	//
	// +required
	ProviderConfigRef string `json:"providerConfigRef"`

	// ProviderType dictates what cloud provider the discovery is working against
	//
	// +kubebuilder:validation:Enum=aws;azure;gcp
	// +kubebuilder:default=aws
	// +optional
	ProviderType string `json:"providerType"`

	// Region A path to the region in the Claim
	//
	// +required
	RegionRef string `json:"regionRef"`

	// VpcName A path to the VPC name in the Claim
	//
	// +required
	VpcNameRef string `json:"vpcRef"`
}
