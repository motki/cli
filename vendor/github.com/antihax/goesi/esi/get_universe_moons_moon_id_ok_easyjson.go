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

func easyjsonE06267a0DecodeGithubComAntihaxGoesiEsi(in *jlexer.Lexer, out *GetUniverseMoonsMoonIdOkList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(GetUniverseMoonsMoonIdOkList, 0, 1)
			} else {
				*out = GetUniverseMoonsMoonIdOkList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 GetUniverseMoonsMoonIdOk
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
func easyjsonE06267a0EncodeGithubComAntihaxGoesiEsi(out *jwriter.Writer, in GetUniverseMoonsMoonIdOkList) {
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
func (v GetUniverseMoonsMoonIdOkList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE06267a0EncodeGithubComAntihaxGoesiEsi(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetUniverseMoonsMoonIdOkList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE06267a0EncodeGithubComAntihaxGoesiEsi(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetUniverseMoonsMoonIdOkList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE06267a0DecodeGithubComAntihaxGoesiEsi(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetUniverseMoonsMoonIdOkList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE06267a0DecodeGithubComAntihaxGoesiEsi(l, v)
}
func easyjsonE06267a0DecodeGithubComAntihaxGoesiEsi1(in *jlexer.Lexer, out *GetUniverseMoonsMoonIdOk) {
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
		case "moon_id":
			out.MoonId = int32(in.Int32())
		case "name":
			out.Name = string(in.String())
		case "position":
			(out.Position).UnmarshalEasyJSON(in)
		case "system_id":
			out.SystemId = int32(in.Int32())
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
func easyjsonE06267a0EncodeGithubComAntihaxGoesiEsi1(out *jwriter.Writer, in GetUniverseMoonsMoonIdOk) {
	out.RawByte('{')
	first := true
	_ = first
	if in.MoonId != 0 {
		const prefix string = ",\"moon_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.MoonId))
	}
	if in.Name != "" {
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	if true {
		const prefix string = ",\"position\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Position).MarshalEasyJSON(out)
	}
	if in.SystemId != 0 {
		const prefix string = ",\"system_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.SystemId))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetUniverseMoonsMoonIdOk) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE06267a0EncodeGithubComAntihaxGoesiEsi1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetUniverseMoonsMoonIdOk) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE06267a0EncodeGithubComAntihaxGoesiEsi1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetUniverseMoonsMoonIdOk) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE06267a0DecodeGithubComAntihaxGoesiEsi1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetUniverseMoonsMoonIdOk) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE06267a0DecodeGithubComAntihaxGoesiEsi1(l, v)
}
