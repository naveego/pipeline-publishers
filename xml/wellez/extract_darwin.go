package wellez

import (
	"errors"
	"os/exec"

	"github.com/naveego/api/pipeline/publisher"
)

func extractFiles(ctx publisher.Context, tmpDir string, files fileInfos) error {

	ftpFilePwd, valid := ctx.GetStringSetting("ftp_file_pwd")
	if !valid {
		return errors.New("Expected setting for 'ftp_file_pwd' but it was not set or not a valid string")
	}

	for _, file := range files {

		c := exec.Command("unzip", "-P", ftpFilePwd, file.FileName)
		c.Dir = tmpDir
		err := c.Run()

		if err != nil {
			return err
		}

		ctx.Logger.Infof("Successfully extracted file '%s'", file.FileName)
	}

	return nil

}
