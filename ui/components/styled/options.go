package styled

type Options uint

const (
	Primary Options = 1 << iota
	Secondary
	Danger
	Warn
	Success

	Text
	Bordered

	Sm
	Lg
)

func (o Options) Primary() Options   { return o | Primary }
func (o Options) Secondary() Options { return o | Secondary }
func (o Options) Danger() Options    { return o | Danger }
func (o Options) Warn() Options      { return o | Warn }
func (o Options) Success() Options   { return o | Success }
func (o Options) Text() Options      { return o | Text }
func (o Options) Bordered() Options  { return o | Bordered }
func (o Options) Sm() Options        { return o | Sm }
func (o Options) Lg() Options        { return o | Lg }

func (o Options) Has(opt Options) bool {
	return o&opt != 0
}
