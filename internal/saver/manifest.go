package saver

import (
	"github.com/tvanriel/ps-bot-2/internal/randstr"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func convertJob(params SaveParams, config *Configuration) *batchv1.Job {

	id := randstr.Concat(
		randstr.Randstr(randstr.Lowercase, 1),
		randstr.Randstr(randstr.Lowercase+randstr.Numbers, 5),
	)

	backoffLimit := int32(6)
	ttlSeconds := int32(3600 * 24 * 2)

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ps-bot-download-" + params.GuildID + id,
			Labels: map[string]string{
				"app.kubernetes.io/part-of":   "ps-bot",
				"app.kubernetes.io/name":      "ps-bot",
				"app.kubernetes.io/component": "downloader",
				"app.kubernetes.io/instance":  "downloader",
				"app.kubernetes.io/version":   "latest",
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            &backoffLimit,
			TTLSecondsAfterFinished: &ttlSeconds,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
					Containers: []v1.Container{
						{
							Name:  "downloader",
							Image: "mitaka8/dca-encoder:latest",
							Env: []v1.EnvVar{
								{
									Name:  "SOURCE",
									Value: params.URL,
								},
								{
									Name:  "TARGET",
									Value: params.Target(config),
								},
								{
									Name:  "FILENAME",
									Value: params.SoundName,
								},
								{
									Name:  "REDIS_PAYLOAD",
									Value: params.TextMessage,
								},
							},
							EnvFrom: []v1.EnvFromSource{{
								SecretRef: &v1.SecretEnvSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: config.SecretName,
									},
								},
							}},
						},
					},
				},
			},
		},
	}
}
