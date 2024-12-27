package file

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var FS FileSystem = osFS{}

type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	io.WriterAt
	Stat() (os.FileInfo, error)
	Name() string
	ReadDir(n int) ([]os.DirEntry, error)
	Readdir(count int) ([]os.FileInfo, error)
	Readdirnames(n int) ([]string, error)
}

type FileSystem interface {
	Open(name string) (File, error)
	Create(name string) (File, error)
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	Stat(name string) (os.FileInfo, error)
	Remove(name string) error
	CreateTemp(dir, pattern string) (File, error)
	MkdirAll(path string, perm os.FileMode) error
	RemoveAll(path string) error
	ReadDir(name string) ([]os.DirEntry, error)
	MkdirTemp(dir, pattern string) (name string, err error)
	TempDir() string
}

// osFS implements FileSystem using the local disk.
type osFS struct{}

func (osFS) Open(name string) (File, error)   { return os.Open(name) }
func (osFS) Create(name string) (File, error) { return os.Create(name) }
func (osFS) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return os.OpenFile(name, flag, perm)
}
func (osFS) Stat(name string) (os.FileInfo, error)         { return os.Stat(name) }
func (osFS) Remove(name string) error                      { return os.Remove(name) }
func (osFS) CreateTemp(dir, pattern string) (File, error)  { return os.CreateTemp(dir, pattern) }
func (osFS) MkdirAll(path string, perms os.FileMode) error { return os.MkdirAll(path, perms) }
func (osFS) RemoveAll(path string) error                   { return os.RemoveAll(path) }
func (osFS) ReadDir(name string) ([]os.DirEntry, error)    { return os.ReadDir(name) }
func (osFS) MkdirTemp(dir, pattern string) (name string, err error) {
	return os.MkdirTemp(dir, pattern)
}
func (osFS) TempDir() string { return os.TempDir() }

// GetContent returns the contents for a File
func GetContent(filename string) ([]byte, error) {
	file, err := FS.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	filesize := fileInfo.Size()
	buffer := make([]byte, filesize)

	bytesRead, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}

	if bytesRead != int(filesize) {
		return nil, errors.New("didn't read all of the File")
	}

	return buffer, nil
}

// GetLines returns lines in a File
func GetLines(filename string) ([]string, error) {
	buffer, err := GetContent(filename)
	if err != nil {
		return nil, err
	}

	rows := strings.Split(string(buffer), "\n")

	return rows, nil
}

// RemoveFile deletes a File
func RemoveFile(filename string) error {
	return FS.Remove(filename)
}

// RemoveAll deletes a File or Directory and all of its contents
func RemoveAll(path string) error {
	var err error
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}
	}
	return FS.RemoveAll(path)
}

// CreateFile opens a file for writing
func CreateFile(path string) (File, error) {
	var err error
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}
	}
	return FS.Create(path)
}

// OpenFile opens a file for reading
func OpenFile(path string) (File, error) {
	var err error
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}
	}
	return FS.Open(path)
}

// CreateTempFile creates a temp File and returns the name
func CreateTempFile(prefix string) (string, error) {
	file, err := FS.CreateTemp("", prefix)
	if err != nil {
		return "", err
	}
	return filepath.Join(FS.TempDir(), file.Name()), nil
}

func IsDir(path string) (bool, error) {
	s, err := FS.Stat(path)
	if err != nil {
		return false, err
	}
	return s.IsDir(), nil
}

