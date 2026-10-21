package helpers

import "fmt"

// formatFileSize converts byte counts into a more human-readable format such as KB or MB.
func (f *FileStats) formatFileSize(bytes int64) string {
	const unit = 1024

	if bytes < unit {
		return fmt.Sprintf("%d bytes", bytes)
	}

	// Determine the appropriate unit for the size.
	units := []string{"bytes", "kb", "MB", "GB", "TB"}
	exp := 0
	size := float64(bytes)

	for size >= unit && exp < len(units)-1 {
		size /= unit
		exp++
	}

	// If the size is still in bytes, show it without decimals.
	if exp == 0 {
		return fmt.Sprintf("%d %s", int64(size), units[exp])
	}

	// For larger units, show two decimal places for readability.
	return fmt.Sprintf("%.2f %s", size, units[exp])
}
