# function-network-discovery

A [Crossplane] Composition Function for discovery of VPC architecture

Provide this function with a VPC name, region and provider configuration to use
and it will discover components of the VPC including:

- VPC ID
- CIDR Block
- Additional CIDRS
- Subnets
- Route Tables
- Internet Gateway
- NAT Gateways
- VPC Peering connections
- Transit gateways
- Security groups

This information will then be patched to the status of the XR. To understand the
structure required for the XR status, see [package/composite](./package/composite/)

## Composition integration

This function is placed in the pipeline with a reference to the cluster object
for that composition and an additional reference of where to patch information
about the subnets it is generating for that provider.

This should be specified in your composition, for example

```yaml
  - step: network-discovery
    functionRef:
      name: function-network-discovery
    input:
      apiVersion: nd.fn.giantswarm.io
      kind: Input
      metadata:
        namespace: crossplane
      spec:
        vpcNameRef: spec.vpcs
        regionRef: spec.region
        providerConfigRef: spec.providerConfigRef.name
        patchTo: status.vpcs
```

## Input parameters

- `enabledRef` **optional** Reference to a boolean parameter that optionally
  tells the function to skip discovery. Use this in complex composition
  structures where discovery may or may not be required.
- `providerConfigRef` **required** A reference to an AWS providerConfig
- `regionRef` **required** The default region being used by the XR
- `vpcNameRef` **required** a path to a location on the XR containing the name
  of one or more VPCs. The referenced location may be a single string or a list
  of objects
- `groupingTagRef` **optional** If specified, the location of the reference
  will be used as a tag filter for grouping subnets and route tables together

### vpcNameRef

If the location pointed to by `vpcNameRef` is a list, it must match the
following format:

- `name` **required** The name of the VPC to discover
- `region` **optional** The region to discover the VPC in - if not defined falls
  back to the default region specified above
- `providerConfigRef` **optional** A provider config reference to use for
  discovery of this specific VPC. Useful for cross account VPC discovery

### groupingTagRef

The location for `groupingTagRef` should match the following format:

- `key` string the key for the tag

The value of the tag key on the AWS resource should be numeric. If it is not it
is ignored.

```yaml
tags:
  subnetsets.xnetworks.crossplane.giantswarm.io: 1
```

> [!NOTE]
> This is **not** an AWS tag filter. It is used to group the output of subnets
> and route tables into sets were defined together. If not defined, a single
> list will be output
>
> eg.
>
> ```yaml
> subnets:
> - subnet-1: sn-123456
>   subnet-2: sn-234567
>   ...
>   subnet-10: sn-012345
> ```
>
> If defined, this would otherwise result in the following:
>
> ```yaml
> subnets:
> - subnet-1: sn-123456
>   subnet-2: sn-234567
>   subnet-3: sn-345678
> - subnet-4: sn-456789
>   subnet-5: sn-567890
>   subnet-6: sn-678901
> ```

For information such as transit gateways, nat gateways and peering connections
a unique name tag is expected to prevent resources overwriting each other.

If a name tag cannot be found, the ID will not be returned for that item so if
you are expecting an id to be returned when it isn't appearing in the status,
check that a unique name tag is assigned to the resource in AWS.
