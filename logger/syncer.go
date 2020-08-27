package logger

// Syncer defines how to write log message
type Syncer string

const (
	// Std syncer writes to stdout
	Std Syncer = "std"
	// File syncer writes to file
	File Syncer = "file"
)

// Enabled function determines exact syncer
func (s Syncer) Enabled(syncer Syncer) bool {
	return syncer == s
}
