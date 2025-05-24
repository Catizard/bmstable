package internal

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type ImportHeaderVo struct {
	DataURL     string `json:"data_url"`
	Name        string `json:"name"`
	OriginalURL string `json:"original_url"`
	Symbol      string `json:"symbol"`
	HeaderURL   string
	LevelOrders []any `json:"level_order"`
	RawCourses  []any `json:"course"`
	Courses     []ImportCourseVo
}

type ImportCourseVo struct {
	Name       string        `json:"name"`
	Md5        []string      `json:"md5"`
	Sha256     []string      `json:"sha256"`
	Constraint []string      `json:"constraint"`
	Charts     []ChartInfoVo `json:"charts"`
	HeaderID   uint
}

type ChartInfoVo struct {
	Title    string
	SubTitle string
	Artist   string
	Sha256   string `json:"sha256"`
	Md5      string `json:"md5"`
}

// Parse field 'RawCourses' into 'Courses'
//
// Possible layouts:
//  1. courses is an array of valid courses
//  2. courses is a two-dimensional array, every element might be an array of valid courses
//  3. courses is an array of a wrapped struct, the real courses are laid inside 'charts' field
//
// For Item3, see pushupChartsHashField for details
func (header *ImportHeaderVo) ParseRawCourses() error {
	if len(header.RawCourses) == 0 {
		return nil // Okay dokey
	}
	if _, isNested := header.RawCourses[0].([]any); isNested {
		for i := range header.RawCourses {
			innerArray := header.RawCourses[i].([]any)
			for _, data := range innerArray {
				courseInfo := ImportCourseVo{}
				if err := mapstructure.Decode(data, &courseInfo); err != nil {
					return err
				}
				if err := courseInfo.pushupChartsHashField(); err != nil {
					return fmt.Errorf("course: %s", err)
				}
				header.Courses = append(header.Courses, courseInfo)
			}
		}
	} else {
		for _, data := range header.RawCourses {
			courseInfo := ImportCourseVo{}
			if err := mapstructure.Decode(data, &courseInfo); err != nil {
				return err
			}
			if err := courseInfo.pushupChartsHashField(); err != nil {
				return fmt.Errorf("course: %s", err)
			}
			header.Courses = append(header.Courses, courseInfo)
		}
	}

	return nil
}

// Some tables' courses are defined in an inner field `charts`, this function is 'pushing' them up
func (courseInfo *ImportCourseVo) pushupChartsHashField() error {
	if len(courseInfo.Charts) > 0 {
		// `charts` may provide `sha256` or `md5`
		firstChartDef := courseInfo.Charts[0]
		if firstChartDef.Md5 != "" {
			courseInfo.Md5 = make([]string, 0)
			for _, chart := range courseInfo.Charts {
				courseInfo.Md5 = append(courseInfo.Md5, chart.Md5)
			}
		} else if firstChartDef.Sha256 != "" {
			courseInfo.Sha256 = make([]string, 0)
			for _, chart := range courseInfo.Charts {
				courseInfo.Sha256 = append(courseInfo.Sha256, chart.Sha256)
			}
		} else {
			return fmt.Errorf("course: no sha256 or md5 provides")
		}
	}
	return nil
}
