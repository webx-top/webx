package com

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

// ==============================
// string
// ==============================
func Int64(i interface{}) int64 {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		log.Printf("string[%s] covert int64 fail. %s", in, err)
		return 0
	}
	return out
}

func Int(i interface{}) int {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.Atoi(in)
	if err != nil {
		log.Printf("string[%s] covert int fail. %s", in, err)
		return 0
	}
	return out
}

func Int32(i interface{}) int32 {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseInt(in, 10, 32)
	if err != nil {
		log.Printf("string[%s] covert int32 fail. %s", in, err)
		return 0
	}
	return int32(out)
}

func Float32(i interface{}) float32 {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseFloat(in, 32)
	if err != nil {
		log.Printf("string[%s] covert float32 fail. %s", in, err)
		return 0
	}
	return float32(out)
}

func Float64(i interface{}) float64 {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseFloat(in, 64)
	if err != nil {
		log.Printf("string[%s] covert float64 fail. %s", in, err)
		return 0
	}
	return out
}

func Str(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

// ==============================
// slice
// ==============================
type reducetype func(interface{}) interface{}
type filtertype func(interface{}) bool

func InSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func InSliceIface(v interface{}, sl []interface{}) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func SliceRandList(min, max int) []int {
	if max < min {
		min, max = max, min
	}
	length := max - min + 1
	t0 := time.Now()
	rand.Seed(int64(t0.Nanosecond()))
	list := rand.Perm(length)
	for index, _ := range list {
		list[index] += min
	}
	return list
}

func SliceMerge(slice1, slice2 []interface{}) (c []interface{}) {
	c = append(slice1, slice2...)
	return
}

func SliceReduce(slice []interface{}, a reducetype) (dslice []interface{}) {
	for _, v := range slice {
		dslice = append(dslice, a(v))
	}
	return
}

func SliceRand(a []interface{}) (b interface{}) {
	randnum := rand.Intn(len(a))
	b = a[randnum]
	return
}

func SliceSum(intslice []int64) (sum int64) {
	for _, v := range intslice {
		sum += v
	}
	return
}

func SliceFilter(slice []interface{}, a filtertype) (ftslice []interface{}) {
	for _, v := range slice {
		if a(v) {
			ftslice = append(ftslice, v)
		}
	}
	return
}

func SliceDiff(slice1, slice2 []interface{}) (diffslice []interface{}) {
	for _, v := range slice1 {
		if !InSliceIface(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

func SliceIntersect(slice1, slice2 []interface{}) (diffslice []interface{}) {
	for _, v := range slice1 {
		if !InSliceIface(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

func SliceChunk(slice []interface{}, size int) (chunkslice [][]interface{}) {
	if size >= len(slice) {
		chunkslice = append(chunkslice, slice)
		return
	}
	end := size
	for i := 0; i <= (len(slice) - size); i += size {
		chunkslice = append(chunkslice, slice[i:end])
		end += size
	}
	return
}

func SliceRange(start, end, step int64) (intslice []int64) {
	for i := start; i <= end; i += step {
		intslice = append(intslice, i)
	}
	return
}

func SlicePad(slice []interface{}, size int, val interface{}) []interface{} {
	if size <= len(slice) {
		return slice
	}
	for i := 0; i < (size - len(slice)); i++ {
		slice = append(slice, val)
	}
	return slice
}

func SliceUnique(slice []interface{}) (uniqueslice []interface{}) {
	for _, v := range slice {
		if !InSliceIface(v, uniqueslice) {
			uniqueslice = append(uniqueslice, v)
		}
	}
	return
}

func SliceShuffle(slice []interface{}) []interface{} {
	for i := 0; i < len(slice); i++ {
		a := rand.Intn(len(slice))
		b := rand.Intn(len(slice))
		slice[a], slice[b] = slice[b], slice[a]
	}
	return slice
}

func SliceInsert(slice, insertion []interface{}, index int) []interface{} {
	result := make([]interface{}, len(slice)+len(insertion))
	at := copy(result, slice[:index])
	at += copy(result[at:], insertion)
	copy(result[at:], slice[index:])
	return result
}

//SliceRomove(a,4,5) //a[4]
func SliceRemove(slice []interface{}, start int, args ...int) []interface{} {
	var end int
	if len(args) == 0 {
		end = start + 1
	} else {
		end = args[0]
	}
	return append(slice[:start], slice[end:]...)
}

// ==============================
// http
// ==============================
func HttpPost(client *http.Client, url string, body []byte, header http.Header) (io.ReadCloser, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	for k, vs := range header {
		req.Header[k] = vs
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &RemoteError{req.URL.Host, err}
	}
	if resp.StatusCode == 200 {
		return resp.Body, nil
	}
	resp.Body.Close()
	if resp.StatusCode == 404 { // 403 can be rate limit error.  || resp.StatusCode == 403 {
		err = NotFoundError{"Resource not found: " + url}
	} else {
		err = &RemoteError{req.URL.Host, fmt.Errorf("get %s -> %d", url, resp.StatusCode)}
	}
	return nil, err
}

func HttpPostBytes(client *http.Client, url string, body []byte, header http.Header) ([]byte, error) {
	rc, err := HttpPost(client, url, body, header)
	if err != nil {
		return nil, err
	}
	p, err := ioutil.ReadAll(rc)
	rc.Close()
	return p, nil
}

func HttpPostJSON(client *http.Client, url string, body []byte, header http.Header) ([]byte, error) {
	if header == nil {
		header = http.Header{}
	}
	header.Add("Content-Type", "application/json")
	p, err := HttpPostBytes(client, url, body, header)
	if err != nil {
		return []byte{}, err
	}
	return p, nil
}

// ==============================
// file
// ==============================
// SaveFile saves content type '[]byte' to file by given path.
// It returns error when fail to finish operation.
func SaveFile(filePath string, b []byte) (int, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(b)
}

// SaveFileS saves content type 'string' to file by given path.
// It returns error when fail to finish operation.
func SaveFileS(filePath string, s string) (int, error) {
	return SaveFile(filePath, []byte(s))
}

// ReadFile reads data type '[]byte' from file by given path.
// It returns error when fail to finish operation.
func ReadFile(filePath string) ([]byte, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte(""), err
	}
	return b, nil
}

// ReadFileS reads data type 'string' from file by given path.
// It returns error when fail to finish operation.
func ReadFileS(filePath string) (string, error) {
	b, err := ReadFile(filePath)
	return string(b), err
}

// Unzip unzips .zip file to 'destPath' and returns sub-directories.
// It returns error when fail to finish operation.
func Unzip(srcPath, destPath string) ([]string, error) {
	// Open a zip archive for reading
	r, err := zip.OpenReader(srcPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	dirs := make([]string, 0, 5)
	// Iterate through the files in the archive
	for _, f := range r.File {
		// Get files from archive
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}

		dir := path.Dir(f.Name)
		// Create directory before create file
		os.MkdirAll(destPath+"/"+dir, os.ModePerm)
		dirs = AppendStr(dirs, dir)

		if f.FileInfo().IsDir() {
			continue
		}

		// Write data to file
		fw, _ := os.Create(path.Join(destPath, f.Name))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(fw, rc)

		if fw != nil {
			fw.Close()
		}
		if err != nil {
			return nil, err
		}
	}
	return dirs, nil
}

func TarGz(srcDirPath string, destFilePath string) error {
	fw, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer fw.Close()

	// Gzip writer
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// Tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Check if it's a file or a directory
	f, err := os.Open(srcDirPath)
	if err != nil {
		return err
	}
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	if fi.IsDir() {
		// handle source directory
		fmt.Println("Cerating tar.gz from directory...")
		if err := tarGzDir(srcDirPath, path.Base(srcDirPath), tw); err != nil {
			return err
		}
	} else {
		// handle file directly
		fmt.Println("Cerating tar.gz from " + fi.Name() + "...")
		if err := tarGzFile(srcDirPath, fi.Name(), tw, fi); err != nil {
			return err
		}
	}
	fmt.Println("Well done!")
	return err
}

// Deal with directories
// if find files, handle them with tarGzFile
// Every recurrence append the base path to the recPath
// recPath is the path inside of tar.gz
func tarGzDir(srcDirPath string, recPath string, tw *tar.Writer) error {
	// Open source diretory
	dir, err := os.Open(srcDirPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	// Get file info slice
	fis, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		// Append path
		curPath := srcDirPath + "/" + fi.Name()
		// Check it is directory or file
		if fi.IsDir() {
			// Directory
			// (Directory won't add unitl all subfiles are added)
			fmt.Printf("Adding path...%s\n", curPath)
			tarGzDir(curPath, recPath+"/"+fi.Name(), tw)
		} else {
			// File
			fmt.Printf("Adding file...%s\n", curPath)
		}

		tarGzFile(curPath, recPath+"/"+fi.Name(), tw, fi)
	}
	return err
}

// Deal with files
func tarGzFile(srcFile string, recPath string, tw *tar.Writer, fi os.FileInfo) error {
	if fi.IsDir() {
		// Create tar header
		hdr := new(tar.Header)
		// if last character of header name is '/' it also can be directory
		// but if you don't set Typeflag, error will occur when you untargz
		hdr.Name = recPath + "/"
		hdr.Typeflag = tar.TypeDir
		hdr.Size = 0
		//hdr.Mode = 0755 | c_ISDIR
		hdr.Mode = int64(fi.Mode())
		hdr.ModTime = fi.ModTime()

		// Write hander
		err := tw.WriteHeader(hdr)
		if err != nil {
			return err
		}
	} else {
		// File reader
		fr, err := os.Open(srcFile)
		if err != nil {
			return err
		}
		defer fr.Close()

		// Create tar header
		hdr := new(tar.Header)
		hdr.Name = recPath
		hdr.Size = fi.Size()
		hdr.Mode = int64(fi.Mode())
		hdr.ModTime = fi.ModTime()

		// Write hander
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}

		// Write file data
		_, err = io.Copy(tw, fr)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnTarGz ungzips and untars .tar.gz file to 'destPath' and returns sub-directories.
// It returns error when fail to finish operation.
func UnTarGz(srcFilePath string, destDirPath string) ([]string, error) {
	// Create destination directory
	os.Mkdir(destDirPath, os.ModePerm)

	fr, err := os.Open(srcFilePath)
	if err != nil {
		return nil, err
	}
	defer fr.Close()

	// Gzip reader
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	// Tar reader
	tr := tar.NewReader(gr)

	dirs := make([]string, 0, 5)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}

		// Check if it is directory or file
		if hdr.Typeflag != tar.TypeDir {
			// Get files from archive
			// Create directory before create file
			dir := path.Dir(hdr.Name)
			os.MkdirAll(destDirPath+"/"+dir, os.ModePerm)
			dirs = AppendStr(dirs, dir)

			// Write data to file
			fw, _ := os.Create(destDirPath + "/" + hdr.Name)
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				return nil, err
			}
		}
	}
	return dirs, nil
}

func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// search a file in paths.
// this is offen used in search config file in /etc ~/
func SearchFile(filename string, paths ...string) (fullpath string, err error) {
	for _, path := range paths {
		if fullpath = filepath.Join(path, filename); FileExists(fullpath) {
			return
		}
	}
	err = errors.New(fullpath + " not found in paths")
	return
}

// like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
func GrepFile(patten string, filename string) (lines []string, err error) {
	re, err := regexp.Compile(patten)
	if err != nil {
		return
	}

	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	lines = make([]string, 0)
	reader := bufio.NewReader(fd)
	prefix := ""
	for {
		byteLine, isPrefix, er := reader.ReadLine()
		if er != nil && er != io.EOF {
			return nil, er
		}
		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		}

		line = prefix + line
		if re.MatchString(line) {
			lines = append(lines, line)
		}
		if er == io.EOF {
			break
		}
	}
	return lines, nil
}
