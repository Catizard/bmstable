package bmstable

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Catizard/bmstable/internal"
)

// Stores the meta data of a difficult table
type DifficultTable struct {
	HeaderURL   string `json:"header_url"`
	DataURL     string `json:"data_url"`
	Name        string `json:"name"`
	OriginalURL string `json:"original_url"`
	Symbol      string `json:"symbol"`
	LevelOrder  []string
	Contents    []DifficultTableData
	Courses     []CourseInfo
}

// Reprents one song arranged in a difficult table
type DifficultTableData struct {
	Artist   string `json:"artist"`
	Comment  string `json:"comment"`
	Level    string `json:"level"`
	Lr2BmsID string `json:"lr2_bmdid"`
	Md5      string `json:"md5"`
	NameDiff string `json:"name_diff"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	URLDiff  string `json:"url_diff"`
	Sha256   string `json:"sha256"`
}

// Stores the meta data of a course
type CourseInfo struct {
	Name       string   `json:"name"`
	Md5        []string `json:"md5"`
	Sha256     []string `json:"sha256"`
	Constraint []string `json:"constraint"`
}

func ParseFromURL(url string) (DifficultTable, error) {
	header, err := parseHeaderFromURL(url)
	if err != nil {
		return DifficultTable{}, err
	}
	data, err := parseDataFromURL(header.DataURL)
	if err != nil {
		return DifficultTable{}, err
	}
	header.Contents = data
	return *header, nil
}

func parseHeaderFromURL(url string) (*DifficultTable, error) {
	jsonURL := ""
	prefixURL := buildPrefixURL(url)
	if strings.HasSuffix(url, ".html") {
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("net: cannot get contents from %s due to %s", url, err)
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return nil, fmt.Errorf("io: failed to read contents: %s", err)
			}
			line := strings.TrimSpace(scanner.Text())
			// TODO: Any other cases?
			// Its pattern should be <meta name="bmstable" content="xxx.json" />
			if strings.HasPrefix(line, "<meta name=\"bmstable\"") {
				startp := strings.Index(line, "content") + len("content=\"") - 1
				if startp == -1 {
					return nil, fmt.Errorf("unexpected format: cannot find 'content' field in %s", url)
				}
				endp := -1
				// Finds the end position
				first := false
				for i := startp; i < len(line); i++ {
					if line[i] == '"' {
						if !first {
							first = true
						} else {
							endp = i
							break
						}
					}
				}
				if endp == -1 {
					return nil, fmt.Errorf("unexpected format: cannot find 'content' field in %s", url)
				}

				content := line[startp+1 : endp]
				if !strings.HasPrefix(content, "http") {
					// Construct the json url path
					jsonURL = prefixURL + "/" + line[startp+1:endp]
				} else {
					jsonURL = content
				}

				break
			}
		}
	} else if strings.HasSuffix(url, ".json") {
		// Okay dokey
		jsonURL = url
	}

	if jsonURL == "" {
		return nil, fmt.Errorf("parse: cannot build possible json url from %s", url)
	}

	var rawHeader internal.ImportHeaderVo
	if err := fetchJSON(jsonURL, &rawHeader); err != nil {
		return nil, err
	}

	if !strings.HasPrefix(rawHeader.DataURL, "http") {
		rawHeader.DataURL = prefixURL + "/" + rawHeader.DataURL
	}
	if err := rawHeader.ParseRawCourses(); err != nil {
		return nil, err
	}
	return newDifficultTable(&rawHeader), nil
}

func parseDataFromURL(url string) ([]DifficultTableData, error) {
	var data []DifficultTableData
	if err := fetchJSON(url, &data); err != nil {
		return nil, fmt.Errorf("parse: cannot fetch table data from %s due to %s", url, err)
	}
	return data, nil
}

func newDifficultTable(rawHeader *internal.ImportHeaderVo) *DifficultTable {
	courses := make([]CourseInfo, 0)
	for i := range rawHeader.Courses {
		courses = append(courses, *newCourseInfo(&rawHeader.Courses[i]))
	}
	castedLevelOrders := make([]string, 0)
	for _, level := range rawHeader.LevelOrders {
		castedLevel := ""
		if l, ok := level.(string); ok {
			castedLevel = l
		} else if l, ok := level.(int); ok {
			castedLevel = strconv.Itoa(l)
		} else {
			castedLevel = fmt.Sprintf("%v", level)
		}
		castedLevelOrders = append(castedLevelOrders, castedLevel)
	}
	return &DifficultTable{
		HeaderURL:   rawHeader.HeaderURL,
		DataURL:     rawHeader.DataURL,
		Name:        rawHeader.Name,
		OriginalURL: rawHeader.OriginalURL,
		Symbol:      rawHeader.Symbol,
		LevelOrder:  castedLevelOrders,
		Contents:    make([]DifficultTableData, 0),
		Courses:     courses,
	}
}

func newCourseInfo(rawCourse *internal.ImportCourseVo) *CourseInfo {
	return &CourseInfo{
		Name:       rawCourse.Name,
		Md5:        rawCourse.Md5,
		Sha256:     rawCourse.Sha256,
		Constraint: rawCourse.Constraint,
	}
}

func fetchJSON(url string, v any) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("net: cannot get contents from %s due to %s", url, err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io: cannot read contents from %s", url)
	}
	// Hack \ufeff, \r, \n out, especially for PMS tables
	replacer := strings.NewReplacer("\r", "", "\n", "", "\ufeff", "")
	body = []byte(replacer.Replace(string(body)))

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("json: %s", err)
	}
	return nil
}

func buildPrefixURL(url string) string {
	splitURL := strings.Split(url, "/")
	splitURL[len(splitURL)-1] = ""
	return strings.Join(splitURL, "/")
}
