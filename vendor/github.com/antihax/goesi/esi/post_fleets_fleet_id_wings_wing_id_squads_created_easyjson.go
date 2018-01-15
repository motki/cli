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

func easyjson8e08947bDecodeGithubComAntihaxGoesiEsi(in *jlexer.Lexer, out *PostFleetsFleetIdWingsWingIdSquadsCreatedList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(PostFleetsFleetIdWingsWingIdSquadsCreatedList, 0, 8)
			} else {
				*out = PostFleetsFleetIdWingsWingIdSquadsCreatedList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 PostFleetsFleetIdWingsWingIdSquadsCreated
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
func easyjson8e08947bEncodeGithubComAntihaxGoesiEsi(out *jwriter.Writer, in PostFleetsFleetIdWingsWingIdSquadsCreatedList) {
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
func (v PostFleetsFleetIdWingsWingIdSquadsCreatedList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8e08947bEncodeGithubComAntihaxGoesiEsi(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostFleetsFleetIdWingsWingIdSquadsCreatedList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8e08947bEncodeGithubComAntihaxGoesiEsi(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostFleetsFleetIdWingsWingIdSquadsCreatedList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8e08947bDecodeGithubComAntihaxGoesiEsi(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostFleetsFleetIdWingsWingIdSquadsCreatedList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8e08947bDecodeGithubComAntihaxGoesiEsi(l, v)
}
func easyjson8e08947bDecodeGithubComAntihaxGoesiEsi1(in *jlexer.Lexer, out *PostFleetsFleetIdWingsWingIdSquadsCreated) {
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
		case "squad_id":
			out.SquadId = int64(in.Int64())
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
func easyjson8e08947bEncodeGithubComAntihaxGoesiEsi1(out *jwriter.Writer, in PostFleetsFleetIdWingsWingIdSquadsCreated) {
	out.RawByte('{')
	first := true
	_ = first
	if in.SquadId != 0 {
		const prefix string = ",\"squad_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.SquadId))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostFleetsFleetIdWingsWingIdSquadsCreated) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8e08947bEncodeGithubComAntihaxGoesiEsi1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostFleetsFleetIdWingsWingIdSquadsCreated) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8e08947bEncodeGithubComAntihaxGoesiEsi1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostFleetsFleetIdWingsWingIdSquadsCreated) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8e08947bDecodeGithubComAntihaxGoesiEsi1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostFleetsFleetIdWingsWingIdSquadsCreated) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8e08947bDecodeGithubComAntihaxGoesiEsi1(l, v)
}
