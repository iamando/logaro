package logaro

import (
	"bytes"
	"encoding/json"
)

// isEnabled checks if the given log level is enabled based on the logger's configured level.
// It uses a map to associate the log levels with numeric values.
// The function compares the numeric log levels of the given level and the logger's level.
// Returns true if the given level is enabled (its numeric value is greater than or equal to
// the logger's numeric level value), false otherwise.
// The function allows determining if a log entry with a specific level should be logged
// based on the logger's configured log level.
func (l *Logger) isEnabled(level string) bool {
	levels := map[string]int{
		"fatal": 5,
		"error": 4,
		"warn":  3,
		"info":  2,
		"debug": 1,
	}

	return levels[level] >= levels[l.Level]
}

// mergeFields merges the event fields from the parent logger with the current logger's event fields
// and the additional fields provided as a parameter.
// It creates a new map to hold the merged fields and copies the parent's event fields into it.
// Then it adds the current logger's event fields and the additional fields to the merged map.
// Returns the merged map of event fields, combining the inherited fields from the parent logger
// with the current logger's fields and the additional fields provided as a parameter.
// The function is used to create a consolidated set of event fields for log entries,
// ensuring that all relevant fields are included in the log context.
func (l *Logger) mergeFields(fields map[string]interface{}) map[string]interface{} {
	mergedFields := make(map[string]interface{})

	if l.Parent != nil {
		mergedFields = l.Parent.mergeFields(l.Parent.EventFields)
	}

	for key, val := range l.EventFields {
		mergedFields[key] = val
	}

	for key, val := range fields {
		mergedFields[key] = val
	}

	return mergedFields
}

// serializeEntry applies the logger's serializer function to the log entry.
// If a serializer is set for the logger, it applies the serializer function to the log message
// and the fields of the entry, allowing custom modification or formatting of the log entry.
// Returns the serialized log entry with the log message and fields modified by the serializer,
// ensuring that the log data is transformed according to the specified serialization logic.
// The function is used to customize the serialization process for specific log entries
// based on the serializer function set for the logger.
func (l *Logger) serializeEntry(entry LogEntry) LogEntry {
	if l.Serializer != nil {
		entry.Message = l.Serializer(entry.Message).(string)
		entry.Fields = l.Serializer(entry.Fields).(map[string]interface{})
	}

	return entry
}

func compareLogEntries(a, b LogEntry) bool {
	// Compare Timestamp, Message, Level, and Fields
	return a.Timestamp == b.Timestamp &&
		a.Message == b.Message &&
		a.Level == b.Level &&
		compareFields(a.Fields, b.Fields)
}

func compareFields(a, b map[string]interface{}) bool {
	// Compare the lengths of the fields maps
	if len(a) != len(b) {
		return false
	}

	// Compare each key-value pair in the fields maps
	for key, valA := range a {
		valB, ok := b[key]
		if !ok || !compareFieldValues(valA, valB) {
			return false
		}
	}

	return true
}

func compareFieldValues(a, b interface{}) bool {
	// Marshal and compare the JSON representations of the field values
	bytesA, errA := json.Marshal(a)
	bytesB, errB := json.Marshal(b)
	if errA != nil || errB != nil {
		return false
	}

	return bytes.Equal(bytesA, bytesB)
}
