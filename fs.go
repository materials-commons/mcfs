package mcfs

import "path/filepath"

func EntryType(path string) PMEnum {
	entries := filepath.SplitList(path)
	entriesLen := len(entries)
	ext := filepath.Ext(path)
	switch {
	case len(entries) == 1:
		return PMProject
	case len(entries) == 2:
		return PMDir
	case len(entries) == 3:
		switch entries[2] {
		case "files", "json", "html", "text", "samples":
			return PMDir
		default:
			return PMProject
		}
	case ext != "":
		switch entries[entriesLen-2] {
		case "input_files", "output_files":
			return PMFile
		case "input_samples", "output_samples":
			return PMSample
		case "processes":
			return PMProcess
		case "properties":
			return PMProperty
		default:
			return PMNotFound
		}
	default:
		return PMNotFound
	}
}
