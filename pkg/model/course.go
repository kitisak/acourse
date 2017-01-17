package model

import (
	"time"
)

// Course model
type Course struct {
	Base
	Stampable
	Title            string `datastore:",noindex"`
	ShortDescription string `datastore:",noindex"`
	Description      string `datastore:",noindex"` // Markdown
	Photo            string `datastore:",noindex"` // URL
	Owner            string
	Start            time.Time
	URL              string
	Type             CourseType
	Video            string `datastore:",noindex"` // Cover Video
	Price            float64
	DiscountedPrice  float64
	Options          CourseOption
	Contents         CourseContents    `datastore:",noindex"`
	EnrollDetail     string            `datastore:",noindex"`
	Assignments      CourseAssignments `datastore:",noindex"`
}

// Courses model
type Courses []*Course

// CourseOption type
type CourseOption struct {
	Public     bool
	Enroll     bool `datastore:",noindex"`
	Attend     bool `datastore:",noindex"`
	Assignment bool `datastore:",noindex"`
	Discount   bool
}

// CourseContent type
type CourseContent struct {
	Title       string `datastore:",noindex"`
	Description string `datastore:",noindex"` // Markdown
	Video       string `datastore:",noindex"` // Youtube ID
	DownloadURL string `datastore:",noindex"` // Video download link
}

// CourseContents type
type CourseContents []CourseContent

// CourseType type
type CourseType string

// CourseType
const (
	CourseTypeLive  CourseType = "live"
	CourseTypeVideo CourseType = "video"
	CourseTypeEbook CourseType = "ebook"
)

// CourseAssignment type
type CourseAssignment struct {
	Title       string `datastore:",noindex"`
	Description string `datastore:",noindex"`
	Open        bool   `datastore:",noindex"`
}

// CourseAssignments type
type CourseAssignments []CourseAssignment
