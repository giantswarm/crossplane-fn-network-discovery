package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	xfnaws "github.com/giantswarm/xfnlib/pkg/auth/aws"

	xfnd "github.com/giantswarm/crossplane-fn-network-discovery/pkg/composite/v1beta1"
	inp "github.com/giantswarm/crossplane-fn-network-discovery/pkg/input/v1beta1"
)

const nametag = "Name"

// EC2API Describes the functions required to access data on the AWS EC2 api
type AwsEc2Api interface {
	DescribeVpcs(ctx context.Context,
		params *ec2.DescribeVpcsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error)
	DescribeSubnets(ctx context.Context,
		params *ec2.DescribeSubnetsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error)
	DescribeRouteTables(ctx context.Context,
		params *ec2.DescribeRouteTablesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeRouteTablesOutput, error)
	DescribeSecurityGroups(ctx context.Context,
		params *ec2.DescribeSecurityGroupsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error)
	DescribeNatGateways(ctx context.Context,
		params *ec2.DescribeNatGatewaysInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeNatGatewaysOutput, error)
	DescribeTransitGateways(ctx context.Context,
		params *ec2.DescribeTransitGatewaysInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeTransitGatewaysOutput, error)
	DescribeVpcPeeringConnections(ctx context.Context,
		params *ec2.DescribeVpcPeeringConnectionsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeVpcPeeringConnectionsOutput, error)
}

// Get the EC2 Launch template versions for a given launch template
func GetVpc(c context.Context, api AwsEc2Api, input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return api.DescribeVpcs(c, input)
}

func GetSubnets(c context.Context, api AwsEc2Api, input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return api.DescribeSubnets(c, input)
}

func GetSecurityGroups(c context.Context, api AwsEc2Api, input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	return api.DescribeSecurityGroups(c, input)
}

func GetRouteTables(c context.Context, api AwsEc2Api, input *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
	return api.DescribeRouteTables(c, input)
}

func GetNatGateways(c context.Context, api AwsEc2Api, input *ec2.DescribeNatGatewaysInput) (*ec2.DescribeNatGatewaysOutput, error) {
	return api.DescribeNatGateways(c, input)
}

func GetTransitGateways(c context.Context, api AwsEc2Api, input *ec2.DescribeTransitGatewaysInput) (*ec2.DescribeTransitGatewaysOutput, error) {
	return api.DescribeTransitGateways(c, input)
}

func GetVpcPeeringConnections(c context.Context, api AwsEc2Api, input *ec2.DescribeVpcPeeringConnectionsInput) (*ec2.DescribeVpcPeeringConnectionsOutput, error) {
	return api.DescribeVpcPeeringConnections(c, input)
}

var (
	getEc2Client = func(cfg aws.Config, ep string) AwsEc2Api {
		if ep != "" {
			return ec2.NewFromConfig(cfg, func(o *ec2.Options) {
				o.BaseEndpoint = &ep
			})
		}
		return ec2.NewFromConfig(cfg)
	}

	awsConfig = func(region, provider *string, log logging.Logger) (aws.Config, map[string]string, error) {
		return xfnaws.Config(region, provider, log)
	}
)

// func (f *Function) ReadVpc(vpcName, region, groupTag, providerConfig *string) (vpc xfnd.Vpc, err error) {
func (f *Function) ReadVpc(input *inp.RemoteVpc) (vpc xfnd.AwsVpc, err error) {
	var (
		cfg      aws.Config
		services map[string]string
		vpcInput *ec2.DescribeVpcsInput = &ec2.DescribeVpcsInput{
			Filters: []ec2types.Filter{
				{
					Name:   aws.String("tag:Name"),
					Values: []string{input.Name},
				},
			},
		}
		ec2client AwsEc2Api
	)

	f.log.Info("Reading VPC", "vpc", input.Name, "region", input.Region, "providerConfig", input.ProviderConfig, "groupBy", input.GroupBy)
	// Set up the aws client config
	if cfg, services, err = awsConfig(&input.Region, &input.ProviderConfig, f.log); err != nil {
		err = errors.Wrap(err, "failed to load aws config with region "+input.Region)
		return
	}

	var ep string
	var ok bool
	if _, ok = services["ec2"]; ok {
		ep = services["ec2"]
	}

	f.log.Info("setting up ec2 client to region " + input.Region + " with provider config " + input.ProviderConfig + " and endpoint " + ep)
	ec2client = getEc2Client(cfg, ep)
	vpc, err = f.getVpc(ec2client, vpcInput, &input.GroupBy)
	return
}

