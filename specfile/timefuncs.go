package specfile

import (
	"time"
)

type timeFuncs struct{}

func (t *timeFuncs) Now() time.Time {
	return time.Now()
}

func (t *timeFuncs) Parse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

func (t *timeFuncs) Format(layout string, in time.Time) string {
	return in.Format(layout)
}

// IsAfter returns true if in is after cmp.
//
// The argument order is chosen to make it usable as a template function where
// pipelines send the previous result as the last argument.
func (t *timeFuncs) IsAfter(cmp, in time.Time) bool {
	return in.After(cmp)
}

// IsBefore returns true if in is before cmp.
//
// The argument order is chosen to make it usable as a template function where
// pipelines send the previous result as the last argument.
func (t *timeFuncs) IsBefore(cmp, in time.Time) bool {
	return in.Before(cmp)
}

// IsEqual returns true if in is equal to cmp.
//
// The argument order is chosen to make it usable as a template function where
// pipelines send the previous result as the last argument.
func (t *timeFuncs) IsEqual(cmp, in time.Time) bool {
	return in.Equal(cmp)
}
