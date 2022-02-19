package io

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"svm/parsers"
	"time"
)

type Resource struct {
	Home    string
	Url     string
	Version parsers.Version
}

func createProgressBar(fileSize int, progress *mpb.Progress) *mpb.Bar {
	bar := progress.AddBar(
		int64(fileSize),
		// Modification before progress bar
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"), // Number of Downloads
			decor.Percentage(decor.WCSyncSpace),     // Progress percentage
		),
		// Modification after progress bar
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)
	return bar
}

func DownloadFile(resource Resource) (err error) {
	// Create the file
	tempFilePath := filepath.Join(resource.Home, ".svm", resource.Version.FullVersion+".tmp")
	tarPath := filepath.Join(resource.Home, ".svm", resource.Version.FullVersion+".tgz")
	_ = os.Remove(tempFilePath)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := tempFile.Close(); closeErr != nil {
			log.Fatal(closeErr)
		}
	}()

	// Get the data
	resp, err := http.Get(resource.Url)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Fatal(closeErr)
		}
	}()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Check filesize
	fileSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	progress := mpb.New(mpb.WithWidth(64), mpb.WithRefreshRate(180*time.Millisecond))
	bar := createProgressBar(fileSize, progress)
	reader := bar.ProxyReader(resp.Body)
	defer func() {
		if closeErr := reader.Close(); closeErr != nil {
			log.Fatal(closeErr)
		}
	}()

	// Write the body to temp file
	_, err = io.Copy(tempFile, reader)
	if err != nil {
		return err
	}

	// Rename temp file to final file
	if err = os.Rename(tempFilePath, tarPath); err != nil {
		return err
	}

	progress.Wait()
	return nil
}

// UnzipTar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func UnzipTar(resource Resource) error {
	svmPath := filepath.Join(resource.Home, ".svm")
	tarPath := filepath.Join(resource.Home, ".svm", resource.Version.FullVersion+".tgz")
	tarFile, err := os.Open(tarPath)
	gzr, err := gzip.NewReader(tarFile)
	if err != nil {
		return err
	}
	defer func(gzr *gzip.Reader) {
		if closeErr := gzr.Close(); closeErr != nil {
			log.Fatal(closeErr)
		}
		if removeErr := os.Remove(tarPath); removeErr != nil {
			log.Fatal(removeErr)
		}
	}(gzr)

	tr := tar.NewReader(gzr)

	for {
		header, getNextErr := tr.Next()

		switch {

		// if no more files are found return
		case getNextErr == io.EOF:
			return nil

		// return any other error
		case getNextErr != nil:
			return getNextErr

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(svmPath, header.Name)

		// the following switch could also be done using fi.Mode()
		// Not sure if there is a benefit of using one or the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if dir and it doesn't exist create it
		case tar.TypeDir:
			if _, typeDirErr := os.Stat(target); typeDirErr != nil {
				if mkDirErr := os.MkdirAll(target, 0755); mkDirErr != nil {
					return mkDirErr
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, openFileErr := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if openFileErr != nil {
				return openFileErr
			}

			// copy over contents
			if _, copyErr := io.Copy(f, tr); copyErr != nil {
				return copyErr
			}

			// manually close here after each file operation; deferring would cause each file close
			// to wait until all operations have completed.
			if err = f.Close(); err != nil {
				return err
			}
		}
	}
}

func RenameUnzipped(resource Resource) error {
	sparkUnzippedPath := filepath.Join(resource.Home, ".svm", resource.Version.DownloadVersion)
	sparkFinalPath := filepath.Join(resource.Home, ".svm", resource.Version.FullVersion)
	err := os.Rename(sparkUnzippedPath, sparkFinalPath)
	if err != nil {
		return err
	}
	fmt.Printf("Installed %s successfully\n", resource.Version.FullVersion)
	return nil
}

func CreateSVMDirectory(dirname string) error {
	svmPath := filepath.Join(dirname, ".svm")
	_, err := os.Stat(svmPath)
	if os.IsNotExist(err) {
		if _, createFolderErr := os.Create(svmPath); createFolderErr != nil {
			return createFolderErr
		}
	}
	return nil
}

//func CheckIfExists(dirname string) (bool, error) {
//
//}
