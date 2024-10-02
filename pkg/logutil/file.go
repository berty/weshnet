package logutil

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"go.uber.org/multierr"
	"moul.io/u"

	"berty.tech/weshnet/pkg/errcode"
)

func newFileWriteCloser(target, kind string) (io.WriteCloser, error) {
	var filename string
	switch {
	case strings.HasSuffix(target, ".log"): // use the indicated 'path' as filename
		filename = target
	default: // automatically create a new file in the 'path' directory following a pattern
		startTime := time.Now().Format(filePatternDateLayout)
		filename = filepath.Join(
			target,
			fmt.Sprintf("%s-%s.log", kind, startTime),
		)
		// run gc
		{
			err := LogfileGC(target, 20)
			if err != nil {
				return nil, err
			}
		}
	}

	if dir := filepath.Dir(filename); !u.DirExists(dir) {
		err := os.MkdirAll(dir, 0o711)
		if err != nil {
			return nil, errcode.ErrCode_TODO.Wrap(err)
		}
	}

	var writer io.WriteCloser
	if u.FileExists(filename) {
		var err error
		writer, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return nil, errcode.ErrCode_TODO.Wrap(err)
		}
	} else {
		var err error
		writer, err = os.Create(filename)
		if err != nil {
			return nil, errcode.ErrCode_TODO.Wrap(err)
		}
	}

	return writer, nil
}

type Logfile struct {
	Dir    string
	Name   string
	Size   int64
	Kind   string
	Time   time.Time
	Latest bool
	Errs   error `json:"Errs,omitempty"`
}

func (l Logfile) Path() string {
	return filepath.Join(l.Dir, l.Name)
}

const filePatternDateLayout = "2006-01-02T15-04-05.000"

var filePatternRegex = regexp.MustCompile(`(?m)^(.*)-(\d{4}-\d{2}-\d{2}T\d{2}-\d{2}-\d{2}.\d{3}).log$`)

func LogfileList(logDir string) ([]*Logfile, error) {
	files, err := os.ReadDir(logDir)
	if err != nil {
		return nil, errcode.ErrCode_TODO.Wrap(err)
	}

	infos := make([]fs.FileInfo, 0, len(files))
	for _, entry := range files {
		info, err := entry.Info()
		if err != nil {
			return nil, errcode.ErrCode_TODO.Wrap(err)
		}
		infos = append(infos, info)
	}

	logfiles := []*Logfile{}
	for _, info := range infos {
		sub := filePatternRegex.FindStringSubmatch(info.Name())
		if sub == nil {
			continue
		}
		t, err := time.Parse(filePatternDateLayout, sub[2])
		var errs error
		if err != nil {
			errs = multierr.Append(errs, err)
		}

		// use os.Stat to get the file size (updated than fs.FileInfo.Size()
		filepath := filepath.Join(logDir, info.Name())
		fi, err := os.Stat(filepath)
		if err != nil {
			errs = multierr.Append(errs, err)
		}

		logfiles = append(logfiles, &Logfile{
			Dir:  logDir,
			Name: info.Name(),
			Size: fi.Size(),
			Kind: sub[1],
			Time: t,
			Errs: errs,
		})
	}

	// compute latest
	if len(logfiles) > 0 {
		var maxTime time.Time
		for _, file := range logfiles {
			if file.Time.After(maxTime) {
				maxTime = file.Time
			}
		}
		for _, file := range logfiles {
			if file.Time == maxTime {
				file.Latest = true
			}
		}
	}

	return logfiles, nil
}

func CurrentLogfilePath(target string) (string, error) {
	filename := ""

	switch {
	case strings.HasSuffix(target, ".log"): // use the indicated 'path' as filename
		filename = target
	default: // find the latest log file in the 'path' directory following a pattern
		logfileList, err := LogfileList(target)
		if err != nil {
			return "", err
		}

		for _, logfile := range logfileList {
			if logfile.Latest {
				filename = logfile.Path()
				break
			}
		}
	}

	return filename, nil
}

func LogfileGC(logDir string, max int) error {
	if !u.DirExists(logDir) {
		return nil
	}
	files, err := LogfileList(logDir)
	if err != nil {
		return errcode.ErrCode_TODO.Wrap(err)
	}
	if len(files) < max {
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Time.Before(files[j].Time)
	})

	var errs error
	for i := 0; i < len(files)-max; i++ {
		err := os.Remove(files[i].Path())
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}
	return errs
}
