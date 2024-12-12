package models

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

var (
	TargetDmm TargetType = "dmm"
	TargetMgs TargetType = "mgs"

	dmmUrl = url.URL{
		Scheme: "https",
		Host:   "pics.dmm.co.jp",
	}

	mgsUrl = url.URL{
		Scheme: "https",
		Host:   "image.mgstage.com",
	}
)

type TargetType string

type DownloadTarget struct {
	Source TargetType `json:"source"  form:"source"`
	Group  string     `json:"group"  form:"group"`
	Number string     `json:"number"  form:"number"`
	Name   string     `json:"name"  form:"name"`

	localPath  string
	localFiles []string

	category string
	sep      string
}

func (d *DownloadTarget) SetLocalPathBase(basePath string) {
	if len(basePath) <= 0 {
		return
	}

	d.localPath = basePath
}

func (d *DownloadTarget) Sanitize() {
	// Replace all occurrences of \t with a single space
	d.Group = strings.ReplaceAll(d.Group, "\t", " ")
	// Replace multiple spaces with a single space
	d.Group = strings.Join(strings.Fields(d.Group), " ")
	// Remove " " prefix and Suffix
	d.Group = strings.TrimPrefix(d.Group, " ")
	d.Group = strings.TrimSuffix(d.Group, " ")

	d.Group = strings.ToLower(d.Group)

	// Replace all occurrences of \t with a single space
	d.Name = strings.ReplaceAll(d.Name, "\t", " ")
	// Replace multiple spaces with a single space
	d.Name = strings.Join(strings.Fields(d.Name), " ")
	// Remove " " prefix and Suffix
	d.Name = strings.TrimPrefix(d.Name, " ")
	d.Name = strings.TrimSuffix(d.Name, " ")
}

func (d *DownloadTarget) BuildTitlePath(cat, sep string) *url.URL {
	var copiedURL url.URL
	switch d.Source {
	case TargetDmm:
		copiedURL = dmmUrl
		copiedURL.Path = d.BuildDmmTitlePath(cat, sep)
	case TargetMgs:
		copiedURL = mgsUrl
		copiedURL.Path = d.BuildMgsTitlePath()
	}
	return &copiedURL
}

func (d *DownloadTarget) BuildDmmTitlePath(cat, sep string) string {
	return path.Join("digital", cat, d.Group+sep+d.Number, d.Group+sep+d.Number+"pl.jpg")
}

func (d *DownloadTarget) BuildSubPath(cat string, sep string, cnt int, hd string) *url.URL {
	var copiedURL url.URL

	switch d.Source {
	case TargetDmm:
		copiedURL = dmmUrl
		copiedURL.Path = path.Join(
			"digital",
			cat,
			fmt.Sprint(d.Group, sep, d.Number),
			fmt.Sprint(d.Group, sep, d.Number+hd, "-", cnt, ".jpg"),
		)
	case TargetMgs:
		copiedURL = mgsUrl
		copiedURL.Path = path.Join(
			"images/prestige",
			d.Group,
			d.Number,
			fmt.Sprint("cap_e_", cnt, "_", d.Group, "-", d.Number, ".jpg"),
		)
	default:
	}

	return &copiedURL
}

func (d *DownloadTarget) BuildMgsTitlePath() string {
	return path.Join("images/prestige", d.Group, d.Number, "pb_e_"+d.Group+"-"+d.Number+".jpg")
}

func (d *DownloadTarget) BuildFolderName() (withoutName, withName string) {

	tmpGroup := d.Group
	// Extract the part of Group after the last '_'
	groupParts := strings.Split(tmpGroup, "_")
	tmpGroup = groupParts[len(groupParts)-1]

	groupParts2 := strings.Split(tmpGroup, "-")
	tmpGroup = groupParts[len(groupParts2)-1]

	// Remove numeric prefix if tmpGroup starts with digits
	tmpGroup = strings.TrimLeftFunc(tmpGroup, unicode.IsDigit)
	tmpGroup = strings.ToUpper(tmpGroup)

	//first part: "/ref/[ABC-123]"
	withoutName = filepath.Join(d.localPath, fmt.Sprintf("[%s-%s]", tmpGroup, d.Number))
	// Construct the destination folder path: /ref/[ABC-123]other
	withName = withoutName + d.Name

	return
}

