package commands

import (
	"encoding/json"
	"strings"

	"github.com/tvanriel/cloudsdk/kubernetes"
	"github.com/tvanriel/ps-bot-2/internal/queues"
	"github.com/tvanriel/ps-bot-2/internal/randstr"
	"go.uber.org/zap"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SaverConfiguration struct {
	SecretName string
	BucketName string
}

type SaveCommand struct {
	secretName string
	bucketName string
	kubernetes *kubernetes.KubernetesClient
	log        *zap.Logger
}

func (s *SaveCommand) SkipsPrefix() bool {
	return false
}

func NewSaveCommand(k *kubernetes.KubernetesClient, config *SaverConfiguration, l *zap.Logger) *SaveCommand {
	return &SaveCommand{
		kubernetes: k,
		bucketName: config.BucketName,
		secretName: config.SecretName,
		log:        l,
	}
}

func (s *SaveCommand) Name() string {
	return "save"
}

func (s *SaveCommand) Apply(ctx *Context) error {

	if len(ctx.Args) < 1 {
		ctx.Reply("Usage: save <name> - Saves the attachment as a ps command")
		return nil
	}

	if len(ctx.Message.Attachments) != 1 {
		ctx.Reply("You must provide an attachment")
		return nil
	}

	url := ctx.Message.Attachments[0].URL
	guildId := ctx.Message.GuildID
	soundName := ctx.Args[0]

	s.log.Info("Pushing job to Kubernetes",
		zap.String("url", url),
		zap.String("guildId", guildId),
		zap.String("soundName", soundName),
	)

        amqpBody, err := json.Marshal(queues.QueuedMessage{
                ChannelID: ctx.Message.ChannelID,
                Content: strings.Join([]string{
                        "Saved sound ",
                        soundName,
                        ".",
                }, ""),
        })
        if err != nil {
                return err
        }

	return s.kubernetes.RunJob(convertJob(url, guildId, soundName, s.secretName, s.bucketName, string(amqpBody)))
}

func convertJob(url, namespace, name, secretName, bucketName, amqpBody string) *batchv1.Job {

        id := randstr.Concat(
                randstr.Randstr(randstr.Lowercase, 1),
                randstr.Randstr(randstr.Lowercase + randstr.Numbers, 5),
                )
        

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ps-bot-download-" + namespace + id,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "downloader",
							Image: "mitaka8/dca-encoder:latest",
							Env: []v1.EnvVar{
								{
									Name:  "SOURCE",
									Value: url,
								},
								{
									Name:  "TARGET",
									Value: bucketName + "/" + namespace,
								},
								{
									Name:  "FILENAME",
									Value: name,
								},
								{
									Name:  "POST_HOOK",
									Value: "1",
								},
                                                                {
                                                                        Name: "AMQP_BODY",
                                                                        Value: amqpBody,
                                                                },
							},
							EnvFrom: []v1.EnvFromSource{{
								SecretRef: &v1.SecretEnvSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: secretName,
									},
								},
							}},
						},
					},
					RestartPolicy: v1.RestartPolicyOnFailure,
				},
			},
		},
	}
}
