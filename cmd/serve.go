/*
Copyright © 2024 Gekko Wrld

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"bufio"
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

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "View Your Memo in the browser",
	Long:  `View Your Memo in your favourite broswer!`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Print("NotYetImplemented!")
			displayIndex()
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
	Title      string
	Main       template.HTML
	StyleSheet template.CSS
  ScriptSheet template.JS
  FaviconImage template.HTML
}

func serveStaticFile(fileType string) string {
	homeFiles := filepath.Join(getKeyValue("configDir").(string))
	readFile, err := os.Open(filepath.Join(homeFiles, "staticfiles"))
	if err != nil {
		log.Print(err)
		return ""
	}

  if fileType == "favicon" {
    faviconPath := filepath.Join(homeFiles, "assets", "favicon.ico")
    faviconLink := fmt.Sprintf("<link rel=\"shortcut icon\" href=\"%s\" type=\"image/x-icon\">", faviconPath)
    return faviconLink
  }

	var fileContent string
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// This file is expected to be small, so it is reasonable to
	// read and subsequently process it
	for fileScanner.Scan() {
		text := fileScanner.Text()
		// Check if the path is absolute or not
		filelocation := text
		if !filepath.IsAbs(filelocation) {
			filelocation = filepath.Join(homeFiles, text)
		}
		fileExt := filepath.Ext(filelocation)
    // Remove the . before actually doing anything
    var withoutDot string
    if fileExt != "" {
      withoutDot = fileExt[1:]
    }
    if withoutDot == fileType {
			fileByteCont, err := os.ReadFile(filelocation)
			if err != nil {
				log.Print(err)
				continue
			}
			fileContent += string(fileByteCont)
		}
	}

	readFile.Close()

	return fileContent
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		displayCustom404(w, r)
		return
	}
	homeFiles := filepath.Join(getKeyValue("configDir").(string), "assets")
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
			forwardContent += fmt.Sprintf("<a class=\"main-link\" href=\"/view?id=%s\">%s - %s</a><br/>", fileNumber[1], fileNumber[1], fileTitle)
		}
	}

	data := inputData{Title: "Home", Main: template.HTML(forwardContent), StyleSheet: template.CSS(serveStaticFile("css")), ScriptSheet: template.JS(serveStaticFile("js")), FaviconImage: template.HTML(serveStaticFile("favicon"))}
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
	homeFiles := filepath.Join(getKeyValue("configDir").(string), "assets")

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
	homeFiles := filepath.Join(getKeyValue("configDir").(string), "assets")

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
	log.Print("Starting on Server: 4000")
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
