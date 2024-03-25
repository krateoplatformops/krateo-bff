//go:build !ignore_autogenerated

/*
Copyright 2023 Krateo SRL.

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
	"github.com/krateoplatformops/krateo-bff/apis/core"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Action) DeepCopyInto(out *Action) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Action.
func (in *Action) DeepCopy() *Action {
	if in == nil {
		return nil
	}
	out := new(Action)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataItem) DeepCopyInto(out *DataItem) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataItem.
func (in *DataItem) DeepCopy() *DataItem {
	if in == nil {
		return nil
	}
	out := new(DataItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FormTemplate) DeepCopyInto(out *FormTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FormTemplate.
func (in *FormTemplate) DeepCopy() *FormTemplate {
	if in == nil {
		return nil
	}
	out := new(FormTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FormTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FormTemplateList) DeepCopyInto(out *FormTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FormTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FormTemplateList.
func (in *FormTemplateList) DeepCopy() *FormTemplateList {
	if in == nil {
		return nil
	}
	out := new(FormTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FormTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FormTemplateSpec) DeepCopyInto(out *FormTemplateSpec) {
	*out = *in
	if in.SchemaDefinitionRef != nil {
		in, out := &in.SchemaDefinitionRef, &out.SchemaDefinitionRef
		*out = new(core.Reference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FormTemplateSpec.
func (in *FormTemplateSpec) DeepCopy() *FormTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(FormTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FormTemplateStatus) DeepCopyInto(out *FormTemplateStatus) {
	*out = *in
	if in.Content != nil {
		in, out := &in.Content, &out.Content
		*out = new(FormTemplateStatusContent)
		(*in).DeepCopyInto(*out)
	}
	if in.Actions != nil {
		in, out := &in.Actions, &out.Actions
		*out = make([]*Action, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Action)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FormTemplateStatus.
func (in *FormTemplateStatus) DeepCopy() *FormTemplateStatus {
	if in == nil {
		return nil
	}
	out := new(FormTemplateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FormTemplateStatusContent) DeepCopyInto(out *FormTemplateStatusContent) {
	*out = *in
	if in.Schema != nil {
		in, out := &in.Schema, &out.Schema
		*out = new(runtime.RawExtension)
		(*in).DeepCopyInto(*out)
	}
	if in.Instance != nil {
		in, out := &in.Instance, &out.Instance
		*out = new(runtime.RawExtension)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FormTemplateStatusContent.
func (in *FormTemplateStatusContent) DeepCopy() *FormTemplateStatusContent {
	if in == nil {
		return nil
	}
	out := new(FormTemplateStatusContent)
	in.DeepCopyInto(out)
	return out
}
