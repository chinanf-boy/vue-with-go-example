package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	fp "path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/chinanf-boy/vue-with-go-example/assets"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
)

var (
	// CTX is a context for webdav vfs
	CTX = context.Background()

	// FS is a virtual memory file system
	FS = webdav.NewMemFS()

	serveCmd = &cobra.Command{
		Use:   "good-serve",
		Short: "Serve 打开服务器 随便哪里",
		Long:  "运行一个简单且高性能的Web服务器。如果不使用--port标志，则默认使用端口8080。",
		Run: func(cmd *cobra.Command, args []string) {

			// Create router
			router := httprouter.New()
			router.GET("/js/*filepath", serveFiles) // 静态资源 路径

			router.GET("/", serveIndexPage)

			// Route for panic
			router.PanicHandler = func(w http.ResponseWriter, r *http.Request, arg interface{}) {
				http.Error(w, fmt.Sprint(arg), 500)
			}

			port, _ := cmd.Flags().GetInt("port")
			url := fmt.Sprintf(":%d", port)
			logrus.Infoln("Serve vue-with-go in", url)
			svr := &http.Server{
				Addr:         url,
				Handler:      router,
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 20 * time.Second,
			}
			logrus.Fatalln(svr.ListenAndServe())
		},
	}
)

func init() {
	serveCmd.Flags().IntP("port", "p", 8070, "Port that used by server")
	rootCmd.AddCommand(serveCmd)
}

func serveFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Read asset path
	path := r.URL.Path
	if path[0:1] == "/" {
		path = path[1:]
	}

	// Load asset
	asset, err := assets.ReadFile(path)
	check(err)

	// Set response header content type
	ext := fp.Ext(path)
	mimeType := mime.TypeByExtension(ext)
	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}

	// Serve asset
	buffer := bytes.NewBuffer(asset)
	io.Copy(w, buffer)
}

func serveIndexPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check token
	// err := checkToken(r)
	// if err != nil {
	// 	redirectPage(w, r, "/login")
	// 	return
	// }

	asset, _ := assets.ReadFile("index.html")
	w.Header().Set("Content-Type", "text/html")
	buffer := bytes.NewBuffer(asset)
	io.Copy(w, buffer)

}
func check(err error) {
	if err != nil {
		panic(err)
	}
}
