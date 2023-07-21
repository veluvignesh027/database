package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

func calculateDirectorySize(dirPath string) (int, error) {
	var dirSize int

	finfo, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Println(err)
	}
	for _, file := range finfo {
		dirSize = dirSize + int(file.Size())
	}
	log.Println("Size of Directory ", dirPath, " : ", dirSize)
	return dirSize, err
}

func deletefile(dirPath string) error {
	var oldestFile string
	var oldestModTime time.Time

	fileinfo, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Println(err)
	}
	for _, info := range fileinfo {
		modtime := info.ModTime()
		if oldestModTime.IsZero() || modtime.Before(oldestModTime) {
			oldestModTime = modtime
			oldestFile = path.Join(dirPath, info.Name())
		}
	}
	log.Println("Oldest file : ", oldestFile, "...Ready to Delete...")
	return os.Remove(oldestFile)
}

func freespace(path string, maxlimit int, ch chan int) {
	for {
		dirSize, err := calculateDirectorySize(path)
		if err != nil {
			log.Println(err)
		}

		if dirSize > maxlimit {
			for dirSize > maxlimit {
				log.Println("Directory size is greater than ", maxlimit, ".. Deleting older contents..")
				err := deletefile(path)
				if err != nil {
					log.Println(err)
				}
				log.Println("File Deleted sucessfully..!")

				dirSize, err = calculateDirectorySize(path)
				if err != nil {
					log.Println(err)
				}
			}
			log.Println("Space Cleared....")
		} else {
			log.Println("Enough space is there...!")
		}
		time.Sleep(time.Second * 2)
	}
	<-ch
}

func main() {
	limit := int(10)
	dirPath := "C:/Users/vb/go/db/testingdir"
	ch := make(chan int)
	go freespace(dirPath, limit, ch)
	ch <- 10
}
