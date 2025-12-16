package magick

import (
	"io"
	"os"
	"os/exec"
)

type Magick struct {
	convertPath string
	tmpDir      string
}

func New() (im *Magick, err error) {
	im = new(Magick)

	if im.convertPath, err = exec.LookPath("convert"); err != nil {
		return nil, err
	}

	return im, nil
}

func (im *Magick) Startup() (err error) {
	im.tmpDir, err = os.MkdirTemp("", "hyperuplink-*")
	return err
}

func (im *Magick) Shutdown() (err error) {
	if im.tmpDir != "" {
		err = os.RemoveAll(im.tmpDir)
	}
	return err
}

func (im *Magick) Convert(src io.ReadSeekCloser) (dest io.ReadSeekCloser, destName string, err error) {
	var tmpSrcFile, tmpDestFile *os.File
	var tmpSrcName, tmpDestName string

	if tmpSrcFile, err = os.CreateTemp(im.tmpDir, "convert.*.bin"); err != nil {
		return nil, "", err
	}
	tmpSrcName = tmpSrcFile.Name()

	if _, err = io.Copy(tmpSrcFile, src); err != nil {
		tmpSrcFile.Close()
		if tmpSrcName != "" {
			os.Remove(tmpSrcName)
		}
		return nil, "", err
	}

	tmpSrcFile.Close()

	if tmpDestFile, err = os.CreateTemp(im.tmpDir, "convert.*.webp"); err != nil {
		if tmpSrcName != "" {
			os.Remove(tmpSrcName)
		}
		return nil, "", err
	}
	tmpDestName = tmpDestFile.Name()

	tmpDestFile.Close()

	cmd := exec.Command(im.convertPath, tmpSrcName, tmpDestName)
	err = cmd.Run()
	if tmpSrcName != "" {
		os.Remove(tmpSrcName)
	}
	if err != nil {
		if tmpDestName != "" {
			os.Remove(tmpDestName)
		}
		return nil, "", err
	}

	if tmpDestFile, err = os.Open(tmpDestName); err != nil {
		if tmpDestName != "" {
			os.Remove(tmpDestName)
		}
		return nil, "", err
	}

	return tmpDestFile, tmpDestFile.Name(), nil
}