func (f *Function) getVpc(client AwsEc2Api, input *ec2.DescribeVpcsInput, groupTag *string) (v xfnd.AwsVpc, err error) {
	var (
		vpcOutput   *ec2.DescribeVpcsOutput
		subnetInput *ec2.DescribeSubnetsInput
	)
	vpcOutput, err = GetVpc(context.Background(), client, input)
	if err != nil {
		fmt.Println("Got an error retrieving information about your VPC endpoint:")
		fmt.Println(err)
		return
	}
	if len(vpcOutput.Vpcs) == 0 {
		err = errors.New("VPC not found")
		return
	}
	f.log.Info("Processing VPC", "vpc", *vpcOutput.Vpcs[0].VpcId)

	subnetInput = &ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{*vpcOutput.Vpcs[0].VpcId},
			},
		},
	}

	var subnets map[string]xfnd.AwsSubnet
	var count int
	{
		count, subnets, err = f.getSubnets(client, subnetInput, groupTag)
		if err != nil {
			return
		}
	}

	var (
		publicSubnets         []xfnd.StatusSubnets     = make([]xfnd.StatusSubnets, count)
		privateSubnets        []xfnd.StatusSubnets     = make([]xfnd.StatusSubnets, count)
		publicRouteTables     []xfnd.StatusRouteTables = make([]xfnd.StatusRouteTables, count)
		privateRouteTables    []xfnd.StatusRouteTables = make([]xfnd.StatusRouteTables, count)
		natGateways           map[string]string        = make(map[string]string, count)
		transitGateways       map[string]string        = make(map[string]string, count)
		vpcPeeringConnections map[string]string        = make(map[string]string, count)
		igw                   string
	)
	{
		for n, sn := range subnets {
			var g int = sn.SubnetSet
			if publicSubnets[g] == nil {
				publicSubnets[g] = make(map[string]string)
			}

			if privateSubnets[g] == nil {
				privateSubnets[g] = make(map[string]string)
			}

			if publicRouteTables[g] == nil {
				publicRouteTables[g] = make(map[string]string)
			}

			if privateRouteTables[g] == nil {
				privateRouteTables[g] = make(map[string]string)
			}

			if sn.IsPublic {
				publicSubnets[g][n] = sn.ID
			} else {
				privateSubnets[g][n] = sn.ID
			}

			if sn.InternetGateway != "" {
				igw = sn.InternetGateway
			}

			if sn.NatGateways != nil {
				for nat, natgw := range sn.NatGateways {
					f.log.Info("Processing NAT Gateway", "nat", nat, "natgw", natgw)
					natGateways[nat] = natgw
				}
			}

			if sn.TransitGateways != nil {
				for tgw, tgwgw := range sn.TransitGateways {
					f.log.Info("Processing Transit Gateway", "tgw", tgw, "tgwgw", tgwgw)
					transitGateways[tgw] = tgwgw
				}
			}

			if sn.VpcPeeringConnections != nil {
				for vp, peering := range sn.VpcPeeringConnections {
					f.log.Info("Processing VPC Peering Connection", "vp", vp, "peering", peering)
					vpcPeeringConnections[vp] = peering
				}
			}

			for n, rt := range sn.RouteTables {
				if rt.IsPublic {
					publicRouteTables[g][n] = rt.ID
				} else {
					privateRouteTables[g][n] = rt.ID
				}
			}
		}
	}

	var securitygroups map[string]string
	{
		securitygroups, err = f.getSecurityGroups(client, *vpcOutput.Vpcs[0].VpcId)
		if err != nil {
			return
		}
	}

	var additionalCidrBlocks []string = make([]string, 0)
	{
		for _, cidr := range vpcOutput.Vpcs[0].CidrBlockAssociationSet {
			if *cidr.CidrBlock != *vpcOutput.Vpcs[0].CidrBlock {
				additionalCidrBlocks = append(additionalCidrBlocks, *cidr.CidrBlock)
			}
		}
	}

	v = xfnd.AwsVpc{
		AdditionalCidrBlocks:  additionalCidrBlocks,
		CidrBlock:             *vpcOutput.Vpcs[0].CidrBlock,
		ID:                    *vpcOutput.Vpcs[0].VpcId,
		PublicSubnets:         resize(publicSubnets),
		PrivateSubnets:        resize(privateSubnets),
		PublicRouteTables:     resize(publicRouteTables),
		PrivateRouteTables:    resize(privateRouteTables),
		InternetGateway:       igw,
		NatGateways:           natGateways,
		SecurityGroups:        securitygroups,
		TransitGateways:       transitGateways,
		VpcPeeringConnections: vpcPeeringConnections,
	}

	return v, nil
}

