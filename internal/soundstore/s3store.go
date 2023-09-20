package soundstore

import (
	"context"
	"path/filepath"
	"strings"

	"io"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type SoundStore struct {
	s3     *minio.Client
	bucket string
        log *zap.Logger
}

func NewSoundStore(s3 *minio.Client, config *Configuration, l *zap.Logger) *SoundStore {

	return &SoundStore{
		s3:     s3,
		bucket: config.Bucket,
                log: l,
	}

}

func (s *SoundStore) Find(guildId string, sound string) (io.ReadCloser, error) {
	return s.s3.GetObject(context.Background(), s.bucket, dcafile(guildId, sound), minio.GetObjectOptions{})
}

func (s *SoundStore) List(guildId string) []string {
        s.log.Info("Requesting list from Soundstore",
                zap.String("guild", guildId),
        )

	names := []string{}
	for object := range s.s3.ListObjects(context.Background(), s.bucket, minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    guildId,
	}) {
		if object.Err != nil {
                        s.log.Error("error listing file", 
                                zap.String("key", object.Key), 
                                zap.String("guild", guildId), 
                                zap.Error(object.Err),
                        )
			continue
		}

		names = append(names, strings.TrimSuffix(filepath.Base(object.Key), ".dca"))
	}

        s.log.Info("List",
                zap.String("guild", guildId),
                zap.Int("count", len(names)),
        )

	return names

}

func dcafile(guildId string, file string) string {
	return filepath.Join(
		guildId,
		strings.Join(
			[]string{
				file,
				".dca",
			}, "",
		),
	)
}
