package local

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/services/config"
)

type Local struct {
	cfg config.Storage
}

func New(cfg config.Storage) (st *Local, err error) {
	st = new(Local)

	st.cfg = cfg

	return st, nil
}

func (st *Local) Startup() (err error) {
	if st.cfg.Local.Path == "" {
		return errs.ErrFilePathInvalid
	}

	return os.MkdirAll(st.cfg.Local.Path, 0o755)
}

func (st *Local) Shutdown() (err error) {
	return nil
}

func (st *Local) resolvePath(dest string) (fullPath string, err error) {
	base := filepath.Clean(st.cfg.Local.Path)
	fullPath = filepath.Clean(filepath.Join(base, dest))

	if fullPath != base &&
		!strings.HasPrefix(fullPath, base+string(os.PathSeparator)) {
		return "", errs.ErrFilePathInvalid
	}

	return fullPath, nil
}

func (st *Local) StoreFileName(src string, dest string) (err error) {
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

func (st *Local) StoreFile(src io.ReadSeeker, dest string) (err error) {
	if src == nil || dest == "" {
		return errs.ErrFilePathInvalid
	}

	var fullPath string
	if fullPath, err = st.resolvePath(dest); err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return err
	}

	if _, err = src.Seek(0, io.SeekStart); err != nil {
		return err
	}

	var dstFile *os.File
	if dstFile, err = os.Create(fullPath); err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, src)
	return err
}

func (st *Local) GetFileDownloadURL(dest string) (
	dlurl string,
	abs bool,
	err error,
) {
	if dest == "" {
		return dlurl, false, errs.ErrFilePathInvalid
	}

	return fmt.Sprintf("%s/%s", st.cfg.Local.PublicURI, dest), false, nil
}

func (st *Local) GetFile(dest string) (reader io.ReadCloser, err error) {
	if dest == "" {
		return nil, errs.ErrFilePathInvalid
	}

	var fullPath string
	if fullPath, err = st.resolvePath(dest); err != nil {
		return nil, err
	}

	return os.Open(fullPath)
}

func (st *Local) DeleteFile(dest string) (err error) {
	if dest == "" {
		return errs.ErrFilePathInvalid
	}

	var fullPath string
	if fullPath, err = st.resolvePath(dest); err != nil {
		return err
	}

	return os.Remove(fullPath)
}
