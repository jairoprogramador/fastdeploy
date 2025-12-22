package vos

import "time"

// Commit representa la información esencial de un commit de Git.
type Commit struct {
	// Hash es el SHA completo del commit.
	Hash string
	// Message es el mensaje de commit completo.
	Message string
	// Author es el autor del commit.
	Author string
	// Date es la fecha en que se realizó el commit.
	Date time.Time
}

func (c *Commit) String() string {
	return c.Hash
}