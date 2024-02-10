/*
Copyright © 2024 Gekko Wrld
*/
package cmd

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"

	"os"

	"github.com/spf13/cobra"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var memoNumber int

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "View Your Memo in the browser",
	Long:  `View Your Memo in your favourite broswer!`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			memoNumber, _ = strconv.Atoi(args[0])
			mux := http.NewServeMux()
			mux.HandleFunc("/", displayIndividualFile)
			log.Print("Server started on http://127.0.0.0:4000")
			err := http.ListenAndServe(":4000", mux)
			if err != nil {
				log.Print(err)
			}
		} else {
			// Display The Index File
			displayIndex()
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

type inputData struct {
	Title       string
	Main        template.HTML
	StyleSheet  template.CSS
	ScriptSheet template.JS
}

func serveStaticFile(fileType string) string {
	var fileContent string
	staticFiles := getKeyValue("StaticFiles").(string) // Files directory
	filesInDir, err := os.ReadDir(staticFiles)
	if err != nil {
		log.Print(err)
	}

	jsFiles := regexp.MustCompile(`(?mi)\.js`)
	cssFiles := regexp.MustCompile(`(?mi)\.css`)

	for _, file := range filesInDir {
		fileName := filepath.Join(staticFiles, file.Name())
		if jsFiles.MatchString(fileName) && fileType == "js" {
			fileByteCont, err := os.ReadFile(fileName)
			if err != nil {
				log.Print(err)
				continue
			}
			fileContent += string(fileByteCont)
		}
		if cssFiles.MatchString(fileName) && fileType == "css" {
			fileByteCont, err := os.ReadFile(fileName)
			if err != nil {
				log.Print(err)
				continue
			}
			fileContent += string(fileByteCont)
		}
	}

	return fileContent
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		displayCustom404(w, r)
		return
	}
	homeFiles := filepath.Join(getKeyValue("StaticFiles").(string))
	baseFile := filepath.Join(homeFiles, "base.html")
	ts, err := template.New("base.html").ParseFiles(baseFile)

	// Now Get the files
	memoDir := getKeyValue("MemoDir").(string)
	files, _ := os.ReadDir(memoDir)

	var forwardContent string
	for _, file := range files {
		if !file.IsDir() {
			fileTitle := getFileTitle(filepath.Join(memoDir, file.Name()))
			re := regexp.MustCompile(`^(\d+)-`)
			fileNumber := re.FindSubmatch([]byte(file.Name()))
			forwardContent += fmt.Sprintf("<a class=\"main-link\" href=\"/view?id=%s\">%s (%s)</a><br/>", fileNumber[1], fileTitle, fileNumber[1])
		}
	}

	data := inputData{Title: "Home", Main: template.HTML(forwardContent), StyleSheet: template.CSS(serveStaticFile("css")), ScriptSheet: template.JS(serveStaticFile("js"))}
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func displayIndividualFile(w http.ResponseWriter, r *http.Request) {
	filename := matchMemoNumber(memoNumber)
	content, err := os.ReadFile(filename)
	if err != nil {
		displayCustom404(w, r)
		return
	}

	userHTML := mdToHTML(content)

	ftitle := getFileTitle(filename)

	homeFiles := filepath.Join(getKeyValue("StaticFiles").(string))
	baseFile := filepath.Join(homeFiles, "base.html")
	ts, err := template.New("base.html").ParseFiles(baseFile)

	data := inputData{
		Title:       ftitle,
		Main:        template.HTML(userHTML),
		StyleSheet:  template.CSS(serveStaticFile("css")),
		ScriptSheet: template.JS(serveStaticFile("js")),
	}

	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func viewFile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		displayCustom404(w, r)
		return
	}

	viewMemo := id
	filename := matchMemoNumber(viewMemo)
	content, err := os.ReadFile(filename)
	if err != nil {
		displayCustom404(w, r)
		return
	}
	userHTML := mdToHTML(content)
	homeFiles := filepath.Join(getKeyValue("StaticFiles").(string))

	baseFile := filepath.Join(homeFiles, "base.html")
	ts, err := template.New("base.html").ParseFiles(baseFile)

	data := inputData{Title: getFileTitle(matchMemoNumber(id)), Main: template.HTML(userHTML), StyleSheet: template.CSS(serveStaticFile("css")), ScriptSheet: template.JS(serveStaticFile("js"))}
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func displayCustom404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	htmlCode := `
  	<div class="custom_404">
		<div class="custom_border_text">
			<p>
				404 Page Not Found
			</p>
			<p>You can go <a href="/">Home</a> to view all Memo Listings</p>
		</div>
	</div>
  `
	homeFiles := filepath.Join(getKeyValue("StaticFiles").(string))

	baseFile := filepath.Join(homeFiles, "base.html")
	ts, err := template.New("base.html").ParseFiles(baseFile)

	data := inputData{Title: "404 Page Not Found", Main: template.HTML(htmlCode), StyleSheet: template.CSS(serveStaticFile("css"))}
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func displayIndex() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/view", viewFile)
	log.Print("Serve opened at http://127.0.0.0:4000/")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

func mdToHTML(md []byte) []byte {

	// No checks or sanitization provided yet!
	// It is expected to run on users system so most of the attacks
	// are not an imminent threat, for now

	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