func resize[T []xfnd.StatusSubnets | []xfnd.StatusRouteTables](s T) T {
	var (
		max int
	)
	switch v := any(s).(type) {
	case []xfnd.StatusSubnets:
		for _, sn := range v {
			if len(sn) > 0 {
				max++
			}
		}
		s = any(v[:max]).(T)
	case []xfnd.StatusRouteTables:
		for _, rt := range v {
			if len(rt) > 0 {
				max++
			}
		}
		s = any(v[:max]).(T)
	}
	return s
}

func (f *Function) getSubnets(client AwsEc2Api, input *ec2.DescribeSubnetsInput, groupTag *string) (count int, subnets map[string]xfnd.AwsSubnet, err error) {
	f.log.Info("Getting subnets")
	subnets = make(map[string]xfnd.AwsSubnet)

	var subnetOutput *ec2.DescribeSubnetsOutput
	{
		subnetOutput, err = GetSubnets(context.Background(), client, input)
		if err != nil {
			return
		}
	}

	var groups map[int]bool = make(map[int]bool)
	groups[0] = true

	for _, sn := range subnetOutput.Subnets {
		var name string
		var subnetSet int = 0

		{
			for _, tag := range sn.Tags {
				if *tag.Key == nametag {
					name = *tag.Value
				}

				if *tag.Key == *groupTag {
					if i, e := strconv.Atoi(*tag.Value); e == nil {
						subnetSet = i
						groups[i] = true
					}
				}

			}
		}

		f.log.Info("Processing subnet", "sn", *sn.SubnetId, "name", name)
		var s xfnd.AwsSubnet = xfnd.AwsSubnet{
			ID:                  *sn.SubnetId,
			AvailabilityZone:    *sn.AvailabilityZone,
			CidrBlock:           *sn.CidrBlock,
			IsPublic:            false,
			IsIpv6:              false,
			MapPublicIPOnLaunch: sn.MapPublicIpOnLaunch,
			SubnetSet:           subnetSet,
		}

		s.RouteTables = make(map[string]xfnd.AwsRouteTable)
		s.NatGateways = make(map[string]string)
		s.TransitGateways = make(map[string]string)
		s.VpcPeeringConnections = make(map[string]string)

		var routeTables *ec2.DescribeRouteTablesOutput
		{
			routeTables, err = GetRouteTables(context.TODO(), client, &ec2.DescribeRouteTablesInput{
				Filters: []ec2types.Filter{
					{
						Name:   aws.String("association.subnet-id"),
						Values: []string{*sn.SubnetId},
					},
				},
			})
			if err != nil {
				f.log.Info("Got an error retrieving information about your route tables", "error", err)
				return
			}
		}

		if len(routeTables.RouteTables) == 0 {
			f.log.Info("No route tables found for subnet", "sn", *sn.SubnetId)
			return 0, nil, errors.New("No route tables found for subnet")
		}

		for _, rt := range routeTables.RouteTables {
			var (
				rtblName     string
				associations []xfnd.AwsAssociation
			)
			{
				for _, tag := range rt.Tags {
					if *tag.Key == nametag {
						rtblName = *tag.Value
					}
				}

				f.log.Info("Processing route table", "rt", *rt.RouteTableId, "name", rtblName)
				if len(rt.Routes) == 0 {
					f.log.Info("No routes found for route table", "rt", *rt.RouteTableId)
					return 0, nil, errors.New("No routes found for route table")
				}

				for _, assoc := range rt.Associations {
					if assoc.SubnetId != nil && *assoc.SubnetId != *sn.SubnetId {
						continue
					}
					f.log.Info("Processing association", "assoc", *assoc.RouteTableAssociationId)
					var a xfnd.AwsAssociation = xfnd.AwsAssociation{
						ID: *assoc.RouteTableAssociationId,
					}

					associations = append(associations, a)
					if assoc.GatewayId != nil && strings.HasPrefix(*assoc.GatewayId, "igw-") {
						s.IsPublic = true
						s.InternetGateway = *assoc.GatewayId
					}
				}

				for _, r := range rt.Routes {
					r := r
					if r.GatewayId != nil && strings.HasPrefix(*r.GatewayId, "igw-") {
						s.IsPublic = true
						s.InternetGateway = *r.GatewayId
					}

					if r.NatGatewayId != nil {
						var ngwname string
						ngwname, err = f.getNatGateway(client, *r.NatGatewayId)
						if err != nil {
							f.log.Info("Error getting NAT Gateway - skipping", "error", err)
						}
						if ngwname != "" {
							s.NatGateways[ngwname] = *r.NatGatewayId
						}
					}

					if r.TransitGatewayId != nil {
						var tgwname string
						tgwname, err = f.getTransitGateway(client, *r.TransitGatewayId)
						if err != nil {
							f.log.Info("Error getting Transit Gateway - skipping", "error", err)
						}
						if tgwname != "" {
							s.TransitGateways[tgwname] = *r.TransitGatewayId
						}

					}

					if r.VpcPeeringConnectionId != nil {
						var pcname string
						pcname, err = f.getVpcPeeringConnection(client, *r.VpcPeeringConnectionId)
						if err != nil {
							f.log.Info("Error getting VPC Peering Connection - skipping", "error", err)
						}

						if pcname != "" {
							s.VpcPeeringConnections[pcname] = *r.VpcPeeringConnectionId
						}

					}
				}
			}

			var rtbl xfnd.AwsRouteTable = xfnd.AwsRouteTable{
				ID:           *rt.RouteTableId,
				Associations: associations,
				IsPublic:     s.IsPublic,
				SubnetSet:    subnetSet,
			}
			rtbl.Routes = make(map[string]xfnd.AwsRoute)
			s.RouteTables[rtblName] = rtbl
		}
		subnets[name] = s
	}

	for _, g := range groups {
		if g {
			count++
		}
	}

	return count, subnets, nil
}

