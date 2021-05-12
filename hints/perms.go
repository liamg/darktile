package hints

import (
	"fmt"
	"regexp"
	"strings"
)

type fileType uint8

const (
	file fileType = iota
	directory
	characterSpecialFile
)

type perms struct {
	Type   fileType
	Owner  access
	Group  access
	World  access
	SetUID bool
	SetGID bool
	Sticky bool
}

type accessType uint8

const (
	owner accessType = iota
	group
	world
)

type access struct {
	Type    accessType
	Read    bool
	Write   bool
	Execute bool
}

func (p perms) Numeric() string {
	return p.Owner.Numeric() + p.Group.Numeric() + p.World.Numeric()
}

func (a access) Nice() string {
	all := []string{}
	if a.Read {
		all = append(all, "read")
	}
	if a.Write {
		all = append(all, "write")
	}
	if a.Execute {
		all = append(all, "execute")
	}

	return strings.Join(all, ", ")
}

func (a access) Numeric() string {
	var n uint8
	if a.Read {
		n += 4
	}
	if a.Write {
		n += 2
	}
	if a.Execute {
		n++
	}
	return fmt.Sprintf("%d", n)
}

func parsePermissionString(s string) (perms, error) {
	if !isPermString(s) {
		return perms{}, fmt.Errorf("Invalid permission string")
	}
	p := perms{}
	switch s[0] {
	case 'c':
		p.Type = characterSpecialFile
	case 'd':
		p.Type = directory
	default:
		p.Type = file
	}

	p.SetUID = s[3] == 's' || s[3] == 'S'
	p.SetGID = s[6] == 's' || s[6] == 'S'
	p.Sticky = s[9] == 't' || s[9] == 'T'

	p.Owner.Type = owner
	p.Owner.Read = s[1] == 'r'
	p.Owner.Write = s[2] == 'w'
	p.Owner.Execute = s[3] == 'x' || s[3] == 's'
	p.Group.Type = group
	p.Group.Read = s[4] == 'r'
	p.Group.Write = s[5] == 'w'
	p.Group.Execute = s[6] == 'x' || s[6] == 's'
	p.World.Type = world
	p.World.Read = s[7] == 'r'
	p.World.Write = s[8] == 'w'
	p.World.Execute = s[9] == 'x'

	return p, nil
}

func init() {
	hinters = append(hinters, hintPerms)
}

func hintPerms(word string, context string, wordX uint16, wordY uint16) *Hint {
	item := NewHint(word, context, wordX, wordY)

	if wordX == 0 {

		p, err := parsePermissionString(word)
		if err != nil {
			return nil
		}

		typ := "file"
		switch p.Type {
		case directory:
			typ = "directory"
		case characterSpecialFile:
			typ = "character special file"
		}

		item.Description = fmt.Sprintf(`Permissions: 
  Type:    %s
  Numeric: %s
  Owner:   %s
  Group:   %s
  World:   %s
  Setuid:  %t
  Setgid:  %t
  Sticky:  %t
  `,
			typ,
			p.Numeric(),
			p.Owner.Nice(),
			p.Group.Nice(),
			p.World.Nice(),
			p.SetUID,
			p.SetGID,
			p.Sticky,
		)

		return item
	}

	return nil
}

func isPermString(s string) bool {
	re := regexp.MustCompile("[cdl\\-sS]{1}[sStTrwx\\-]{9}")
	return re.MatchString(s)
}
