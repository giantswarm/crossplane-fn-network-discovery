// Package v1beta1 contains the definition of the XR requirements for using this function
//
// +kubebuilder:object:generate=true
// +groupName=networkdiscovery.fn.giantswarm.io
// +versionName=v1beta1
package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWS is an object that holds VPCs
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:categories=crossplane;composition

type Aws struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// The VPCs defined in this AWS account
	//
	// +mapType=granular
	Vpcs map[string]AwsVpc `json:"vpcs"`
}

// StatusSubnets is a map of subnets and their status
//
// +mapType=atomic
type StatusSubnets map[string]StatusSubnetDetails

type StatusSubnetDetails struct {
	// The ARN of the subnet
	//
	// +optional
	ARN string `json:"arn"`

	// The ID of the subnet
	//
	// +required
	ID string `json:"id"`
}

// StatusRouteTables is a map of route tables and their status
//
// +mapType=atomic
type StatusRouteTables map[string]StatusRouteTableDetails

type StatusRouteTableDetails struct {
	// The ID of the route table
	//
	// +required
	ID string `json:"id"`
}

// Vpc holds VPC information
//
// +structType=granular
type AwsVpc struct {
	// A list of additional VPC CIDR blocks defined in this VPC
	// +listType=atomic
	// +optional
	AdditionalCidrBlocks []string `json:"additionalCidrBlocks,omitempty"`

	// The Ipv4 cidr block defined for this VPC
	// +optional
	CidrBlock string `json:"cidrBlock,omitempty"`

	// ID The VPC ID
	// +kubebuilder:validation:Required
	// +required
	ID string `json:"id,omitempty"`

	// The internet gateway defined in this VPC
	// +optional
	InternetGateway string `json:"internetGateway,omitempty"`

	// A map of NAT gateways defined in this VPC
	// +mapType=atomic
	// +optional
	NatGateways map[string]string `json:"natGateways,omitempty"`

	// The owner of the current VPC
	// +optional
	Owner string `json:"owner,omitempty"`

	// The provider config used to look up this VPC
	// +optional
	ProviderConfig string `json:"providerConfig,omitempty"`

	// A map of private subnets defined in this VPC
	// +listType=atomic
	// +optional
	PrivateSubnets []StatusSubnets `json:"privateSubnets,omitempty"`

	// A list of maps of public subnets defined in this VPC
	// +listType=atomic
	// +optional
	PublicSubnets []StatusSubnets `json:"publicSubnets,omitempty"`

	// A map of private route tables defined in this VPC
	// +listType=atomic
	// +optional
	PrivateRouteTables []StatusRouteTables `json:"privateRouteTables,omitempty"`

	// A map of public route tables defined in this VPC
	// +listType=atomic
	// +optional
	PublicRouteTables []StatusRouteTables `json:"publicRouteTables,omitempty"`

	// The region this VPC is located in
	// +optional
	Region string `json:"region,omitempty"`

	// A map of security groups defined in this VPC
	// +mapType=atomic
	// +optional
	SecurityGroups map[string]string `json:"securityGroups,omitempty"`

	// A map of transit gateways defined in this VPC
	// +mapType=atomic
	// +optional
	TransitGateways map[string]TransitGateway `json:"transitGateways,omitempty"`

	// A map of VPC peering connections defined in this VPC
	// +mapType=atomic
	// +optional
	VpcPeeringConnections map[string]PeeringConnection `json:"vpcPeeringConnections,omitempty"`
}

type PeeringConnection struct {
	// The ID of the VPC peering connection
	//
	// +optional
	ID string `json:"id"`

	// The ARN of the VPC peering connection
	//
	// +optional
	ARN string `json:"arn"`
}

type TransitGateway struct {
	// The ARN of the transit gateway
	//
	// +optional
	ARN string `json:"arn"`

	// The ID of the transit gateway
	//
	// +optional
	ID string `json:"id"`

	// TransitGatewayAttachments The IDs of the transit gateway attachment(s)
	// associated with this transit gateway
	//
	// +optional
	Attachments map[string]TransitGatewayAttachment `json:"attachments"`

	// TransitGatewayRouteTables The IDs of the transit gateway route table(s)
	// associated with this transit gateway
	//
	// +optional
	RouteTables map[string]TransitGatewayRouteTable `json:"routeTables"`
}

type TransitGatewayAttachment struct {
	// The ID of the transit gateway attachment
	//
	// +optional
	ID string `json:"id"`

	// The ID of the resource that the transit gateway is attached to
	//
	// +optional
	ResourceID string `json:"resourceId"`

	// The associated route table ID
	//
	// +optional
	RouteTableID string `json:"routeTableId"`

	// The type of the transit gateway attachment
	//
	// +optional
	Type string `json:"type"`
}

type TransitGatewayRouteTable struct {
	// The ID of the transit gateway route table
	//
	// +optional
	ID string `json:"id"`

	// Is this the default route table for the transit gateway
	//
	// +optional
	DefaultAssociation bool `json:"defaultAssociation"`

	// Is this the default propagation route table for the transit gateway
	//
	// +optional
	DefaultPropagation bool `json:"defaultPropagation"`
}