// WriteContent writes bytes to a File replacing the existing File or creating new
func WriteContent(filename string, data []byte) error {
	dir := filepath.Dir(filename)
	if _, err := FS.Stat(dir); errors.Is(err, os.ErrNotExist) {
		_ = FS.MkdirAll(dir, 0700) // Create your File
	}

	f, err := FS.Create(filename)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	defer func() {
		_ = f.Close()
	}()

	_, err = w.Write(data)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

// Copy contents of src file to dest file
func Copy(src, dest string) error {
	content, err := GetContent(src)
	if err != nil {
		return err
	}
	return WriteContent(dest, content)
}

// RelFileExists return nil, if path is relative to current working directory, and error otherwise
func RelFileExists(path string) error {
	if filepath.IsAbs(path) {
		return fmt.Errorf("illegal file path: %s", path)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(absPath, filepath.Clean(cwd)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", path)
	}

	fileInfo, err := FS.Stat(absPath)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("path is dir: %s", path)
	}

	return nil
}

// FileExists returns whether a file exists, or error during processing
func FileExists(path string) bool {
	val, err := IsDir(path)
	if err != nil {
		return false
	}
	return !val
}

// DirExists returns whether a file exists, or error during processing
func DirExists(path string) bool {
	val, err := IsDir(path)
	if err != nil {
		return false
	}
	return val
}

// CopyDir will copy the src dir to a dest dir
func CopyDir(srcPath, destPath string) error {
	var err error
	if !filepath.IsAbs(srcPath) {
		srcPath, err = filepath.Abs(srcPath)
		if err != nil {
			return err
		}
	}

	if !filepath.IsAbs(destPath) {
		destPath, err = filepath.Abs(destPath)
		if err != nil {
			return err
		}
	}

	err = FS.MkdirAll(destPath, 0755)
	if err != nil {
		return err
	}

	srcDir, err := FS.Open(srcPath)
	if err != nil {
		return err
	}

	fi, err := srcDir.Stat()
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return errors.New(fmt.Sprintf("%s is not a directory", srcPath))
	}

	entries, err := srcDir.ReadDir(-1)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			err = CopyDir(filepath.Join(srcPath, entry.Name()), filepath.Join(destPath, entry.Name()))
			if err != nil {
				return err
			}
		} else {
			err = Copy(filepath.Join(srcPath, entry.Name()), filepath.Join(destPath, entry.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetDirEntryNames returns a recursive list of dir entry paths, including the srcPath or not
func GetDirEntryNames(srcPath string, includeSrcPath bool) ([]string, error) {
	srcDir, err := FS.Open(srcPath)
	if err != nil {
		return nil, err
	}

	fi, err := srcDir.Stat()
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, errors.New(fmt.Sprintf("%s is not a directory", srcPath))
	}

	entryNames := make([]string, 0, 100)

	entries, err := srcDir.ReadDir(-1)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirEntryNames, err := GetDirEntryNames(entry.Name(), true)
			if err != nil {
				return nil, err
			}
			for _, den := range dirEntryNames {
				if includeSrcPath {
					entryNames = append(entryNames, filepath.Join(srcPath, den))
				} else {
					entryNames = append(entryNames, den)
				}
			}
		} else {
			if includeSrcPath {
				entryNames = append(entryNames, filepath.Join(srcPath, entry.Name()))
			} else {
				entryNames = append(entryNames, entry.Name())
			}
		}
	}

	return entryNames, nil
}

// CreateTempDir creates a temp directory using the pattern to generate a name
func CreateTempDir(dir, pattern string) (name string, err error) {
	return FS.MkdirTemp(dir, pattern)
}

// RemoveEntriesFromDir removes directory contents
func RemoveEntriesFromDir(dir string) error {
	files, err := FS.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		err = FS.RemoveAll(filepath.Join(dir, f.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

// GetDirnames return paths for directories within a directory
func GetDirnames(dir string) ([]string, error) {
	names := make([]string, 0, 20)
	files, err := FS.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			dirName := filepath.Join(dir, f.Name())
			names = append(names, dirName)
		}
	}
	return names, nil
}

// Unzip unzips an archive specified by srcFilePath into destDirPath
func Unzip(srcFilePath, destDirPath string) error {
	srcFile, err := OpenFile(srcFilePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = srcFile.Close()
	}()

	fi, err := srcFile.Stat()
	if err != nil {
		return err
	}

	r, err := zip.NewReader(srcFile, fi.Size())
	if err != nil {
		return err
	}

	err = FS.MkdirAll(destDirPath, 0755)
	if err != nil {
		return err
	}

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			_ = rc.Close()
		}()

		path := filepath.Join(destDirPath, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(destDirPath)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			err = FS.MkdirAll(path, 0755)
			if err != nil {
				return err
			}
		} else {
			err = FS.MkdirAll(filepath.Dir(path), 0755)
			if err != nil {
				return err
			}
			f, err := FS.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return err
			}
			defer func() {
				_ = f.Close()
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
