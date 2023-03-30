package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/secopsbear/sb-http/serve"
	"github.com/spf13/cobra"
)

var port string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Basic server with upload and download feature",
	Long:  `Basic server with upload and download feature`,
	Run: func(cmd *cobra.Command, args []string) {

		myIpCommand := "myips | awk '{print $1\" \"$3}'"
		_, err := exec.LookPath("myips")
		if err != nil {
			myIpCommand = "ip a | grep -w inet | awk '{print $NF \" \" $2}' | awk -F \"/\" '{print $1}' | grep -v lo"
		}
		var execCmd *exec.Cmd
		// if runtime.GOOS == "windows" {
		// 	execCmd = exec.Command("cmd", "/c", myIpCommand)
		// }
		if runtime.GOOS == "linux" {

			execCmd = exec.Command("bash", "-c", myIpCommand)
			out, err := execCmd.CombinedOutput()
			if err != nil {
				log.Fatalf("Something went wrong: %s \n", err)
			}
			fmt.Print(string(out))
		}

		fmt.Println("Starting Server...")
		fmt.Println("/upload - to upload files")

		mux := http.NewServeMux()
		fs := http.FileServer(http.Dir("."))
		mux.Handle("/", addUploadLink(http.StripPrefix("/", fs)))
		mux.HandleFunc("/upload", UploadHandler)

		fmt.Println("Started listening on :" + port)

		// Listen on default port 8099
		http.ListenAndServe(":"+port, RequestLogger(mux))
	},
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		serve.UploadFilesGet(w, r)
	case "POST":
		uploadFiles(w, r)
	}
}

// addUploadLink adds upload link to the response
func addUploadLink(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h.ServeHTTP(w, r)
		ss := fmt.Sprintln("<a href=\"/\">Home</a>&nbsp;&nbsp;&nbsp;&nbsp;<a href=\"/upload\">Upload files</a>&nbsp;<link rel=\"shortcut icon\" href=\"#\"/><br>")
		buff := &bytes.Buffer{}
		buff.WriteString(ss)
		if _, err := io.Copy(w, buff); err != nil {
			log.Printf("Failed to send out response: %v", err)
		}
	})
}

func uploadFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	// Add links to upload and home
	ss := fmt.Sprintln("&#60;<a href=\"/upload\">Back</a>&nbsp;&nbsp;&nbsp;&nbsp;<a href=\"/\">Home</a><link rel=\"icon\" href=\"data:,\"/><br>")
	buff := &bytes.Buffer{}
	buff.WriteString(ss)
	if _, err := io.Copy(w, buff); err != nil {
		log.Printf("Failed to send out response upload: %v", err)
	}

	serve.UploadFilesPOST(w, r)

}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVarP(&port, "port", "p", "8099", "Enter the port number")

}

// Request logger function and wrap it on the HTTP request multiplexer
func RequestLogger(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		targetMux.ServeHTTP(w, r)

		// log request by who(IP address)
		requesterIP := r.RemoteAddr

		log.Printf(
			"%s\t%s \t%s\t%v",
			r.Method,
			r.RequestURI,
			requesterIP,
			time.Since(start),
		)
	})
}
