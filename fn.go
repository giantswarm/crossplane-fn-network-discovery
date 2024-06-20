package main

import (
	"context"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/response"
	"github.com/giantswarm/xfnlib/pkg/composite"
	"k8s.io/apimachinery/pkg/runtime"

	fnc "github.com/giantswarm/crossplane-fn-network-discovery/pkg/composite/v1beta1"
	inp "github.com/giantswarm/crossplane-fn-network-discovery/pkg/input/v1beta1"
)

const composedName = "crossplane-fn-network-discovery"

// RunFunction runs the composition Function to generate subnets from the given cluster
func (f *Function) RunFunction(_ context.Context, req *fnv1beta1.RunFunctionRequest) (rsp *fnv1beta1.RunFunctionResponse, err error) {
	f.log.Info("preparing function", composedName, req.GetMeta().GetTag())
	rsp = response.To(req, response.DefaultTTL)

	var (
		composed       *composite.Composition
		input          inp.Input
		vpcs           map[string]fnc.Vpc = make(map[string]fnc.Vpc)
		search         []inp.RemoteVpc    = make([]inp.RemoteVpc, 0)
		region         string
		providerConfig string
	)

	// The composite resource that actually exists.
	oxr, err := request.GetObservedCompositeResource(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrap(err, "cannot get observed composite resource"))
		return rsp, nil
	}

	if composed, err = composite.New(req, &input, &oxr); err != nil {
		response.Fatal(rsp, errors.Wrap(err, "error setting up function "+composedName))
		return rsp, nil
	}

	if input.Spec == nil {
		response.Fatal(rsp, &composite.MissingSpec{})
		return rsp, nil
	}

	enabled, err := f.getBooleanFromPaved(oxr.Resource, input.Spec.EnabledRef)
	if err != nil {
		f.log.Info("cannot get enabled state from input", "error", err)
		response.Fatal(rsp, errors.Wrap(err, "cannot get enabled state from input"))
		return rsp, nil
	}

	if !enabled {
		f.log.Info("Function is disabled for composition")
		return rsp, nil
	}

	if err = f.getStringFromPaved(oxr.Resource, input.Spec.RegionRef, &region); err != nil {
		f.log.Info("cannot get region from input", "error", err)
		response.Fatal(rsp, errors.Wrap(err, "cannot get region from input"))
		return rsp, nil
	}
	f.log.Info("Region", "region", region)

	var groupTag string
	_ = f.getStringFromPaved(oxr.Resource, input.Spec.GroupByRef, &groupTag)

	if err = f.getStringFromPaved(oxr.Resource, input.Spec.ProviderConfigRef, &providerConfig); err != nil {
		f.log.Info("cannot get provider config from input", "error", err)
		response.Fatal(rsp, errors.Wrap(err, "cannot get provider config from input"))
		return rsp, nil
	}
	f.log.Info("ProviderConfig", "pc", providerConfig)

	if err = f.getValueInto(oxr.Resource, input.Spec.VpcNameRef, region, providerConfig, groupTag, &search); err != nil {
		f.log.Info("cannot get VPC name from input", "error", err)
		response.Fatal(rsp, errors.Wrap(err, "cannot get VPC name from input"))
		return rsp, nil
	}

	for _, n := range search {
		var vpc fnc.Vpc
		if vpc, err = f.ReadVpc(&n); err != nil {
			f.log.Info("cannot read VPC", "error", err, "name", n.Name, "region", n.Region, "providerConfig", n.ProviderConfig)
			continue
		}

		// Copy the  provider config and region from the search input so the
		// composition doesn't have to re-match it on cross-account lookups.
		vpc.Region = n.Region
		vpc.ProviderConfig = n.ProviderConfig
		vpcs[n.Name] = vpc
	}
	f.log.Info("VPCs", "vpcs", vpcs)

	if err = f.patchFieldValueToObject(input.Spec.PatchTo, vpcs, composed.DesiredComposite.Resource); err != nil {
		f.log.Info("cannot patch VPCs to composite", "error", err)
		response.Fatal(rsp, errors.Wrapf(err, "cannot render ToComposite patch %q", input.Spec.PatchTo))
		return rsp, nil
	}

	if err = composed.ToResponse(rsp); err != nil {
		f.log.Info("cannot convert composition to response", "error", err)
		response.Fatal(rsp, errors.Wrapf(err, "cannot convert composition to response %T", rsp))
		return
	}

	return rsp, nil
}

// get array from paved
func (f *Function) getValueInto(req runtime.Object, ref, region, providerConfig, groupBy string, value *[]inp.RemoteVpc) (err error) {
	var paved *fieldpath.Paved
	if paved, err = fieldpath.PaveObject(req); err != nil {
		return
	}

	var s string
	if s, err = paved.GetString(ref); err != nil {
		err = paved.GetValueInto(ref, &value)
		for i := range *value {
			if (*value)[i].Region == "" {
				(*value)[i].Region = region
			}

			if (*value)[i].ProviderConfig == "" {
				(*value)[i].ProviderConfig = providerConfig
			}

			if (*value)[i].GroupBy == "" {
				(*value)[i].GroupBy = groupBy
			}
		}
		return
	}
	input := inp.RemoteVpc{
		GroupBy:        groupBy,
		Name:           s,
		Region:         region,
		ProviderConfig: providerConfig,
	}
	*value = append(*value, input)
	return
}

// get string from paved
func (f *Function) getStringFromPaved(req runtime.Object, ref string, value *string) (err error) {
	var paved *fieldpath.Paved
	if paved, err = fieldpath.PaveObject(req); err != nil {
		return
	}

	*value, err = paved.GetString(ref)
	return
}

// get boolean from paved
func (f *Function) getBooleanFromPaved(req runtime.Object, ref string) (value bool, err error) {
	var paved *fieldpath.Paved
	if paved, err = fieldpath.PaveObject(req); err != nil {
		return
	}

	value, err = paved.GetBool(ref)
	return
}

// patchFieldValueToObject is used to push information onto the XR status
func (f *Function) patchFieldValueToObject(fieldPath string, value map[string]fnc.Vpc, to runtime.Object) (err error) {
	var paved *fieldpath.Paved
	if paved, err = fieldpath.PaveObject(to); err != nil {
		return
	}

	if err = paved.SetValue(fieldPath, value); err != nil {
		return
	}

	return runtime.DefaultUnstructuredConverter.FromUnstructured(paved.UnstructuredContent(), to)
}
