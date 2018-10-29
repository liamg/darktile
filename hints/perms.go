package hints

import (
	"fmt"
	"regexp"
	"strings"
)

type perms struct {
	IsDirectory bool
	Owner       access
	Group       access
	World       access
}

type access struct {
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
	p.IsDirectory = s[0] == 'd'

	p.Owner.Read = s[1] == 'r'
	p.Owner.Write = s[2] == 'w'
	p.Owner.Execute = s[3] == 'x'
	p.Group.Read = s[4] == 'r'
	p.Group.Write = s[5] == 'w'
	p.Group.Execute = s[6] == 'x'
	p.World.Read = s[7] == 'r'
	p.World.Write = s[8] == 'w'
	p.World.Execute = s[9] == 'x'

	return p, nil
}

func init() {
	hinters = append(hinters, hintPerms)
}

func hintPerms(word string, context string, wordX uint16, wordY uint16) *Hint {

	item := &Hint{
		Line:   context,
		Word:   word,
		StartX: wordX,
		StartY: wordY,
	}

	if wordX == 0 {

		p, err := parsePermissionString(word)
		if err != nil {
			return nil
		}

		typ := "file"
		if p.IsDirectory {
			typ = "directory"
		}

		item.Description = fmt.Sprintf(`Permissions: 
  Type:    %s
  Numeric: %s
  Owner:   %s
  Group:   %s
  World:   %s`,
			typ,
			p.Numeric(),
			p.Owner.Nice(),
			p.Group.Nice(),
			p.World.Nice(),
		)

		return item
	}

	return nil
}

func isPermString(s string) bool {
	re := regexp.MustCompile("[dl\\-sS]{1}[sSrwx\\-]{9}")
	return re.MatchString(s)
}
