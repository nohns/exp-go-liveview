package glive

import "github.com/a-h/templ"

type View interface {
	Render() templ.Component
	ID() string
}

type ViewLifecycleMount interface {
	OnMount(*Session)
}

type ViewLifecycleUnmount interface {
	OnUnmount()
}

type ViewBuilder func() View
