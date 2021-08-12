package main

import "strings"

const PathSeparator = "/"

type Path struct {
	Path string
	ID   string
}

func NewPath(p string) *Path {
	var id string
	p = strings.Trim(p, PathSeparator)
	s := strings.Split(p, PathSeparator)
	if len(s) > 1 {
		id = s[len(s)-1] // Last string is id.
		p = strings.Join(s[:len(s)-1], PathSeparator)
	}
	return &Path{Path: p, ID: id}
}

// Return true if the path has an ID.
func (p *Path) HasID() bool {
	return len(p.ID) > 0
}
