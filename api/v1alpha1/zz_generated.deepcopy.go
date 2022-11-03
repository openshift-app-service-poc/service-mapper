//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespacedName) DeepCopyInto(out *NamespacedName) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespacedName.
func (in *NamespacedName) DeepCopy() *NamespacedName {
	if in == nil {
		return nil
	}
	out := new(NamespacedName)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProxableService) DeepCopyInto(out *ProxableService) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProxableService.
func (in *ProxableService) DeepCopy() *ProxableService {
	if in == nil {
		return nil
	}
	out := new(ProxableService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ProxableService) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProxableServiceList) DeepCopyInto(out *ProxableServiceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ProxableService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProxableServiceList.
func (in *ProxableServiceList) DeepCopy() *ProxableServiceList {
	if in == nil {
		return nil
	}
	out := new(ProxableServiceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ProxableServiceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProxableServiceSpec) DeepCopyInto(out *ProxableServiceSpec) {
	*out = *in
	out.ServiceResourceMapRef = in.ServiceResourceMapRef
	out.ServiceInstance = in.ServiceInstance
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProxableServiceSpec.
func (in *ProxableServiceSpec) DeepCopy() *ProxableServiceSpec {
	if in == nil {
		return nil
	}
	out := new(ProxableServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProxableServiceStatus) DeepCopyInto(out *ProxableServiceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProxableServiceStatus.
func (in *ProxableServiceStatus) DeepCopy() *ProxableServiceStatus {
	if in == nil {
		return nil
	}
	out := new(ProxableServiceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceKindReference) DeepCopyInto(out *ServiceKindReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceKindReference.
func (in *ServiceKindReference) DeepCopy() *ServiceKindReference {
	if in == nil {
		return nil
	}
	out := new(ServiceKindReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceProxy) DeepCopyInto(out *ServiceProxy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceProxy.
func (in *ServiceProxy) DeepCopy() *ServiceProxy {
	if in == nil {
		return nil
	}
	out := new(ServiceProxy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceProxy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceProxyList) DeepCopyInto(out *ServiceProxyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ServiceProxy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceProxyList.
func (in *ServiceProxyList) DeepCopy() *ServiceProxyList {
	if in == nil {
		return nil
	}
	out := new(ServiceProxyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceProxyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceProxySpec) DeepCopyInto(out *ServiceProxySpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceProxySpec.
func (in *ServiceProxySpec) DeepCopy() *ServiceProxySpec {
	if in == nil {
		return nil
	}
	out := new(ServiceProxySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceProxyStatus) DeepCopyInto(out *ServiceProxyStatus) {
	*out = *in
	out.Binding = in.Binding
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceProxyStatus.
func (in *ServiceProxyStatus) DeepCopy() *ServiceProxyStatus {
	if in == nil {
		return nil
	}
	out := new(ServiceProxyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceProxyStatusBinding) DeepCopyInto(out *ServiceProxyStatusBinding) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceProxyStatusBinding.
func (in *ServiceProxyStatusBinding) DeepCopy() *ServiceProxyStatusBinding {
	if in == nil {
		return nil
	}
	out := new(ServiceProxyStatusBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceReference) DeepCopyInto(out *ServiceReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceReference.
func (in *ServiceReference) DeepCopy() *ServiceReference {
	if in == nil {
		return nil
	}
	out := new(ServiceReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceResourceMap) DeepCopyInto(out *ServiceResourceMap) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceResourceMap.
func (in *ServiceResourceMap) DeepCopy() *ServiceResourceMap {
	if in == nil {
		return nil
	}
	out := new(ServiceResourceMap)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceResourceMap) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceResourceMapList) DeepCopyInto(out *ServiceResourceMapList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ServiceResourceMap, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceResourceMapList.
func (in *ServiceResourceMapList) DeepCopy() *ServiceResourceMapList {
	if in == nil {
		return nil
	}
	out := new(ServiceResourceMapList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceResourceMapList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceResourceMapSpec) DeepCopyInto(out *ServiceResourceMapSpec) {
	*out = *in
	out.ServiceKindReference = in.ServiceKindReference
	if in.ServiceMap != nil {
		in, out := &in.ServiceMap, &out.ServiceMap
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceResourceMapSpec.
func (in *ServiceResourceMapSpec) DeepCopy() *ServiceResourceMapSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceResourceMapSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceResourceMapStatus) DeepCopyInto(out *ServiceResourceMapStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceResourceMapStatus.
func (in *ServiceResourceMapStatus) DeepCopy() *ServiceResourceMapStatus {
	if in == nil {
		return nil
	}
	out := new(ServiceResourceMapStatus)
	in.DeepCopyInto(out)
	return out
}
