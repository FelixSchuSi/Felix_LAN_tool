package main

import (
	_ "embed"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

//go:embed static/bootstrap.min.css
var bootstrap_content []byte

func index_handler(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir("./")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("couldnt read directory './'."))
		return
	}

	var filesHtml []string

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		var fileHtml = `
<li id="list-item" class="list-group-item list-group-item-action">
	<span class="nameElem">` + e.Name() + `</span>
	<span class="btnElem">
		<form action="/download/` + e.Name() + `" method="POST">
			<button class="btn btn-outline-dark">Download</button>
		</form>
	</span>
</li>`
		filesHtml = append(filesHtml, fileHtml)
	}

	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<link rel="stylesheet" href="/static/bootstrap.min.css">
		<style>
			.btnElem {
				display: inline-flex;
				flex-shrink: 0;
				width: 4rem;
				height: 2rem;
				justify-content: center;
				align-items: center;
				float: right;
				margin-right: 1rem;
			}
			#list-item,
			#upload-form {
				display: flex;
				align-items: center;
				justify-content: space-between;
			} 

			#upload-form {
				width: 100%%;
			}

			#ul {
				margin-top: 0;
			}
		</style>
		<title>Felix LAN Tool</title>
	</head>

	<body>
		<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
			<a class="navbar-brand" href="/">Felix LAN Tool 🐷🐷🐷</a>
		</nav>
		<main class="container">
		<ul id="ul" class="list-group mt-3">
			`+strings.Join(filesHtml, "")+`

			<li id="list-item" class="list-group-item list-group-item-action">
				<form id="upload-form" enctype="multipart/form-data" action="/upload" method="post">
					<input type="file" name="myFile" />
					<button class="btn btn-outline-dark" type="submit">Upload</button>
				</form>
			</li>
		</ul>
		</main>
  </body>
</html>`)
}

func setHeadersForDownload(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment")
		fmt.Println("Someone is downloading", r.URL.Path)
		fs.ServeHTTP(w, r)
	}
}

func upload_handler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(100000 << 20)
	formFile, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error retrieving the file")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Upload failed: error retrieving the file"))
		return
	}
	defer formFile.Close()

	fmt.Println("Someone is uploading", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)

	fileBytes, err := io.ReadAll(formFile)
	if err != nil {
		fmt.Println("Reading files from multipartform failed")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		fmt.Println("Reading files from multipartform failed")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !strings.HasSuffix(handler.Filename, ".zip") {
		fmt.Println("Upload failed: you can only upload '.zip' files")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Upload failed: you can only upload '.zip' files"))
		return
	}

	files, err := os.ReadDir("./")
	if err != nil {
		fmt.Println("Cannot read dir")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Upload failed"))
		return
	}
	for _, existingFileName := range files {
		if existingFileName.Name() == handler.Filename {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Upload failed: file already exists"))
			return
		}
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}

	err = os.WriteFile("./"+handler.Filename, fileBytes, 0644)
	if err != nil {
		fmt.Println("Writing file failed")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully uploaded file: "+handler.Filename+"\n")
}

func bootstrap_handler(w http.ResponseWriter, r *http.Request, bootstrap_content []byte) {
	w.Header().Set("Content-Type", "text/css")
	w.WriteHeader(http.StatusOK)
	w.Write(bootstrap_content)
}

func print_local_ip_addresses() {
	addrs, _ := net.InterfaceAddrs()
	println("running on:")
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && ipnet.IP.IsPrivate() == true {
				println("  http://" + ipnet.IP.String())
			}
		}
	}
}

func main() {
	http.HandleFunc("/static/bootstrap.min.css", func(w http.ResponseWriter, r *http.Request) { bootstrap_handler(w, r, bootstrap_content) })
	http.HandleFunc("/", index_handler)

	fs := http.FileServer(http.Dir("./"))

	http.Handle("/download/", http.StripPrefix("/download/", setHeadersForDownload(fs)))
	http.HandleFunc("/upload", upload_handler)

	print_local_ip_addresses()

	http.ListenAndServe("0.0.0.0:80", nil)
}
