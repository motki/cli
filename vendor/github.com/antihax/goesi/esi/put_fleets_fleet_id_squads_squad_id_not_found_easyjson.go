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

func easyjson7e20b68eDecodeGithubComAntihaxGoesiEsi(in *jlexer.Lexer, out *PutFleetsFleetIdSquadsSquadIdNotFoundList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(PutFleetsFleetIdSquadsSquadIdNotFoundList, 0, 4)
			} else {
				*out = PutFleetsFleetIdSquadsSquadIdNotFoundList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 PutFleetsFleetIdSquadsSquadIdNotFound
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
func easyjson7e20b68eEncodeGithubComAntihaxGoesiEsi(out *jwriter.Writer, in PutFleetsFleetIdSquadsSquadIdNotFoundList) {
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
func (v PutFleetsFleetIdSquadsSquadIdNotFoundList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7e20b68eEncodeGithubComAntihaxGoesiEsi(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PutFleetsFleetIdSquadsSquadIdNotFoundList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7e20b68eEncodeGithubComAntihaxGoesiEsi(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PutFleetsFleetIdSquadsSquadIdNotFoundList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7e20b68eDecodeGithubComAntihaxGoesiEsi(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PutFleetsFleetIdSquadsSquadIdNotFoundList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7e20b68eDecodeGithubComAntihaxGoesiEsi(l, v)
}
func easyjson7e20b68eDecodeGithubComAntihaxGoesiEsi1(in *jlexer.Lexer, out *PutFleetsFleetIdSquadsSquadIdNotFound) {
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
func easyjson7e20b68eEncodeGithubComAntihaxGoesiEsi1(out *jwriter.Writer, in PutFleetsFleetIdSquadsSquadIdNotFound) {
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
func (v PutFleetsFleetIdSquadsSquadIdNotFound) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7e20b68eEncodeGithubComAntihaxGoesiEsi1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PutFleetsFleetIdSquadsSquadIdNotFound) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7e20b68eEncodeGithubComAntihaxGoesiEsi1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PutFleetsFleetIdSquadsSquadIdNotFound) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7e20b68eDecodeGithubComAntihaxGoesiEsi1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PutFleetsFleetIdSquadsSquadIdNotFound) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7e20b68eDecodeGithubComAntihaxGoesiEsi1(l, v)
}
