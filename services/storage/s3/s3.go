package s3

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrusme/hyperuplink/errs"
	"github.com/mrusme/hyperuplink/services/config"
	"github.com/rhnvrm/simples3"
)

type S3 struct {
	cfg    config.Storage
	client *simples3.S3
}

func New(cfg config.Storage) (st *S3, err error) {
	st = new(S3)

	st.cfg = cfg

	return st, nil
}

func (st *S3) Startup() (err error) {
	st.client = simples3.New(
		st.cfg.S3.Region,
		st.cfg.S3.AccessKey,
		st.cfg.S3.SecretKey,
	)
	if st.cfg.S3.Endpoint != "" {
		st.client.SetEndpoint(st.cfg.S3.Endpoint)
	}
	return nil
}

func (st *S3) Shutdown() (err error) {
	return nil
}

func (st *S3) getBucketObjectFile(dest string) (bucket string, objKey string, fileName string, err error) {
	separator := string(os.PathSeparator)
	cleanDest := filepath.Clean(dest)
	destSplit := strings.Split(cleanDest, separator)
	destSplitLen := len(destSplit)
	if destSplitLen < 2 {
		return bucket, objKey, fileName, errs.ErrFilePathInvalid
	}
	bucket = destSplit[0]
	objKey = strings.Join(destSplit[1:], "/")
	fileName = destSplit[(destSplitLen - 1)]

	return bucket, objKey, fileName, nil
}

func (st *S3) getFileContentType(file io.ReadSeeker) (contentType string, err error) {
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	contentType = http.DetectContentType(buffer)
	return contentType, nil
}

func (st *S3) StoreFileName(src string, dest string) (err error) {
	if src == "" || dest == "" {
		return errs.ErrFilePathInvalid
	}

	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	return st.StoreFile(file, dest)
}

func (st *S3) StoreFile(src io.ReadSeeker, dest string) (err error) {
	if src == nil || dest == "" {
		return errs.ErrFilePathInvalid
	}

	var bucket, objKey, fileName string
	if bucket, objKey, fileName, err = st.getBucketObjectFile(dest); err != nil {
		return err
	}

	var contentType string
	if contentType, err = st.getFileContentType(src); err != nil {
		return err
	}

	_, err = st.client.FileUpload(simples3.UploadInput{
		Bucket:      bucket,
		ObjectKey:   objKey,
		FileName:    fileName,
		ContentType: contentType,
		Body:        src,
	})

	return err
}

func (st *S3) GetFileDownloadURL(dest string) (dlurl string, err error) {
	if dest == "" {
		return dlurl, errs.ErrFilePathInvalid
	}

	if st.cfg.S3.PublicDownload {
		return fmt.Sprintf("%s/%s", st.cfg.S3.PublicURL, dest), nil
	}

	if st.cfg.S3.PresignedDownload {
		var bucket, objKey string
		if bucket, objKey, _, err = st.getBucketObjectFile(dest); err != nil {
			return dlurl, err
		}

		dlurl = st.client.GeneratePresignedURL(simples3.PresignedInput{
			Bucket:        bucket,
			ObjectKey:     objKey,
			Method:        "GET",
			ExpirySeconds: 60, // TODO: Make configurable
		})

		return dlurl, nil
	}

	// TODO: If neither public nor presigned downloads are possible, download the
	// file temporarily into Redis and provide a temporary download URL that lets
	// the client get it through a special Route. Set a TTL in Redis.
	return dlurl, errs.ErrNotImplemented
}
