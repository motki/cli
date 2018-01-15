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

func easyjsonFeb661e9DecodeGithubComAntihaxGoesiEsi(in *jlexer.Lexer, out *GetCharactersCharacterIdStatsPveList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(GetCharactersCharacterIdStatsPveList, 0, 2)
			} else {
				*out = GetCharactersCharacterIdStatsPveList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 GetCharactersCharacterIdStatsPve
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
func easyjsonFeb661e9EncodeGithubComAntihaxGoesiEsi(out *jwriter.Writer, in GetCharactersCharacterIdStatsPveList) {
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
func (v GetCharactersCharacterIdStatsPveList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonFeb661e9EncodeGithubComAntihaxGoesiEsi(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetCharactersCharacterIdStatsPveList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonFeb661e9EncodeGithubComAntihaxGoesiEsi(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetCharactersCharacterIdStatsPveList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonFeb661e9DecodeGithubComAntihaxGoesiEsi(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetCharactersCharacterIdStatsPveList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonFeb661e9DecodeGithubComAntihaxGoesiEsi(l, v)
}
func easyjsonFeb661e9DecodeGithubComAntihaxGoesiEsi1(in *jlexer.Lexer, out *GetCharactersCharacterIdStatsPve) {
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
		case "dungeons_completed_agent":
			out.DungeonsCompletedAgent = int64(in.Int64())
		case "dungeons_completed_distribution":
			out.DungeonsCompletedDistribution = int64(in.Int64())
		case "missions_succeeded":
			out.MissionsSucceeded = int64(in.Int64())
		case "missions_succeeded_epic_arc":
			out.MissionsSucceededEpicArc = int64(in.Int64())
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
func easyjsonFeb661e9EncodeGithubComAntihaxGoesiEsi1(out *jwriter.Writer, in GetCharactersCharacterIdStatsPve) {
	out.RawByte('{')
	first := true
	_ = first
	if in.DungeonsCompletedAgent != 0 {
		const prefix string = ",\"dungeons_completed_agent\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.DungeonsCompletedAgent))
	}
	if in.DungeonsCompletedDistribution != 0 {
		const prefix string = ",\"dungeons_completed_distribution\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.DungeonsCompletedDistribution))
	}
	if in.MissionsSucceeded != 0 {
		const prefix string = ",\"missions_succeeded\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.MissionsSucceeded))
	}
	if in.MissionsSucceededEpicArc != 0 {
		const prefix string = ",\"missions_succeeded_epic_arc\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.MissionsSucceededEpicArc))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetCharactersCharacterIdStatsPve) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonFeb661e9EncodeGithubComAntihaxGoesiEsi1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetCharactersCharacterIdStatsPve) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonFeb661e9EncodeGithubComAntihaxGoesiEsi1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetCharactersCharacterIdStatsPve) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonFeb661e9DecodeGithubComAntihaxGoesiEsi1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetCharactersCharacterIdStatsPve) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonFeb661e9DecodeGithubComAntihaxGoesiEsi1(l, v)
}
