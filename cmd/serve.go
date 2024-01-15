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
	"github.com/microcosm-cc/bluemonday"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "View Your Memo in the browser",
	Long:  `View Your Memo in your favourite broswer!`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			viewMemo, _ := strconv.Atoi(args[0])
			filename := matchMemoNumber(viewMemo)
			content, _ := os.ReadFile(filename)
			userHTML := mdToHTML(content)
			finHTML := injectAdditionalHtml(userHTML, filename)
			startServer(finHTML)
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
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		displayCustom404(w, r)
		return
	}
	homeFiles := filepath.Join(getKeyValue("MemoDir").(string), "assets")
	baseFile := filepath.Join(homeFiles, "base.html")
	cssContent, _ := os.ReadFile(filepath.Join(homeFiles, "base.css"))
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

	data := inputData{Title: "Home", Main: template.HTML(forwardContent), StyleSheet: template.CSS(cssContent)}
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
	homeFiles := filepath.Join(getKeyValue("MemoDir").(string), "assets")

	baseFile := filepath.Join(homeFiles, "base.html")
	ts, err := template.New("base.html").ParseFiles(baseFile)

	cssContent, _ := os.ReadFile(filepath.Join(homeFiles, "base.css"))
	data := inputData{Title: getFileTitle(matchMemoNumber(id)), Main: template.HTML(userHTML), StyleSheet: template.CSS(cssContent)}
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
	homeFiles := filepath.Join(getKeyValue("MemoDir").(string), "assets")

	baseFile := filepath.Join(homeFiles, "base.html")
	ts, err := template.New("base.html").ParseFiles(baseFile)

	cssContent, _ := os.ReadFile(filepath.Join(homeFiles, "base.css"))
	data := inputData{Title: "404 Page Not Found", Main: template.HTML(htmlCode), StyleSheet: template.CSS(cssContent)}
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

// The code inside here should be replaced!
func injectAdditionalHtml(bodyContent []byte, filename string) string {
	assetsDir := filepath.Join(getKeyValue("configDir").(string), "assets")
	getTopFileContent, _ := os.ReadFile(filepath.Join(assetsDir, "top.html"))
	fileTitle := getFileTitle(filename)

	// Replace the <title> tag with the fileTitle using regular expression
	titlePattern := regexp.MustCompile(`<title>\s*(.*)\s*</title>`)
	topHtmlWithFileTitle := titlePattern.ReplaceAllString(string(getTopFileContent), "<title>"+fileTitle+"</title>")

	// Inject bodyContent to the <body> tag using regular expression
	sanitizedHTML := bluemonday.UGCPolicy().SanitizeBytes(bodyContent)
	bodyPattern := regexp.MustCompile(`<body>\s*(.*)\s*</body>`)
	fullHtml := bodyPattern.ReplaceAllString(topHtmlWithFileTitle, "<body>\n"+string(sanitizedHTML)+"</body>")

	return fullHtml
}

func mdToHTML(md []byte) []byte {
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
func startServer(htmlContent string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve the generated HTML content
		fmt.Fprintf(w, htmlContent)
	})

	// Start the server on port 8080
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
