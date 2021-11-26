package iphostacl

const (
	LevelPackage uint = iota
	LevelUserProxy
	LevelUserBackconnect
	LevelUser
	LevelServerProxy
	LevelServerBackconnect
	LevelServer
	LevelGlobalProxy
	LevelGlobalBackconnect
	LevelGlobal
	LevelNotPresent
)

var strLevels = []string{
	"Package level",
	"User level - proxy only",
	"User level - backconnect only",
	"User level",
	"Server level - proxy only",
	"Server level - backconnect only",
	"Server level",
	"Global level - proxy only",
	"Global level - backconnect only",
	"Global level",
	"Not present in graylists",
}

func LevelToString(level uint) string {
	if level > LevelNotPresent {
		return ""
	}
	return strLevels[level]
}
