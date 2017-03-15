package wellez

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dutchcoders/goftp"
	"github.com/naveego/api/pipeline/publisher"
	"github.com/naveego/api/types/pipeline"
)

type Publisher struct {
}

// NewPublisher creates a new Wellcast publisher instance
func NewPublisher() publisher.Publisher {
	return &Publisher{}
}

func (p *Publisher) TestConnection(ctx publisher.Context, connSettings map[string]interface{}) (bool, string, error) {
	return true, "", nil
}

func (p *Publisher) Shapes(ctx publisher.Context) (pipeline.ShapeDefinitions, error) {
	return nil, nil
}

func (p *Publisher) Publish(ctx publisher.Context, shape pipeline.ShapeDefinition, dataTransport publisher.DataTransport) {
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
	defer ftp.Close()

	err = ftp.Login(ftpUser, ftpPwd)
	if err != nil {
		return fileInfos, errors.New("Could not login to FTP server: " + err.Error())
	}

	err = ftp.Cwd("/")
	if err != nil {
		return fileInfos, errors.New("Could not set CWD to '/': " + err.Error())
	}

	files, err := ftp.List("")
	if err != nil {
		return fileInfos, errors.New("Could not list contents of FTP folder: " + err.Error())
	}

	now := time.Now()
	currentYear, currentMonth, currentDay := now.Date()
	lastSunday := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, time.UTC)

	switch now.Weekday() {
	case time.Monday:
		lastSunday = lastSunday.AddDate(0, 0, -1)
	case time.Tuesday:
		lastSunday = lastSunday.AddDate(0, 0, -2)
	case time.Wednesday:
		lastSunday = lastSunday.AddDate(0, 0, -3)
	case time.Thursday:
		lastSunday = lastSunday.AddDate(0, 0, -4)
	case time.Friday:
		lastSunday = lastSunday.AddDate(0, 0, -5)
	case time.Saturday:
		lastSunday = lastSunday.AddDate(0, 0, -6)
	}

	for _, file := range files {
		info, err := parseFileInfo(file)
		if err != nil {
			ctx.Logger.Warnf("Could not parse file info: %v", err)
			continue
		}

		if info.FileName == "." || info.FileName == ".." {
			continue
		}

		if info.ModifiedOn == lastSunday || info.ModifiedOn.After(lastSunday) {

			ctx.Logger.Infof("Retrieving file '%s' from FTP", info.FileName)
			fileInfos = append(fileInfos, info)

			_, err := ftp.Retr(info.FileName, func(r io.Reader) error {

				var file *os.File
				if file, err = os.Create(filepath.Join(tmpDir, info.FileName)); err != nil {
					return err
				}

				if _, err := io.Copy(file, r); err != nil {
					return err
				}

				if err := file.Close(); err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				return fileInfos, err
			}
		}
	}

	sort.Sort(fileInfos)
	return fileInfos, nil

}

type fileInfo struct {
	FileName     string
	LocalDirName string // Used to obfuscate the file names
	ModifiedOn   time.Time
}

type fileInfos []fileInfo

func (f fileInfos) Len() int {
	return len(f)
}

func (f fileInfos) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f fileInfos) Less(i, j int) bool {
	f1 := f[i]
	f2 := f[j]
	return f1.ModifiedOn.Before(f2.ModifiedOn)
}

func parseFileInfo(info string) (fileInfo, error) {
	f := fileInfo{}

	parts := strings.Split(info, ";")

	f.FileName = strings.TrimSpace(parts[len(parts)-1])
	f.LocalDirName, _ = generateRandomString(16)

	modifyStr := parts[0]
	modifyParts := strings.Split(modifyStr, "=")
	dateStr := modifyParts[1]
	yearStr := dateStr[:4]
	monthStr := dateStr[4:6]
	dayStr := dateStr[6:8]
	hourStr := dateStr[8:10]
	minuteStr := dateStr[10:12]
	secStr := dateStr[12:14]

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return f, err
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return f, err
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return f, err
	}

	hour, err := strconv.Atoi(hourStr)
	if err != nil {
		return f, err
	}

	minute, err := strconv.Atoi(minuteStr)
	if err != nil {
		return f, err
	}

	sec, err := strconv.Atoi(secStr)
	if err != nil {
		return f, err
	}

	f.ModifiedOn = time.Date(year, time.Month(month), day, hour, minute, sec, 0, time.UTC)

	return f, nil
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
