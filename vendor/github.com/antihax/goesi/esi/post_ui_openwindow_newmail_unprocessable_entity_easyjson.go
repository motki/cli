// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package esi

import (
	json "encoding/json"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson7f10966bDecodeGithubComAntihaxGoesiEsi(in *jlexer.Lexer, out *PostUiOpenwindowNewmailUnprocessableEntityList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(PostUiOpenwindowNewmailUnprocessableEntityList, 0, 4)
			} else {
				*out = PostUiOpenwindowNewmailUnprocessableEntityList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 PostUiOpenwindowNewmailUnprocessableEntity
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson7f10966bEncodeGithubComAntihaxGoesiEsi(out *jwriter.Writer, in PostUiOpenwindowNewmailUnprocessableEntityList) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v PostUiOpenwindowNewmailUnprocessableEntityList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7f10966bEncodeGithubComAntihaxGoesiEsi(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostUiOpenwindowNewmailUnprocessableEntityList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7f10966bEncodeGithubComAntihaxGoesiEsi(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostUiOpenwindowNewmailUnprocessableEntityList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7f10966bDecodeGithubComAntihaxGoesiEsi(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostUiOpenwindowNewmailUnprocessableEntityList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7f10966bDecodeGithubComAntihaxGoesiEsi(l, v)
}
func easyjson7f10966bDecodeGithubComAntihaxGoesiEsi1(in *jlexer.Lexer, out *PostUiOpenwindowNewmailUnprocessableEntity) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "error":
			out.Error_ = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson7f10966bEncodeGithubComAntihaxGoesiEsi1(out *jwriter.Writer, in PostUiOpenwindowNewmailUnprocessableEntity) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Error_ != "" {
		const prefix string = ",\"error\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Error_))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostUiOpenwindowNewmailUnprocessableEntity) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7f10966bEncodeGithubComAntihaxGoesiEsi1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostUiOpenwindowNewmailUnprocessableEntity) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7f10966bEncodeGithubComAntihaxGoesiEsi1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostUiOpenwindowNewmailUnprocessableEntity) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7f10966bDecodeGithubComAntihaxGoesiEsi1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostUiOpenwindowNewmailUnprocessableEntity) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7f10966bDecodeGithubComAntihaxGoesiEsi1(l, v)
}
