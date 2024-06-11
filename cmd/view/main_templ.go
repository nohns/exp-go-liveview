// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.719
package view

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"

import "github.com/nohns/go-liveview"
import "github.com/nohns/go-liveview/internal/rng"
import "strconv"
import "time"
import "fmt"

type main struct {
	id  string
	ses *glive.Session
	t   *time.Ticker

	count int
	user  string
}

func Main(user string) glive.ViewBuilder {
	return glive.ViewBuilder(func() glive.View {
		id, _ := rng.ViewID()
		return &main{id: id, user: user}
	})
}

func (v *main) OnMount(ses *glive.Session) {
	v.ses = ses

	v.t = time.NewTicker(1 * time.Second)
	go func() {
		for range v.t.C {
			v.count++
			fmt.Printf("count now %d\n", v.count)
			ses.Rerender(v)
		}
	}()
}

func (v *main) OnUnmount() {
	v.t.Stop()
}

func (v *main) ID() string {
	return v.id
}

func (v *main) RootAttr() templ.Attributes {
	return map[string]any{
		"data-glive-view": v.id,
	}
}

func (v *main) Render() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*templ.DiffBuffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<main")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderAttributes(ctx, templ_7745c5c3_Buffer, v.RootAttr())
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(">Hello ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(v.user)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `cmd/view/main.templ`, Line: 54, Col: 16}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteDynamic(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" count is now ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(v.count))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `cmd/view/main.templ`, Line: 54, Col: 55}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteDynamic(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("!</main>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
