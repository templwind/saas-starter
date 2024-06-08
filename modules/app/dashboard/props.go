package dashboard

import (
	"net/http"

	"github.com/templwind/sass-starter/internal/config"

	"github.com/a-h/templ"
	"github.com/templwind/templwind"
)

type Props struct {
	Request *http.Request
	Config  *config.Config
	ID      string
}

// New creates a new component
func New(props ...templwind.OptFunc[Props]) templ.Component {
	return templwind.New(defaultProps, tpl, props...)
}

// NewWithProps creates a new component with the given props
func NewWithProps(props *Props) templ.Component {
	return templwind.NewWithProps(tpl, props)
}

// WithProps builds the propsions with the given props
func WithProps(props ...templwind.OptFunc[Props]) *Props {
	return templwind.WithProps(defaultProps, props...)
}

func defaultProps() *Props {
	return &Props{}
}

func WithRequest(r *http.Request) templwind.OptFunc[Props] {
	return func(o *Props) {
		o.Request = r
	}
}

func WithConfig(cfg *config.Config) templwind.OptFunc[Props] {
	return func(o *Props) {
		o.Config = cfg
	}
}

func WithID(id string) templwind.OptFunc[Props] {
	return func(o *Props) {
		o.ID = id
	}
}
