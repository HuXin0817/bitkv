package bitkv

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	dirPath      = "/tmp/bitkv"
	dataFilePath = filepath.Join(dirPath, "bitkv.data")
	dataFile     *os.File
	tmpFilePath  = filepath.Join(dirPath, "bitkv.data.tmp")
	tmpFile      *os.File

	fileSize int
	fileLock sync.Mutex

	buffer = newBucket()
	cache  = newBucket()
)

func init() {
	if err := initDB(); err != nil {
		fmt.Printf("[ERROR] an error occurs when initial the databases, err: %s\n", err)
		os.Exit(0)
	}
}

func initDB() (err error) {
	if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}

	if dataFile, err = os.OpenFile(dataFilePath, os.O_CREATE|os.O_RDWR, 0644); err != nil {
		return err
	}

	info, err := dataFile.Stat()
	if err != nil {
		return err
	}

	fileSize = int(info.Size())

	var content = make([]byte, fileSize)
	if _, err = dataFile.Read(content); err != nil {
		return err
	}

	if cache, err = replay(content); err != nil {
		return err
	}

	if fileSize > (cache.bits * 2) {
		if err = Merge(); err != nil {
			return err
		}
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			if len(buffer.m) > 0 {
				c := buffer.Export()
				fileLock.Lock()
				if n, err := dataFile.WriteAt(c, int64(fileSize)); err != nil {
					fmt.Printf("[ERROR] %v\n", err)
				} else {
					fileSize += n
				}
				fileLock.Unlock()
			}

			if fileSize > (cache.bits << 4) {
				if err := Merge(); err != nil {
					fmt.Printf("[ERROR] %v\n", err)
				}
			}

			fmt.Printf("[INFO] cache bit: %dkb, file bit: %dkb\n", cache.bits>>10, fileSize>>10)
		}
	}()

	return nil
}

// Put a k-v into database
func Put(k, v string) {
	if cache.Put(k, v) {
		go buffer.Put(k, v)
	}
}

// Get the value of the key
func Get(k string) (v string) {
	return cache.Get(k)
}

// Merge the data file
func Merge() (err error) {
	c := cache.Export()
	if tmpFile, err = os.OpenFile(tmpFilePath, os.O_CREATE|os.O_RDWR, 0644); err != nil {
		return err
	}

	if _, err = tmpFile.Write(c); err != nil {
		return err
	}

	if err = tmpFileToData(); err != nil {
		return err
	}

	return nil
}

func tmpFileToData() (err error) {
	fileLock.Lock()
	defer fileLock.Unlock()

	fileSize = cache.bits
	if err = os.Remove(dataFilePath); err != nil {
		return err
	}

	if err = os.Rename(tmpFilePath, dataFilePath); err != nil {
		return err
	}

	return nil
}
