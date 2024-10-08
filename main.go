package main

import (
	"errors"
	"github.com/gopherjs/gopherjs/js"
	"html/template"
	"reflect"
	"strings"
)

func throw(err error) {
	js.Global.Call("$throw", js.Global.Get("Error").New(err.Error()))
}

func typeof(something *js.Object) string {
	return js.Global.Call("$typeOf", something).String()
}

func templateFromJs(this *js.Object) *template.Template {
	return this.Get("_ptr").Interface().(wrappedTemplate).t
}

func templateToJs(t *template.Template, this *js.Object) {
	this.Set("_ptr", js.MakeWrapper(wrappedTemplate{t}))
}

func templateNew(this *js.Object, arguments []*js.Object) any {
	t := templateFromJs(this)
	if len(arguments) != 1 || typeof(arguments[0]) != "string" {
		throw(errors.New("invalid arguments"))
		return nil
	}
	this2 := js.Global.Get("Object").Call("create", js.Module.Get("Template").Get("prototype"))
	t2 := t.New(arguments[0].String())
	templateToJs(t2, this2)
	return this2
}

func templateFuncs(this *js.Object, arguments []*js.Object) any {
	t := templateFromJs(this)
	if len(arguments) != 1 || typeof(arguments[0]) != "object" {
		throw(errors.New("invalid arguments"))
		return nil
	}
	m, ok := arguments[0].Interface().(map[string]any)
	if !ok {
		throw(errors.New("invalid arguments"))
	}
	mfuncs := map[string]any{}
	for k, v := range m {
		mfuncs[k] = func(args ...any) (ret any, err error) {
			defer func() {
				e := recover()
				if e == nil {
					return
				}
				if e1, ok := e.(*js.Error); ok {
					err = e1
				} else {
					panic(e)
				}
			}()
			argsValue := make([]reflect.Value, len(args))
			for i, arg := range args {
				argsValue[i] = reflect.ValueOf(arg)
			}
			retSlice := reflect.ValueOf(v).Call(argsValue)
			if len(retSlice) != 1 {
				return nil, err
			}
			return retSlice[0].Interface(), nil
		}
	}
	t.Funcs(mfuncs)
	return this
}

func templateParse(this *js.Object, arguments []*js.Object) any {
	t := templateFromJs(this)
	if len(arguments) != 1 || typeof(arguments[0]) != "string" {
		throw(errors.New("invalid arguments"))
		return nil
	}
	_, err := t.Parse(arguments[0].String())
	if err != nil {
		throw(err)
		return nil
	}
	return this
}

func templateOption(this *js.Object, arguments []*js.Object) any {
	t := templateFromJs(this)
	options := make([]string, len(arguments))
	for i, arg := range arguments {
		if typeof(arguments[0]) != "string" {
			throw(errors.New("invalid arguments"))
		}
		options[i] = arg.String()
	}
	t.Option(options...)
	return this
}

func templateDefinedTemplates(this *js.Object, arguments []*js.Object) any {
	t := templateFromJs(this)
	if len(arguments) != 0 {
		throw(errors.New("invalid arguments"))
		return nil
	}
	return t.DefinedTemplates()
}

func templateExecute(this *js.Object, arguments []*js.Object) any {
	t := templateFromJs(this)
	if len(arguments) != 1 {
		throw(errors.New("invalid arguments"))
		return nil
	}
	var builder strings.Builder
	err := t.Execute(&builder, arguments[0].Interface())
	if err != nil {
		throw(err)
		return nil
	}

	return builder.String()
}

func main() {
	js.Module.Get("exports").Set("Template", js.MakeFunc(func(this *js.Object, arguments []*js.Object) any {
		if len(arguments) != 1 || typeof(arguments[0]) != "string" {
			throw(errors.New("invalid arguments"))
			return nil
		}
		templateToJs(template.New(arguments[0].String()), this)
		return this
	}))

	proto := js.Module.Get("exports").Get("Template").Get("prototype")
	proto.Set("option", js.MakeFunc(templateOption))
	proto.Set("definedTemplates", js.MakeFunc(templateDefinedTemplates))
	proto.Set("new", js.MakeFunc(templateNew))
	proto.Set("funcs", js.MakeFunc(templateFuncs))
	proto.Set("parse", js.MakeFunc(templateParse))
	proto.Set("execute", js.MakeFunc(templateExecute))
}

type wrappedTemplate struct {
	t *template.Template
}
