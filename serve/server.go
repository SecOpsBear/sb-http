package serve

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const tmpl = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>Upload File</title>
  </head>
  <body>
  &#60;<a href="/">Home</a>
  <link rel="shortcut icon" href="#" />
  <br><br>
    <form
      enctype="multipart/form-data"
      action="/upload"
      method="post"
    >
      <input type="file" name="multiplefiles" id="multiplefiles" multiple />
      <input type="submit" name="submit" value="upload" />
    </form>
  </body>
</html>
`

// Compile templates on start of the application
var templates, _ = template.New("upload.html").Parse(tmpl)

func UploadFilesGet(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "upload.html", nil)
}

func UploadFilesPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// w.Header().Set("link", "rel=\"shortcut icon\" href=\"#\"")
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	//get the fileheaders
	files := r.MultipartForm.File["multiplefiles"] // grab the filenames
	var st []string
	st = append(st, "<pre>\n<link rel=\"shortcut icon\" href=\"#\" />")

	for i, _ := range files { // loop through the files one by one
		file, err := files[i].Open()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		defer file.Close()

		out, err := os.Create(files[i].Filename)
		if err != nil {
			fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file) // file not files[i] !

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		st = append(st, "Files uploaded successfully : "+files[i].Filename+"\n")
		fmt.Printf("Upload file: %+v\n", files[i].Filename)
	}

	st = append(st, "</pre>\n")
	buff := &bytes.Buffer{}
	buff.WriteString(strings.Join(st, ""))
	if _, err := io.Copy(w, buff); err != nil {
		log.Printf("Failed to send out response: %v", err)
	}
}
