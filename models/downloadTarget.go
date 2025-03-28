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
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	TargetDmm TargetType = "dmm"
	TargetMgs TargetType = "mgs"
	dmmCats              = []string{"digital/video", "digital/amateur", "mono/movie"} //	"digital/e-book"

	dmmSeps = []string{"00", "", "0"}
	mgsCats = []string{"images/prestige", "images/jackson"}
	mgsSeps = []string{""}

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

	//category string
	sep string
}

func (d *DownloadTarget) HadFilesDownloaded() bool {
	if d.localFiles == nil {
		return false
	}
	return len(d.localFiles) > 0
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
	d.Group = strings.TrimPrefix(strings.TrimSuffix(d.Group, " "), " ")

	d.Group = strings.ToLower(d.Group)

	// remove unwanted/repeated chars
	d.sanitizeName()
	// shorten the name to be less than 255 bytes
	d.shortenName()
}

func (d *DownloadTarget) sanitizeName() {

	// Replace all occurrences of \t with a single space
	for _, c := range []string{"\t", "/", "／"} {
		d.Name = strings.ReplaceAll(d.Name, c, " ")
	}

	// Regex to remove content inside square brackets, including the brackets themselves
	re := regexp.MustCompile(`\[[^\]]+\]`)
	d.Name = re.ReplaceAllString(d.Name, "")

	// Replace multiple spaces with a single space
	d.Name = strings.Join(strings.Fields(d.Name), " ")
	// Remove " " prefix and Suffix
	d.Name = strings.TrimPrefix(strings.TrimSuffix(d.Name, " "), " ")
}

func (d *DownloadTarget) shortenName() {

	for {
		words := strings.Fields(d.Name) // Split into words
		l := len(words)

		if l <= 2 || words[l-1] != words[l-2] {
			break
		}

		// remove the last section
		d.Name = strings.Join(words[:l-1], " ")
	}

	//make sure the make not longer than 120?, but keep last element
	for len(d.Name) > 200 {
		words := strings.Fields(d.Name) // Split into words
		l := len(words)

		if l <= 0 {
			return
		}

		switch l {
		case 1:
			runes := []rune(d.Name) // Convert string to rune slice (handles multi-byte characters)
			newCharacters := int(float32(len(runes)) * 0.9)
			d.Name = string(runes[:newCharacters]) // Keep only the first 100 characters
			return
		case 2:
			runes := []rune(words[0]) // Convert string to rune slice (handles multi-byte characters)
			newCharacters := int(float32(len(runes)) * 0.9)
			words[0] = string(runes[:newCharacters]) // Keep only the first 100 characters
		default:
			// l >= 3
			words = append(words[:l-2], words[l-1]) // Remove second-to-last but keep last word
		}

		d.Name = strings.Join(words, " ") // Reconstruct the name
	}
}

func (d *DownloadTarget) BuildTitlePath(cat, sep string) *url.URL {
	var copiedURL url.URL
	switch d.Source {
	case TargetDmm:
		copiedURL = dmmUrl
		copiedURL.Path = path.Join(cat, d.Group+sep+d.Number, d.Group+sep+d.Number+"pl.jpg")
	case TargetMgs:
		copiedURL = mgsUrl
		copiedURL.Path = path.Join(cat, d.Group, d.Number, "pb_e_"+d.Group+"-"+d.Number+".jpg")
	}
	return &copiedURL
}

//func (d *DownloadTarget) BuildDmmTitlePath(cat, sep string) string {
//	return path.Join(cat, d.Group+sep+d.Number, d.Group+sep+d.Number+"pl.jpg")
//}

func (d *DownloadTarget) BuildSubPath(cat string, sep string, cnt int, hd string) *url.URL {
	var copiedURL url.URL

	switch d.Source {
	case TargetDmm:
		copiedURL = dmmUrl
		copiedURL.Path = path.Join(
			cat,
			fmt.Sprint(d.Group, sep, d.Number),
			fmt.Sprint(d.Group, sep, d.Number+hd, "-", cnt, ".jpg"),
		)
	case TargetMgs:
		copiedURL = mgsUrl
		copiedURL.Path = path.Join(
			cat,
			d.Group,
			d.Number,
			fmt.Sprint("cap_e_", cnt, "_", d.Group, "-", d.Number, ".jpg"),
		)
	default:
	}

	return &copiedURL
}

//func (d *DownloadTarget) BuildMgsTitlePath() string {
//	return path.Join("images/prestige", d.Group, d.Number, "pb_e_"+d.Group+"-"+d.Number+".jpg")
//}

func (d *DownloadTarget) BuildFolderName() (withoutName, withName string) {

	tmpGroup := d.Group
	// Extract the part of Group after the last '_'
	groupParts := strings.Split(tmpGroup, "_")
	tmpGroup = groupParts[len(groupParts)-1]

	groupParts2 := strings.Split(tmpGroup, "-")
	tmpGroup = groupParts2[len(groupParts2)-1]

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
	var cats, seps []string
	switch d.Source {
	case TargetMgs:
		cats = mgsCats
		seps = mgsSeps
	case TargetDmm:
		cats = dmmCats
		seps = dmmSeps
	default:
	}

	//download main pic
	for _, cat := range cats {
		for _, sep := range seps {

			// download main pic
			u := d.BuildTitlePath(cat, sep)

			// Get the file name from the URL path
			fileName := path.Base(u.Path)

			localFilepath := path.Join(d.localPath, fileName)
			if err = d.DownloadRemoteFile(*u, localFilepath); err == nil {
				//d.category = cat
				d.sep = sep
				return
			}

			log.Info("Fail to download ", *u)
		}
	}

	return http.ErrMissingFile
}

func (d *DownloadTarget) DownloadSub() (err error) {
	var cats []string
	switch d.Source {
	case TargetMgs:
		cats = mgsCats
	case TargetDmm:
		cats = dmmCats
	default:
	}

	var correctCat, correctHd string
	cnt := 1
outerLoop: // Label for the outermost loop
	for _, hd := range []string{"jp", ""} {
		for _, cat := range cats {
			//for cnt := 0; cnt <= 1; cnt++ {
			u := d.BuildSubPath(cat, d.sep, cnt, hd)
			log.Infof("Trying %v", u.Path)

			// Get the file name from the URL path
			fileName := path.Base(u.Path)
			if err = d.DownloadRemoteFile(*u, path.Join(d.localPath, fileName)); err == nil {
				correctCat = cat
				correctHd = hd
				log.Infof("Catched:%#v", u.String())
				break outerLoop // Exit all loops once condition is met
			}
			//}
		}
	}

	if correctCat == "" {
		log.Infof("Fail to find proper category")
		return http.ErrBodyNotAllowed
	}

	log.Infof("Download params: cat:%v hd:%v", correctCat, correctHd)
	for cnt := 2; cnt <= 50; cnt++ {
		u := d.BuildSubPath(correctCat, d.sep, cnt, correctHd)

		// Get the file name from the URL path
		fileName := path.Base(u.Path)

		if err = d.DownloadRemoteFile(*u, path.Join(d.localPath, fileName)); err != nil {
			//we don't care if 2~n can be retrieved or not as long as the first image does exist
			return nil
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
	for _, fil := range d.localFiles {
		// Extract the file name
		fileName := filepath.Base(fil)
		// Define the destination path
		destPath := filepath.Join(dirPath, fileName)

		//log.Infof("Moving file %v to %v", fil, destPath)
		// Move the file
		if err = os.Rename(fil, destPath); err != nil {
			log.Infof("failed to move file %v to %v:%v", fil, destPath, err)
		}
	}

	return
}
