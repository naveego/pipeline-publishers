package wellez

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	goftp "github.com/jlaffaye/ftp"
	"github.com/naveego/api/pipeline/publisher"
	"github.com/naveego/api/types/pipeline"
)

var fullNameFormat = "Synergy_Resources_20010101-%s.zip"
var nameFormat = "Synergy_Resources_%s.zip"


type Publisher struct {
}

// NewPublisher creates a new Wellcast publisher instance
func NewPublisher() publisher.Publisher {
	return &Publisher{}
}

func (p *Publisher) Shapes(ctx publisher.Context) (map[string]pipeline.Shape, error) {
	return nil, nil
}

func (p *Publisher) Publish(ctx publisher.Context, dataTransport publisher.DataTransport) {
	tmpDir, err := ioutil.TempDir("", "wellezxml_")
	if err != nil {
		ctx.Logger.Warn("Could not create temp directory for file storage: ", err)
		return
	}
	defer func(t string) {
		os.RemoveAll(t)
	}(tmpDir)

	files, err := fetchFiles(ctx, tmpDir)
	if err != nil {
		ctx.Logger.Warn("Could not fetch files from FTP: ", err)
		return
	}

	err = extractFiles(ctx, tmpDir, files)
	if err != nil {
		ctx.Logger.Warn("Could not extract files: ", err)
		return
	}

	for _, file := range files {
		err = processFile(ctx, dataTransport, tmpDir, file)
		if err != nil {
			ctx.Logger.Warnf("Could not process file '%s': %v", file.FileName, err)
		}
	}

}

func fetchFiles(ctx publisher.Context, tmpDir string) (fileInfos, error) {
	fileInfos := fileInfos{}

	ftpAddr, valid := ctx.GetStringSetting("ftp_server")
	if !valid {
		return fileInfos, errors.New("Expected setting for 'ftp_server' but it was not set or not a valid string")
	}

	ftpUser, valid := ctx.GetStringSetting("ftp_user")
	if !valid {
		return fileInfos, errors.New("Expected setting for 'ftp_user' but it was not set or not a valid string")
	}

	ftpPwd, valid := ctx.GetStringSetting("ftp_pwd")
	if !valid {
		return fileInfos, errors.New("Expected setting for 'ftp_pwd' but it was not set or not a valid string")
	}

	ctx.Logger.Infof("Fetching files from FTP: %s", ftpAddr)

	if strings.Contains(ftpAddr, ":") == false {
		ftpAddr = ftpAddr + ":21"
	}

	ftp, err := goftp.Connect(ftpAddr)
	if err != nil {
		return fileInfos, errors.New("Could not connect to FTP Server: " + err.Error())
	}
	defer func() {
		ftp.Logout()
		ftp.Quit()
	}()

	err = ftp.Login(ftpUser, ftpPwd)
	if err != nil {
		return fileInfos, errors.New("Could not login to FTP server: " + err.Error())
	}

	err = ftp.ChangeDir("/XML Exports/Synergy Resources")
	if err != nil {
		return fileInfos, errors.New("Could not set CWD to '/XML Exports/Synergy Resources': " + err.Error())
	}

	var curpath string
	if curpath, err = ftp.CurrentDir(); err != nil {
		return fileInfos, errors.New("Could not get PWD: " + err.Error())
	}

	ctx.Logger.Infof("Current path: %s", curpath)

	now := time.Now()
	currentYear, currentMonth, currentDay := now.Date()
	lastSaturday := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, time.UTC)

	ctx.Logger.Infof("Current Dates: Year=%v, Month=%v, Day=%v", currentYear, currentMonth, currentDay)

	daysToProcess := 0

	switch now.Weekday() {
	case time.Sunday:
		lastSaturday = lastSaturday.AddDate(0, 0, -1)
		daysToProcess = 1
	case time.Monday:
		lastSaturday = lastSaturday.AddDate(0, 0, -2)
		daysToProcess = 2
	case time.Tuesday:
		lastSaturday = lastSaturday.AddDate(0, 0, -3)
		daysToProcess = 3
	case time.Wednesday:
		lastSaturday = lastSaturday.AddDate(0, 0, -4)
		daysToProcess = 4
	case time.Thursday:
		lastSaturday = lastSaturday.AddDate(0, 0, -5)
		daysToProcess = 5
	case time.Friday:
		lastSaturday = lastSaturday.AddDate(0, 0, -6)
		daysToProcess = 6
	}

	var files []string

	if daysToProcess == 0 {
		files = append(files, fmt.Sprintf(fullNameFormat, now.Format("20060102")))
	} else {
		for i := 0; i < daysToProcess; i++ {
			day := lastSaturday.AddDate(0, 0, i+1)
			files = append(files, fmt.Sprintf(nameFormat, day.Format("20060102")))
		}
	}

	for _, file := range files {
		fi, err := downloadFile(tmpDir, ctx, ftp, file)
		fileInfos = append(fileInfos, fi)

		if err != nil {
			return fileInfos, err
		}
	}

	return fileInfos, nil
}

func downloadFile(tmpDir string, ctx publisher.Context, ftp *goftp.ServerConn, file string) (fileInfo, error) {
	ctx.Logger.Infof("Retrieving file '%s' from FTP", file)

	info := fileInfo{
		FileName:     file,
		LocalDirName: file[:(len(file) - len(filepath.Ext(file)))],
	}

	rr, err := ftp.Retr(file)
	if err != nil {
		return info, err
	}
	defer rr.Close()

	os.Mkdir(filepath.Join(tmpDir, info.LocalDirName), 0700)
	ctx.Logger.Infof("Using Temp Dir: %s/%s", tmpDir, info.LocalDirName)
	var outFile *os.File
	if outFile, err = os.Create(filepath.Join(tmpDir, info.LocalDirName, file)); err != nil {
		return info, err
	}

	if _, err := io.Copy(outFile, rr); err != nil {
		return info, err
	}

	if err := outFile.Close(); err != nil {
		return info, err
	}

	return info, nil
}

type fileInfo struct {
	FileName     string
	LocalDirName string // Used to obfuscate the file names
}

type fileInfos []fileInfo
