package model

type Filter struct {
	Limit           int
	Page            int
	IsEncrypted     bool
	GroupDuplicates bool
}
