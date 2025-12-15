package s3

import (
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

func (st *S3) StoreFile(src string, dest string) (err error) {
	if src == "" || dest == "" {
		return errs.ErrFilePathInvalid
	}

	destSplit := filepath.SplitList(dest)
	destSplitLen := len(destSplit)
	if destSplitLen < 2 {
		return errs.ErrFilePathInvalid
	}
	bucket := destSplit[0]
	objKey := strings.Join(destSplit[1:], "/")
	fileName := destSplit[(destSplitLen - 1)]

	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	var contentType string
	if contentType, err = st.getFileContentType(file); err != nil {
		return err
	}

	_, err = st.client.FileUpload(simples3.UploadInput{
		Bucket:      bucket,
		ObjectKey:   objKey,
		FileName:    fileName,
		ContentType: contentType,
		Body:        file,
	})

	return err
}

func (st *S3) getFileContentType(file *os.File) (contentType string, err error) {
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
