package wellez

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/naveego/api/pipeline/publisher"
)

func extractFiles(ctx publisher.Context, tmpDir string, files fileInfos) error {

	ftpFilePwd, valid := ctx.GetStringSetting("ftp_file_pwd")
	if !valid {
		return errors.New("Expected setting for 'ftp_file_pwd' but it was not set or not a valid string")
	}

	for _, file := range files {

		outputDir := filepath.Join(tmpDir, file.LocalDirName)

		errOutput := &bytes.Buffer{}
		c := exec.Command("unzip", "-P", ftpFilePwd, file.FileName, "-d", outputDir)
		c.Dir = tmpDir
		c.Stderr = errOutput
		err := c.Run()

		if err != nil {
			return fmt.Errorf("Extract error: %s", errOutput.String())
		}

		ctx.Logger.Infof("Successfully extracted file '%s'", file.FileName)
	}

	return nil

}
