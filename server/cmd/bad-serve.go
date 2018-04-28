package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	fp "path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/julienschmidt/httprouter"
)

var (
	badserveCmd = &cobra.Command{
		Use:   "bad-serve",
		Short: "Serve 打开服务器 只能在对应项目下运行",
		Long:  "运行一个简单且高性能的Web服务器。如果不使用--port标志，则默认使用端口8080。",
		Run: func(cmd *cobra.Command, args []string) {

			// Create router
			router := httprouter.New()
			router.GET("/js/*filepath", bserveFiles) // 静态资源 路径

			router.GET("/", bserveIndexPage)

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
	badserveCmd.Flags().IntP("port", "p", 8070, "Port that used by server")
	rootCmd.AddCommand(serveCmd)
}

func bserveFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Read asset path
	path := r.URL.Path
	if path[0:1] == "/" {
		path = path[1:]
	}

	// Load asset
	asset, err := bReadFile(path)
	bcheck(err)

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

func bserveIndexPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// bCheck token
	// err := bcheckToken(r)
	// if err != nil {
	// 	redirectPage(w, r, "/login")
	// 	return
	// }

	asset, err := bReadFile("index.html")
	check(err)
	w.Header().Set("Content-Type", "text/html")
	buffer := bytes.NewBuffer(asset)
	io.Copy(w, buffer)

}

// ReadFile read path
func bReadFile(path string) ([]byte, error) {

	f, err := ioutil.ReadFile(Paths + "/client/" + path)
	if err != nil {
		return nil, err
	}

	// buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))

	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	// _, err = buf.ReadFrom(f)
	return f, err
}

func bcheck(err error) {
	if err != nil {
		panic(err)
	}
}
