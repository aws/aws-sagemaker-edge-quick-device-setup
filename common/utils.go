package common

import (
	"archive/tar"
	"aws-sagemaker-edge-quick-device-setup/aws"
	"aws-sagemaker-edge-quick-device-setup/cli"
	"aws-sagemaker-edge-quick-device-setup/constants"
	"compress/gzip"
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Release struct {
	s3Location    string
	sha1_shasum   string
	sha256_shasum string
	sha512_shasum string
	md5_shasum    string
}

func GetAgentRelease(client *s3.Client, bucketName *string, prefix *string) *Release {
	output := aws.ListBucket(client, bucketName, prefix)
	releases := make(map[int]*Release)
	releaseDates := make([]int, 0)
	for _, value := range output.Contents {
		paths := strings.Split(*value.Key, "/")
		version := strings.Split(paths[1], ".")
		if len(version) != 3 {
			continue
		}
		date, err := strconv.Atoi(version[1])
		if err != nil {
			continue
		}
		release, ok := releases[date]

		if !ok {
			release = &Release{}
			releases[date] = release
			releaseDates = append(releaseDates, date)
		}

		if strings.HasSuffix(paths[2], "tgz") || strings.HasSuffix(paths[2], "zip") {
			release.s3Location = *value.Key
		} else if strings.HasSuffix(paths[2], "shasum") {
			if strings.HasPrefix(paths[2], "sha1") {
				release.sha1_shasum = *value.Key
			} else if strings.HasPrefix(paths[2], "sha256") {
				release.sha256_shasum = *value.Key
			} else if strings.HasPrefix(paths[2], "sha512") {
				release.sha512_shasum = *value.Key
			} else if strings.HasPrefix(paths[2], "md5") {
				release.md5_shasum = *value.Key
			}
		}
	}

	sort.Ints(releaseDates)
	latestReleaseDate := releaseDates[len(releaseDates)-1]
	return releases[latestReleaseDate]
}

func DownloadAgent(client *s3.Client, cliArgs *cli.CliArgs) *string {

	arch := cliArgs.TargetPlatform.Arch

	// map target arch to the s3 bucket
	if cliArgs.TargetPlatform.Arch == constants.AMD64 || cliArgs.TargetPlatform.Arch == constants.X86_64 || cliArgs.TargetPlatform.Arch == constants.X64 {
		arch = constants.X64
	} else if cliArgs.TargetPlatform.Arch == constants.I386 || cliArgs.TargetPlatform.Arch == constants.X86 {
		arch = constants.X86
	} else if cliArgs.TargetPlatform.Arch == constants.ARM64 || cliArgs.TargetPlatform.Arch == constants.ARMV8 {
		arch = constants.ARMV8
	}

	agentBucket := fmt.Sprintf("sagemaker-edge-release-store-us-west-2-%s-%s", cliArgs.TargetPlatform.Os, arch)
	s3Prefix := "Releases/"
	release := GetAgentRelease(client, &agentBucket, &s3Prefix)
	agentFile := aws.DownloadFileFromS3(client, &agentBucket, &release.s3Location)
	if strings.HasSuffix(*agentFile, "gz") {
		untar(agentFile, &cliArgs.AgentDirectory)
	} else if strings.HasSuffix(*agentFile, "zip") {
		unzip(agentFile, &cliArgs.AgentDirectory)
	} else {
		log.Fatal("Unsupported agent format!")
	}

	return agentFile
}

func unzip(src *string, dest *string) {
    r, err := zip.OpenReader(*src)
    if err != nil {
        log.Fatal(err)
    }
    defer r.Close()

    for _, f := range r.File {

        // Store filename/path for returning and using later on
        fpath := filepath.Join(*dest, f.Name)

        // Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
        if !strings.HasPrefix(fpath, filepath.Clean(*dest)+string(os.PathSeparator)) {
        	log.Fatal(fmt.Sprintf("%s: illegal file path", fpath))
        }

		fmt.Println(fpath)

        if f.FileInfo().IsDir() {
            // Make Folder
            os.MkdirAll(fpath, os.ModePerm)
            continue
        }

        // Make File
        if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
            log.Fatal(err)
        }

        outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            log.Fatal(err)
        }

        rc, err := f.Open()
        if err != nil {
            log.Fatal(err)
        }

        _, err = io.Copy(outFile, rc)

        // Close the file without defer to close before next iteration of loop
        outFile.Close()
        rc.Close()

        if err != nil {
            log.Fatal(err)
        }
    }
}

func untar(agentFile *string, dest *string) {
	file, err := os.Open(*agentFile)

	if err != nil {
		log.Fatal("Error ", err)
	}

	defer file.Close()
	var fileReader io.ReadCloser = file
	if strings.HasSuffix(*agentFile, "gz") {
		if fileReader, err = gzip.NewReader(file); err != nil {
			log.Fatal("Error ", err)
		}
		defer fileReader.Close()
	}
	tarBallReader := tar.NewReader(fileReader)
	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Error ", err)
		}

		// get the individual filename and extract to the current directory
		if !strings.Contains(header.Name, "..") {
			filename := filepath.Join(*dest, header.Name)

			switch header.Typeflag {
			case tar.TypeDir:
				// handle directory
				fmt.Println("Creating directory :", filename)
				err = os.MkdirAll(filename, os.FileMode(header.Mode)) // or use 0755 if you prefer

				if err != nil {
					log.Fatal("Error ", err)
				}

			case tar.TypeReg:
				// handle normal file
				fmt.Println("Untarring :", filename)
				writer, err := os.Create(filename)

				if err != nil {
					log.Fatal("Error ", err)
				}

				io.Copy(writer, tarBallReader)

				err = os.Chmod(filename, os.FileMode(header.Mode))

				if err != nil {
					log.Fatal("Error ", err)
				}

				writer.Close()
			default:
				fmt.Printf("Unable to untar type : %c in file %s", header.Typeflag, filename)
			}
		}
	}
}

func DownloadSigningRootCert(client *s3.Client, cliArgs *cli.CliArgs) {
	certBucket := "sagemaker-edge-release-store-us-west-2-linux-x64"
	certKey := fmt.Sprintf("Certificates/%s/%s.pem", cliArgs.Region, cliArgs.Region)
	certPath := filepath.Join(cliArgs.AgentDirectory, "certificates", fmt.Sprintf("%s.pem", cliArgs.Region))
	aws.DownloadFileFromS3ToPath(client, &certBucket, &certKey, &certPath)
	os.Chmod(certPath, 0400)
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
