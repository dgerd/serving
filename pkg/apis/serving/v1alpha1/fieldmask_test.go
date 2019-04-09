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
	"testing"

	"github.com/knative/pkg/kmp"
	corev1 "k8s.io/api/core/v1"
)

func TestVolumeMask(t *testing.T) {
	want := &corev1.Volume{
		Name:         "foo",
		VolumeSource: corev1.VolumeSource{},
	}
	in := want

	got := VolumeMask(in)

	if &want == &got {
		t.Errorf("Input and output share addresses. Want different addresses")
	}

	if diff, err := kmp.SafeDiff(want, got); err != nil {
		t.Errorf("Got error comparing output, err = %v", err)
	} else if diff != "" {
		t.Errorf("VolumeMask (-want, +got): %s", diff)
	}
}

func TestSourceMask(t *testing.T) {
	want := &corev1.VolumeSource{
		Secret:    &corev1.SecretVolumeSource{},
		ConfigMap: &corev1.ConfigMapVolumeSource{},
	}
	in := &corev1.VolumeSource{
		Secret:    &corev1.SecretVolumeSource{},
		ConfigMap: &corev1.ConfigMapVolumeSource{},
		NFS:       &corev1.NFSVolumeSource{},
	}

	got := VolumeSourceMask(in)

	if &want == &got {
		t.Errorf("Input and output share addresses. Want different addresses")
	}

	if diff, err := kmp.SafeDiff(want, got); err != nil {
		t.Errorf("Got error comparing output, err = %v", err)
	} else if diff != "" {
		t.Errorf("VolumeMask (-want, +got): %s", diff)
	}
}

func TestContainerMask(t *testing.T) {
	want := &corev1.Container{
		Args:                     []string{"hello"},
		Command:                  []string{"world"},
		Env:                      []corev1.EnvVar{{}},
		EnvFrom:                  []corev1.EnvFromSource{{}},
		Image:                    "python",
		LivenessProbe:            &corev1.Probe{},
		Ports:                    []corev1.ContainerPort{{}},
		ReadinessProbe:           &corev1.Probe{},
		Resources:                corev1.ResourceRequirements{},
		SecurityContext:          &corev1.SecurityContext{},
		TerminationMessagePath:   "/",
		TerminationMessagePolicy: corev1.TerminationMessageReadFile,
		VolumeMounts:             []corev1.VolumeMount{{}},
	}
	in := &corev1.Container{
		Args:                     []string{"hello"},
		Command:                  []string{"world"},
		Env:                      []corev1.EnvVar{{}},
		EnvFrom:                  []corev1.EnvFromSource{{}},
		Image:                    "python",
		LivenessProbe:            &corev1.Probe{},
		Ports:                    []corev1.ContainerPort{{}},
		ReadinessProbe:           &corev1.Probe{},
		Resources:                corev1.ResourceRequirements{},
		SecurityContext:          &corev1.SecurityContext{},
		TerminationMessagePath:   "/",
		TerminationMessagePolicy: corev1.TerminationMessageReadFile,
		VolumeMounts:             []corev1.VolumeMount{{}},
		Name:                     "foo",
		Stdin:                    true,
		StdinOnce:                true,
		TTY:                      true,
	}

	got := ContainerMask(in)

	if &want == &got {
		t.Errorf("Input and output share addresses. Want different addresses")
	}

	if diff, err := kmp.SafeDiff(want, got); err != nil {
		t.Errorf("Got error comparing output, err = %v", err)
	} else if diff != "" {
		t.Errorf("VolumeMask (-want, +got): %s", diff)
	}

}
