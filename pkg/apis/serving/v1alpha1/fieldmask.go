/*
Copyright 2019 The Knative Authors

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
	corev1 "k8s.io/api/core/v1"
)

// KnativeContainer performs a deep copy of the Kubernetes Container object to a new
// Kubernetes Container object bringing over only the fields allowed in the Knative API. This
// does not validate the contents or the bounds of the provided fields.
func KnativeContainer(in *corev1.Container) *corev1.Container {
	if in == nil {
		return nil
	}
	out := new(corev1.Container)
	// Args is allowed
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	// Command is allowed
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	// Env is allowed with restricted fields
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]corev1.EnvVar, len(*in))
		for i := range *in {
			knativeEnvVar(&(*in)[i], &(*out)[i])
		}
	}
	// EnvFrom is allowed
	if in.EnvFrom != nil {
		in, out := &in.EnvFrom, &out.EnvFrom
		*out = make([]corev1.EnvFromSource, len(*in))
		for i := range *in {
			//TODO: Should we really still DeepCopy versus copying each field?
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	// Image is allowed
	out.Image = in.Image
	// ImagePullPolicy is not allowed
	// Lifecycle is not allowed
	// LivenessProbe is allowed with restricted fields
	if in.LivenessProbe != nil {
		in, out := &in.LivenessProbe, &out.LivenessProbe
		*out = new(corev1.Probe)
		(*in).DeepCopyInto(*out)
	}
	// Name is not allowed
	// Ports is allowed with restricted fields
	// TODO: Enforce restrictions
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]corev1.ContainerPort, len(*in))
		copy(*out, *in)
	}
	// ReadinessProbe is allowed with restricted fields
	// TODO: Enforce restrictions
	if in.ReadinessProbe != nil {
		in, out := &in.ReadinessProbe, &out.ReadinessProbe
		*out = new(corev1.Probe)
		(*in).DeepCopyInto(*out)
	}
	// Resource is allowed
	in.Resources.DeepCopyInto(&out.Resources)
	// SecurityContext is allowed with restricted fields
	if in.SecurityContext != nil {
		in, out := &in.SecurityContext, &out.SecurityContext
		*out = new(corev1.SecurityContext)
		(*in).DeepCopyInto(*out)
	}
	// Stdin is not allowed
	// StdinOnce is not allowed
	// TerminationMessagePath is allowed
	out.TerminationMessagePath = in.TerminationMessagePath
	// TerminationMessagePolicy is allowed
	out.TerminationMessagePolicy = in.TerminationMessagePolicy
	// TTY is not allowed
	// VolumeDevices is not allowed
	// VolumeMounts is allowed with restricted fields
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]corev1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return out
}

func knativeEnvVar(in, out *corev1.EnvVar) {
	out.Name = in.Name
	out.Value = in.Value
	if in.ValueFrom != nil {
		in, out := &in.ValueFrom, &out.ValueFrom
		*out = new(corev1.EnvVarSource)
		knativeEnvVarSource(*in, *out)
	}
	return
}

func knativeEnvVarSource(in, out *corev1.EnvVarSource) {
	// FieldRef is not allowed
	out.FieldRef = nil
	// ResourceFieldRef is not allowed
	out.ResourceFieldRef = nil
	if in.ConfigMapKeyRef != nil {
		in, out := &in.ConfigMapKeyRef, &out.ConfigMapKeyRef
		*out = new(corev1.ConfigMapKeySelector)
		(*in).DeepCopyInto(*out)
	}
	if in.SecretKeyRef != nil {
		in, out := &in.SecretKeyRef, &out.SecretKeyRef
		*out = new(corev1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
	return
}
