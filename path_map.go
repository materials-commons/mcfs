package mcfs

type PMEnum int

const (
	PMProject PMEnum = iota
	PMSample
	PMFile
	PMProcess
	PMProperty
	PMDir
	PMNotFound
)

var pmEnumToStr = map[PMEnum]string{
	PMProject: "PMProject",
	PMSample: "PMSample",
	PMFile: "PMFile",
	PMProcess: "PMProcess",
	PMDir: "PMDir",
	PMNotFound: "PMNotFound",
}

func (p PMEnum) String() string {
	if val, found := pmEnumToStr[p]; found {
		return val
	}

	return "UNKNOWN"
}

type PathMapEntry struct {
	Type PMEnum
	ID   string
}

var pathMap map[string]PathMapEntry = make(map[string]PathMapEntry)

func AddProject(path string, ID string) {
	//fmt.Printf("AddProject '%s'\n", path)
	pathMap[path] = PathMapEntry{Type: PMProject, ID: ID}
}

func AddSample(path string, ID string) {
	pathMap[path] = PathMapEntry{Type: PMSample, ID: ID}
}

func AddFile(path string, ID string) {
	pathMap[path] = PathMapEntry{Type: PMFile, ID: ID}
}

func AddProcess(path string, ID string) {
	pathMap[path] = PathMapEntry{Type: PMProcess, ID: ID}
}

func AddProperty(path string, ID string) {
	pathMap[path] = PathMapEntry{Type: PMProperty, ID: ID}
}

func PathEntry(path string) (string, PMEnum) {
	//fmt.Printf("PathEntry '%s'\n", path)
	entry, found := pathMap[path]
	//fmt.Println("  ", entry, found)
	if found {
		return entry.ID, entry.Type
	}

	return "", PMNotFound
}