func (f *Function) getNatGateway(client AwsEc2Api, ngwId string) (name string, err error) {
	f.log.Info("Getting NAT Gateway", "ngw", ngwId)
	ngw, err := GetNatGateways(context.Background(), client, &ec2.DescribeNatGatewaysInput{
		NatGatewayIds: []string{ngwId},
	})
	if err != nil {
		return
	}

	for _, n := range ngw.NatGateways {
		for _, tag := range n.Tags {
			if *tag.Key == nametag {
				name = *tag.Value
			}
		}
	}
	return
}

func (f *Function) getTransitGateway(client AwsEc2Api, tgwId string) (name string, err error) {
	f.log.Info("Getting Transit Gateway", "tgw", tgwId)
	tgw, err := GetTransitGateways(context.Background(), client, &ec2.DescribeTransitGatewaysInput{
		TransitGatewayIds: []string{tgwId},
	})
	if err != nil {
		return
	}

	for _, n := range tgw.TransitGateways {
		for _, tag := range n.Tags {
			if *tag.Key == nametag {
				name = *tag.Value
			}
		}
	}
	return
}

func (f *Function) getVpcPeeringConnection(client AwsEc2Api, pcId string) (name string, err error) {
	f.log.Info("Getting VPC Peering Connection", "pc", pcId)
	pc, err := GetVpcPeeringConnections(context.Background(), client, &ec2.DescribeVpcPeeringConnectionsInput{
		VpcPeeringConnectionIds: []string{pcId},
	})
	if err != nil {
		return
	}

	for _, n := range pc.VpcPeeringConnections {
		for _, tag := range n.Tags {
			if *tag.Key == nametag {
				name = *tag.Value
			}
		}
	}
	return
}

func (f *Function) getSecurityGroups(client AwsEc2Api, vpcId string) (sgs map[string]string, err error) {
	f.log.Info("Getting security groups")
	sgs = make(map[string]string)
	securitygroups, err := GetSecurityGroups(context.Background(), client, &ec2.DescribeSecurityGroupsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcId},
			},
		},
	})

	if err != nil {
		return
	}

	for _, sg := range securitygroups.SecurityGroups {
		f.log.Info("Processing security group", "sg", *sg.GroupId)
		sgs[*sg.GroupName] = *sg.GroupId
	}

	return
}
