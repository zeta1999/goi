package goi

import (
	dom "github.com/TobiasYin/goi/vdom"
	"strings"
)

type Widget interface {
	Pack(ctx Context) Node
}

type Node interface {
	pack() dom.JsDomElement
	getContext() Context
}

type packAble interface {
	getParam() Params
	packChildAble
}

type packChildAble interface {
	getChildren() []Widget
}

func packChildren(able packChildAble, ele *dom.JsDomElement, ctx Context) {
	for _, c := range able.getChildren() {
		if c == nil{
			continue
		}
		//Pack
		v := c.Pack(ctx).pack()
		ele.AppendChild(v)
	}
}

func pack(able packAble, name string, ctx Context) dom.JsDomElement {
	p := able.getParam()
	e := p.packWithName(name)
	packChildren(able, &e, ctx)
	return e
}

type Params struct {
	Style       Style
	Class       []string
	Value       string
	Placeholder string
	OnClick     EventCallBack
	OnInput     EventCallBack
	OnChange    EventCallBack
	OnFocus     EventCallBack
	OnBlur      EventCallBack
	OnKeyDown   EventCallBack
	OnKeyUp     EventCallBack
}

func (e Params) getParam() Params {
	return e
}

func (e Params) packWithName(name string) dom.JsDomElement {
	ele := dom.Dom.CreateElement(name)
	style := e.Style.packStyle()
	if style != "" {
		ele.Set("style", style)
	}
	if len(e.Class) != 0 {
		ele.Set("className", strings.Join(e.Class, " "))
	}
	if e.OnClick != nil {
		ele.Set("onclick", WrapEventCallBack(e.OnClick))
	}
	if e.OnInput != nil {
		ele.Set("oninput", WrapEventCallBack(e.OnInput))
	}
	if e.OnChange != nil {
		ele.Set("onchange", WrapEventCallBack(e.OnChange))
	}
	if e.OnFocus != nil {
		ele.Set("onfocus", WrapEventCallBack(e.OnFocus))
	}
	if e.OnBlur != nil {
		ele.Set("onblur", WrapEventCallBack(e.OnBlur))
	}
	if e.OnKeyDown != nil {
		ele.Set("onkeydown", WrapEventCallBack(e.OnKeyDown))
	}
	if e.OnKeyUp != nil {
		ele.Set("onkeyup", WrapEventCallBack(e.OnKeyUp))
	}
	return ele
}