// AwsSubnet is an object that holds information about a subnet defined in AWS
// +mapType=granular
type AwsSubnet struct {
	// The ARN of the subnet
	// +optional
	ARN string `json:"arn"`

	// ID The subnet ID
	// +kubebuilder:validation:Required
	ID string `json:"id"`

	// AvailabilityZone The availability zone this subnet is located in
	// +optional
	AvailabilityZone string `json:"availabilityZone"`

	// The Ipv4 cidr block defined for this subnet
	// +optional
	CidrBlock string `json:"cidrBlock"`

	// Is this subnet enabled for IPv6
	// +optional
	IsIpv6 bool `json:"isIpV6"`

	// The IPv6 CIDR block (if defined) for this subnet
	// +optional
	Ipv6CidrBlock string `json:"ipv6CidrBlock"`

	// Is this a public subnet. Determined by validating an internet gateway on
	// the subnet route tables
	// +optional
	IsPublic bool `json:"isPublic"`

	// Does this subnet map public IPs to instances started in it
	// +nullable
	MapPublicIPOnLaunch *bool `json:"mapPublicIpOnLaunch,omitempty"`

	// The route tables associated with this subnet
	// +mapType=granular
	RouteTables map[string]AwsRouteTable `json:"routeTables"`

	// The internet gateway associated with this subnet
	// +optional
	InternetGateway string `json:"internetGateway"`

	// A map of NAT gateways associated with this subnet
	// +mapType=granular
	// +optional
	NatGateways map[string]string `json:"natGateways"`

	// The tag value to group subnets by
	// +optional
	SubnetSet int `json:"subnetSet"`

	// A map of transit gateways associated with this subnet
	// +mapType=granular
	// +optional
	TransitGateways map[string]TransitGateway `json:"transitGateways"`

	// A map of VPC peering connections associated with this subnet
	// +mapType=granular
	// +optional
	VpcPeeringConnections map[string]PeeringConnection `json:"vpcPeeringConnections"`
}

// AwsRouteTable is an object that holds information about a route table defined in AWS
// +mapType=granular
type AwsRouteTable struct {
	// The associations defined for this route table
	// +listType=map
	// +listMapKey=id
	Associations []AwsAssociation `json:"associations"`

	// ID The route table ID
	// +kubebuilder:validation:Required
	ID string `json:"id"`

	// Is this a public route table. Determined by validating the
	// existence of an internet gateway
	// +optional
	IsPublic bool `json:"isPublic"`

	// The name of the route table
	// +optional
	Name string `json:"name"`

	// The routes defined in this route table
	// +mapType=granular
	Routes map[string]AwsRoute `json:"routes"`

	// The tag value to group route tables by
	// +optional
	SubnetSet int `json:"subnetSet"`
}

// AwsRoute is an object that holds information about a route defined in AWS
// +mapType=granular
type AwsRoute struct {
	// ID The route ID
	// +kubebuilder:validation:Required
	ID string `json:"id"`

	// The destination CIDR block for this route
	// +optional
	DestinationCidrBlock string `json:"destinationCidrBlock"`

	// The gateway ID for this route
	// +optional
	GatewayID string `json:"gatewayId"`

	// The instance ID for this route
	// +optional
	InstanceID string `json:"instanceId"`

	// The NAT gateway ID for this route
	// +optional
	NatGatewayID string `json:"natGatewayId"`

	// The network interface ID for this route
	// +optional
	NetworkInterfaceID string `json:"networkInterfaceId"`

	// The transit gateway ID for this route
	// +optional
	TransitGateway string `json:"transitGateway"`

	// The VPC peering connection ID for this route
	// +optional
	VpcPeeringConnectionID string `json:"vpcPeeringConnectionId"`

	// The local gateway ID for this route
	// +optional
	LocalGatewayID string `json:"localGatewayId"`

	// The carrier gateway ID for this route
	// +optional
	CarrierGatewayID string `json:"carrierGatewayId"`

	// The prefix list ID for this route
	// +optional
	PrefixListID string `json:"prefixListId"`

	// The egress only internet gateway ID for this route
	// +optional
	EgressOnlyInternetGatewayID string `json:"egressOnlyInternetGatewayId"`

	// The destination IPv6 CIDR block for this route
	// +optional
	DestinationIpv6CidrBlock string `json:"destinationIpv6CidrBlock"`
}

// AwsAssociation is an object that holds information about an association defined in AWS
// +mapType=granular
type AwsAssociation struct {
	// ID The association ID
	// +kubebuilder:validation:Required
	ID string `json:"id"`

	// The association state
	// +optional
	State string `json:"state"`

	// The association main flag
	// +optional
	Main bool `json:"main"`

	// The association route table ID
	// +optional
	RouteTableID string `json:"routeTableId"`

	// The association subnet ID
	// +optional
	SubnetID string `json:"subnetId"`
}
