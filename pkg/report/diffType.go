package report

type DiffType int

const (
	Added DiffType = iota
	Removed
	TypeChanged
	ValueChanged
)

func (d DiffType) String() string {
	return [...]string{"Added", "Removed", "TypeChanged", "ValueChanged"}[d]
}