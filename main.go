package main

import (
	"os"
	"html/template"
	"net/http"
	"log"
	"io/ioutil"
	"strconv"
)

type Directory struct {
	Dir string
	Files []File
	SubDir string
}

type File struct {
	Path string
	Name string
	IsDir bool
}

var rootDir string

func getContentType(file string) (string, error) {
	f,err := os.Open(file)

	if err != nil {
		return "", err
	}

	defer f.Close()

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}

func listDir(directory string) (*Directory, error) {
 	var files []File

	filesInDir, err := ioutil.ReadDir(directory)
	
	if err != nil {
		return nil, err
	}

	for _, file := range filesInDir {
		files = append(files, File{Name: file.Name(), Path: directory+"/"+ file.Name(), IsDir: file.IsDir()})
	}

	return &Directory{Dir: directory, SubDir: "" , Files: files},nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	filePath := rootDir + path
	file, err := os.Open(filePath)
	defer file.Close()	
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	fileInfo, err := file.Stat()


	if fileInfo.IsDir() {
		dir,_ := listDir(filePath)
		if path != "/" {
			dir.SubDir = path
		}
		t, _ := template.ParseFiles("listing.template")
		t.Execute(w, dir)
		return
	}

	
	contentType, _ := getContentType(filePath)

	body, _ := ioutil.ReadFile(filePath)
	w.Header().Add("Content-Type", contentType)
	w.Header().Add("Content-Length", strconv.Itoa(len(body)))
	w.Header().Add("Pragma", "No-Cache")
	w.Header().Add("Cache", "max-age=0")
	w.Write(body)
}

func main() {
	rootDir = os.Args[1]
	port := 8080
	var err error
	
	if len(os.Args) == 3 {
		port, err = strconv.Atoi(os.Args[2])

		if err != nil {
			log.Fatalf("Invalid port number", err)
		}
	}	

	http.HandleFunc("/", handler)
	log.Printf("Minions starting on ::: %d", port)
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), nil))
}
