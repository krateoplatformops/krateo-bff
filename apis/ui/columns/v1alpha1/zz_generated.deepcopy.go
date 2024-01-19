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
	cardtemplatesv1alpha1 "github.com/krateoplatformops/krateo-bff/apis/ui/cardtemplates/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *App) DeepCopyInto(out *App) {
	*out = *in
	if in.Props != nil {
		in, out := &in.Props, &out.Props
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new App.
func (in *App) DeepCopy() *App {
	if in == nil {
		return nil
	}
	out := new(App)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Column) DeepCopyInto(out *Column) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Column.
func (in *Column) DeepCopy() *Column {
	if in == nil {
		return nil
	}
	out := new(Column)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Column) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ColumnList) DeepCopyInto(out *ColumnList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Column, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ColumnList.
func (in *ColumnList) DeepCopy() *ColumnList {
	if in == nil {
		return nil
	}
	out := new(ColumnList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ColumnList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ColumnSpec) DeepCopyInto(out *ColumnSpec) {
	*out = *in
	in.App.DeepCopyInto(&out.App)
	if in.CardTemplateListRef != nil {
		in, out := &in.CardTemplateListRef, &out.CardTemplateListRef
		*out = make([]*core.Reference, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(core.Reference)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ColumnSpec.
func (in *ColumnSpec) DeepCopy() *ColumnSpec {
	if in == nil {
		return nil
	}
	out := new(ColumnSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ColumnStatus) DeepCopyInto(out *ColumnStatus) {
	*out = *in
	if in.Cards != nil {
		in, out := &in.Cards, &out.Cards
		*out = make([]*cardtemplatesv1alpha1.Card, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(cardtemplatesv1alpha1.Card)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ColumnStatus.
func (in *ColumnStatus) DeepCopy() *ColumnStatus {
	if in == nil {
		return nil
	}
	out := new(ColumnStatus)
	in.DeepCopyInto(out)
	return out
}
