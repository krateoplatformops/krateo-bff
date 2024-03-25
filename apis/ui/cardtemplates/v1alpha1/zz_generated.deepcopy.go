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
	runtime "k8s.io/apimachinery/pkg/runtime"
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
func (in *Card) DeepCopyInto(out *Card) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Card.
func (in *Card) DeepCopy() *Card {
	if in == nil {
		return nil
	}
	out := new(Card)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CardTemplate) DeepCopyInto(out *CardTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CardTemplate.
func (in *CardTemplate) DeepCopy() *CardTemplate {
	if in == nil {
		return nil
	}
	out := new(CardTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CardTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CardTemplateList) DeepCopyInto(out *CardTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CardTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CardTemplateList.
func (in *CardTemplateList) DeepCopy() *CardTemplateList {
	if in == nil {
		return nil
	}
	out := new(CardTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CardTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CardTemplateSpec) DeepCopyInto(out *CardTemplateSpec) {
	*out = *in
	out.FormTemplateRef = in.FormTemplateRef
	out.App = in.App
	if in.Iterator != nil {
		in, out := &in.Iterator, &out.Iterator
		*out = new(string)
		**out = **in
	}
	if in.APIList != nil {
		in, out := &in.APIList, &out.APIList
		*out = make([]*core.API, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(core.API)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CardTemplateSpec.
func (in *CardTemplateSpec) DeepCopy() *CardTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(CardTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CardTemplateStatus) DeepCopyInto(out *CardTemplateStatus) {
	*out = *in
	if in.Cards != nil {
		in, out := &in.Cards, &out.Cards
		*out = make([]*Card, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Card)
				**out = **in
			}
		}
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

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CardTemplateStatus.
func (in *CardTemplateStatus) DeepCopy() *CardTemplateStatus {
	if in == nil {
		return nil
	}
	out := new(CardTemplateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FormTemplateRef) DeepCopyInto(out *FormTemplateRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FormTemplateRef.
func (in *FormTemplateRef) DeepCopy() *FormTemplateRef {
	if in == nil {
		return nil
	}
	out := new(FormTemplateRef)
	in.DeepCopyInto(out)
	return out
}
