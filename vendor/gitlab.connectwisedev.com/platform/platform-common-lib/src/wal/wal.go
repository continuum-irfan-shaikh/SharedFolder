package wal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

const separator = '\n'

// WAL - Interface to perform action on wal file
type WAL interface {
	// Write - write messages in the wal file
	Write(...string) error

	// WriteObject - serialize object and write this into wal file
	WriteObject(...interface{}) error

	//Read - x wal files, and convert all the enrties in to records
	Read(batchSize int) ([]Record, error)

	// Flush - Commit / Lock / Rotate current wal file, and
	// create a new one for wal messages
	Flush() error
}

// Record - Holds all the Messages from a file and
// Error is there is problem in reading
type Record struct {
	fileName string

	// Messages - All the messages from a wal file
	Messages []string

	// Err - Any error accured while reading wal file
	Err error
}

// Commit - Commit and Remove old wal files
func (r *Record) Commit() error {
	if r.fileName != "" && r.fileName != "." {
		return os.Remove(r.fileName)
	}
	return nil
}

type walImpl struct {
	conf   *Config
	writer *lumberjack.Logger
}

// Create - returns instance of a wal service
var Create = func(conf *Config) WAL {
	w := &walImpl{conf: conf}
	w.createWriter()
	return w
}

func (w *walImpl) createWriter() {
	w.writer = &lumberjack.Logger{
		Filename:   w.conf.Name,
		MaxSize:    (w.conf.MaxSegments * w.conf.SegmentSizeInKB) / 1000, // megabytes
		MaxBackups: w.conf.MaxFiles,
		MaxAge:     w.conf.MaxAgeDays,
		Compress:   false,
		LocalTime:  false,
	}
}

// Write - write messages in the wal file
func (w *walImpl) Write(messages ...string) error {
	replacer := func(r rune) rune {
		if r == separator {
			return -1
		}
		return r
	}
	buf := make([]byte, 0)
	for _, message := range messages {
		buf = append(buf, strings.Map(replacer, message)...)
		buf = append(buf, separator)
	}
	_, err := w.writer.Write(buf)
	return err
}

// WriteObject - serialize object and write this into wal file
func (w *walImpl) WriteObject(objects ...interface{}) error {
	messages := make([]string, len(objects))
	for index, obj := range objects {
		data, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		messages[index] = string(data)
	}
	return w.Write(messages...)
}

// Read - x wal files, and convert all the enrties in to records
func (w *walImpl) Read(batchSize int) ([]Record, error) {
	size, files, err := w.readFiles(batchSize)
	if err != nil {
		return nil, err
	}

	dir := path.Dir(w.conf.Name)
	records := make([]Record, size)
	for index := 0; index < size; index++ {
		path := path.Join(dir, files[index].Name())
		data, err := ioutil.ReadFile(path)
		msg := strings.Split(string(data), string(separator))
		records[index] = Record{
			fileName: path,
			Messages: msg[:len(msg)-1],
			Err:      err,
		}
	}
	return records, nil
}

// Flush - Commit / Lock / Rotate current wal file, and
// create a new one for wal messages
func (w *walImpl) Flush() error {
	return w.writer.Rotate()
}

func (w *walImpl) readFiles(batchSize int) (int, []os.FileInfo, error) {
	files, err := ioutil.ReadDir(path.Dir(w.conf.Name))
	if err != nil {
		return 0, files, err
	}

	size := len(files) - 1 // to avoid current file
	if batchSize != 0 && size > batchSize {
		size = batchSize
	}
	return size, files, nil
}