func (d *DownloadTarget) DownloadRemoteFile(remoteFileUrl url.URL, localFilepath string) (err error) {

	// Validate by Head
	if err = validateContentLength(remoteFileUrl); err != nil {
		return
	}

	// GET the data
	var resp *http.Response
	resp, err = http.Get(remoteFileUrl.String())
	if err != nil {
		log.Infof("failed to download %v: %s", remoteFileUrl, err.Error())
		return
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	// Save the response body to the local file
	if err = saveToFile(resp.Body, localFilepath); err != nil {
		return
	}

	d.localFiles = append(d.localFiles, localFilepath)

	return
}

// Helper to validate Content-Length
func validateContentLength(remoteFileUrl url.URL) (err error) {

	var headResp *http.Response
	// Make a HEAD request
	if headResp, err = http.Head(remoteFileUrl.String()); err != nil {
		return fmt.Errorf("failed to perform HEAD request: %w", err)
	}
	defer headResp.Body.Close()

	contentLengthStr := headResp.Header.Get("Content-Length")
	if contentLengthStr == "" {
		return fmt.Errorf("Content-Length header is missing for %v", remoteFileUrl)
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		return fmt.Errorf("error parsing Content-Length for %v: %v", remoteFileUrl, err)
	}
	if contentLength <= 2732 {
		return fmt.Errorf("content length too small for %v: %d bytes", remoteFileUrl, contentLength)
	}

	//log.Infof("File size for %v: %d bytes", remoteFileUrl, contentLength)
	return
}

// Helper to save HTTP response body to a file
func saveToFile(body io.Reader, filepath string) (err error) {
	var out *os.File
	if out, err = os.Create(filepath); err != nil {
		return fmt.Errorf("failed to create file %v: %w", filepath, err)
	}
	defer out.Close()

	if _, err = io.Copy(out, body); err != nil {
		return fmt.Errorf("failed to write to file %v: %w", filepath, err)
	}

	return
}

func (d *DownloadTarget) TryDownloadMain() (err error) {
	switch d.Source {
	case TargetMgs:
		return d.tryDownloadMgsMain()
	case TargetDmm:
		return d.tryDownloadDmmMain()
	default:
	}

	return http.ErrMissingFile
}

func (d *DownloadTarget) tryDownloadDmmMain() (err error) {

	//download main pic
	for _, cat := range []string{"video", "amateur"} {
		for _, sep := range []string{"00", "", "0"} {

			// download main pic
			u := d.BuildTitlePath(cat, sep)

			// Get the file name from the URL path
			fileName := path.Base(u.Path)

			localFilepath := path.Join(d.localPath, fileName)
			if err = d.DownloadRemoteFile(*u, localFilepath); err == nil {
				d.category = cat
				d.sep = sep
				return
			}
		}
	}

	return http.ErrMissingFile
}

func (d *DownloadTarget) tryDownloadMgsMain() (err error) {

	//download main pic
	for _, cat := range []string{"video", "amateur"} {
		for _, sep := range []string{"00", "", "0"} {

			// download main pic
			u := d.BuildTitlePath(cat, sep)

			// Get the file name from the URL path
			fileName := path.Base(u.Path)

			localFilepath := path.Join(d.localPath, fileName)
			if err = d.DownloadRemoteFile(*u, localFilepath); err == nil {
				//d.localFiles = append(d.localFiles, localFilepath)
				d.category = cat
				d.sep = sep
				return
			}
		}
	}

	return http.ErrMissingFile
}

func (d *DownloadTarget) DownloadSub() (err error) {
	for _, hd := range []string{"jp", ""} {
		for cnt := 1; cnt <= 30; cnt++ {
			// download main pic
			u := d.BuildSubPath(d.category, d.sep, cnt, hd)

			// Get the file name from the URL path
			fileName := path.Base(u.Path)

			if err = d.DownloadRemoteFile(*u, path.Join(d.localPath, fileName)); err != nil {
				if cnt > 1 {
					return nil
				}

				if hd == "" {
					return nil
				}
				break
			}
		}
	}
	return
}

func (d *DownloadTarget) MoveLocalFilesUnderFolder() (err error) {

	// Construct the destination folder path
	destinationFolderWithoutName, destinationFolderWithName := d.BuildFolderName()

	// Ensure the destination folder exists
	var dirPath string
	for _, dirPath = range []string{destinationFolderWithName, destinationFolderWithoutName} {
		if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
			log.Info("Fail to MkdirAll", dirPath)
			continue
		}

		break
	}
	if err != nil {
		return fmt.Errorf("failed to create destination folder: %w", err)
	}

	// Iterate over the local files and move them
	for _, file := range d.localFiles {
		// Extract the file name
		fileName := filepath.Base(file)
		// Define the destination path
		destPath := filepath.Join(dirPath, fileName)
		// Move the file
		if err = os.Rename(file, destPath); err != nil {
			return fmt.Errorf("failed to move file %s to %s: %w", file, destPath, err)
		}
	}

	return
}
