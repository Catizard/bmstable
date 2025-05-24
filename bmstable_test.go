package bmstable_test

import (
	"testing"

	"github.com/Catizard/bmstable"
)

type diffTableDefinition struct {
	name       string
	url        string
	symbol     string
	hasCourses bool
}

// Real table definition
var realTableDefintion = [...]diffTableDefinition{
	{"通常難易度表", "http://zris.work/bmstable/normal/normal_header.json", "☆", false},
	{"発狂BMS難易度表", "http://zris.work/bmstable/insane/insane_header.json", "★", false},
	{"第三期Overjoy", "http://zris.work/bmstable/overjoy/header.json", "★★", false},
	{"NEW GENERATION 通常難易度表", "http://zris.work/bmstable/normal2/header.json", "▽", true},
	{"NEW GENERATION 発狂難易度表", "http://zris.work/bmstable/insane2/insane_header.json", "▼", true},
	{"Satellite", "http://zris.work/bmstable/satellite/header.json", "sl", true},
	{"Stella", "https://stellabms.xyz/st/table.html", "st", true},
	{"DP Satellite", "http://zris.work/bmstable/dp_satellite/header.json", "DPsl", true},
	{"DP Stella", "http://zris.work/bmstable/dp_stella/header.json", "DPst", false},
	{"δ難易度表", "http://zris.work/bmstable/dp_normal/dpn_header.json", "δ", true},
	{"発狂DP難易度表", "http://zris.work/bmstable/dp_insane/dpi_header.json", "★", true},
	{"DP Overjoy", "http://zris.work/bmstable/dp_overjoy/header.json", "★★", false},
	{"DPBMS白難易度表(通常)", "http://zris.work/bmstable/dp_white/header.json", "白", false},
	{"DPBMS黒難易度表(発狂)", "http://zris.work/bmstable/dp_black/header.json", "黒", false},
	{"発狂DPBMSごった煮難易度表", "http://zris.work/bmstable/dp_zhu/header.json", "★", false},
	{"発狂14keyBMS闇鍋難易度表", "http://zris.work/bmstable/dp_anguo/head14.json", "★", false},
	{"DPBMSと諸感", "http://zris.work/bmstable/dp_zhugan/header.json", "☆", false},
	{"Luminous", "http://zris.work/bmstable/luminous/header.json", "ln", false},
	{"LN難易度", "http://zris.work/bmstable/ln/ln_header.json", "◆", true},
	{"Scramble難易度表", "http://zris.work/bmstable/scramble/header.json", "SB", true},
	{"PMSデータベース(Lv1~45)", "http://zris.work/bmstable/pms_normal/pmsdatabase_header.json", "PLv", false},
	{"発狂PMSデータベース(lv46～)", "https://pmsdifficulty.xxxxxxxx.jp/insane_PMSdifficulty.html", "P●", false},
	{"発狂PMS難易度表", "http://zris.work/bmstable/pms_upper/header.json", "●", true},
	{"PMS Database コースデータ案内所", "http://zris.work/bmstable/pms_course/course_header.json", "Pcourse", true},
	{"Stellalite", "http://zris.work/bmstable/stellalite/Stellalite-header.json", "stl", false},
	{"オマージュBMS難易度表", "http://zris.work/bmstable/homage/header.json", "∽", false},
}

// This is a very stupid test: go through every table and see
// it smokes or not
func TestParseFromURL(t *testing.T) {
	for _, tt := range realTableDefintion {
		t.Run(tt.name, func(t *testing.T) {
			header, err := bmstable.ParseFromURL(tt.url)
			if err != nil {
				t.Fatalf("parse: %s", err)
			}
			if header.HeaderURL != tt.url {
				t.Fatalf("expect headerURL: %s, got %s", tt.url, header.HeaderURL)
			}
			if header.Symbol != tt.symbol {
				t.Fatalf("expect symbol: %s, got %s", tt.symbol, header.Symbol)
			}
			if header.Name != tt.name {
				t.Fatalf("expect name: %s, got %s", tt.name, header.Name)
			}
			if tt.hasCourses && len(header.Courses) == 0 {
				t.Fatalf("expect courses, got nothing")
			}
			if len(header.Contents) == 0 {
				t.Fatalf("expect contents, got nothing")
			}

			courseNameSet := make(map[string]any)
			for _, course := range header.Courses {
				if _, ok := courseNameSet[course.Name]; ok {
					t.Fatalf("course: %s is duplicated", course.Name)
				}
				courseNameSet[course.Name] = new(any)
			}
		})
	}
}
