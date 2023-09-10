package soundstore

import (
	"context"
	"path/filepath"
	"strings"

	"io"

	"github.com/minio/minio-go/v7"
)

type SoundStore struct {
	s3     *minio.Client
	bucket string
}

func NewSoundStore(s3 *minio.Client, config *Configuration) *SoundStore {

	return &SoundStore{
		s3:     s3,
		bucket: config.Bucket,
	}

}

func (s *SoundStore) Find(guildId string, sound string) (io.ReadCloser, error) {
	return s.s3.GetObject(context.Background(), s.bucket, dcafile(guildId, sound), minio.GetObjectOptions{})
}

func (s *SoundStore) List(guildId string) []string {

	names := []string{}
	for object := range s.s3.ListObjects(context.Background(), s.bucket, minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    guildId,
	}) {
		if object.Err != nil {
			continue
		}

		names = append(names, strings.TrimSuffix(filepath.Base(object.Key), ".dca"))
	}

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
