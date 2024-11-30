package models

import (
	"fmt"
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
	Source string `json:"source"  form:"source"`
	Group  string `json:"group"  form:"group"`
	Number string `json:"number"  form:"number"`
	Name   string `json:"name"  form:"name"`

	localFiles []string

	category string
	sep      string
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
	case string(TargetDmm):
		copiedURL = dmmUrl
		copiedURL.Path = d.BuildDmmTitlePath(cat, sep)
	case string(TargetMgs):
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
	case string(TargetDmm):
		copiedURL = dmmUrl
		copiedURL.Path = path.Join(
			"digital",
			cat,
			fmt.Sprint(d.Group, sep, d.Number),
			fmt.Sprint(d.Group, sep, d.Number+hd, "-", cnt, ".jpg"),
		)
	case string(TargetMgs):
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

func (d *DownloadTarget) BuildFolderName(basePath string) string {

	tmpGroup := d.Group
	// Extract the part of Group after the last '_'
	groupParts := strings.Split(tmpGroup, "_")
	tmpGroup = groupParts[len(groupParts)-1]

	groupParts2 := strings.Split(tmpGroup, "-")
	tmpGroup = groupParts[len(groupParts2)-1]

	// Remove numeric prefix if tmpGroup starts with digits
	tmpGroup = strings.TrimLeftFunc(tmpGroup, unicode.IsDigit)

	// Construct the destination folder path
	return filepath.Join(basePath, fmt.Sprintf("[%s-%s]%s", tmpGroup, d.Number, d.Name))
}

func (d *DownloadTarget) DownloadRemoteFile(remoteFileUrl url.URL, localFilepath string) (err error) {

	var headResp *http.Response
	// Make a HEAD request
	headResp, err = http.Head(remoteFileUrl.String())
	if err != nil {
		return
	}
	defer headResp.Body.Close()

	// Check Content-Length header
	contentLengthStr := headResp.Header.Get("Content-Length")
	if contentLengthStr == "" {
		fmt.Println("Content-Length header is missing")
	} else {
		var contentLength int
		contentLength, err = strconv.Atoi(contentLengthStr)
		if err != nil {
			fmt.Println("Error parsing Content-Length:", err)
			return
		} else {
			fmt.Printf("%v File size: %d bytes\n", remoteFileUrl, contentLength)
			if contentLength <= 2732 {
				err = http.ErrContentLength
			}
			return
		}
	}

	// Get the data
	var resp *http.Response
	resp, err = http.Get(remoteFileUrl.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	// Create the file
	var out *os.File
	out, err = os.Create(localFilepath)
	if err != nil {
		return
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	if err == nil {
		d.localFiles = append(d.localFiles, localFilepath)
	}

	return
}

func (d *DownloadTarget) TryDownloadDmmMain() (err error) {

	//download main pic
	for _, cat := range []string{"video", "amateur"} {
		for _, sep := range []string{"00", ""} {

			// download main pic
			u := d.BuildTitlePath(cat, sep)

			// Get the file name from the URL path
			fileName := path.Base(u.Path)

			localFilepath := path.Join("/ref", fileName)
			if err = d.DownloadRemoteFile(*u, localFilepath); err == nil {
				d.localFiles = append(d.localFiles, localFilepath)
				d.category = cat
				d.sep = sep
				return
			}
		}
	}

	return http.ErrMissingFile
}

//func (d *DownloadTarget) TryDownloadMain() (err error) {
//
//	//download main pic; determine from dmm or mgs
//
//	switch d.Source {
//	case string(TargetDmm):
//		for _, sep := range []string{"00", ""} {
//
//			// download main pic
//			u := d.BuildTitlePath(sep)
//
//			// Get the file name from the URL path
//			fileName := path.Base(u.Path)
//
//			localFilepath := path.Join("/ref", fileName)
//			if err = d.DownloadRemoteFile(*u, localFilepath); err == nil {
//				d.localFiles = append(d.localFiles, localFilepath)
//				return
//			}
//		}
//	case string(TargetMgs):
//		// download main pic
//		u := d.BuildTitlePath("")
//
//		// Get the file name from the URL path
//		fileName := path.Base(u.Path)
//		localFilepath := path.Join("/ref", fileName)
//		if err = d.DownloadRemoteFile(*u, localFilepath); err == nil {
//			d.localFiles = append(d.localFiles, localFilepath)
//			return
//		}
//	}
//	return
//}

func (d *DownloadTarget) DownloadSub(localPath string) (err error) {
	for _, hd := range []string{"jp", ""} {
		for cnt := 1; cnt <= 30; cnt++ {
			// download main pic
			u := d.BuildSubPath(d.category, d.sep, cnt, hd)

			// Get the file name from the URL path
			fileName := path.Base(u.Path)

			if err = d.DownloadRemoteFile(*u, path.Join(localPath, fileName)); err != nil {
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

func (d *DownloadTarget) MoveLocalFilesUnderFolder(basePath string) (err error) {

	// Construct the destination folder path
	destinationFolder := d.BuildFolderName(basePath)

	// Ensure the destination folder exists
	err = os.MkdirAll(destinationFolder, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create destination folder: %w", err)
	}

	// Iterate over the local files and move them
	for _, file := range d.localFiles {
		// Extract the file name
		fileName := filepath.Base(file)
		// Define the destination path
		destPath := filepath.Join(destinationFolder, fileName)

		// Move the file
		err := os.Rename(file, destPath)
		if err != nil {
			return fmt.Errorf("failed to move file %s to %s: %w", file, destPath, err)
		}
	}

	return
}
