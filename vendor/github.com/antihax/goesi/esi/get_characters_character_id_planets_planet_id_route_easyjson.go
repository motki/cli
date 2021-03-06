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

func easyjson5a671d6eDecodeGithubComAntihaxGoesiEsi(in *jlexer.Lexer, out *GetCharactersCharacterIdPlanetsPlanetIdRouteList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(GetCharactersCharacterIdPlanetsPlanetIdRouteList, 0, 1)
			} else {
				*out = GetCharactersCharacterIdPlanetsPlanetIdRouteList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 GetCharactersCharacterIdPlanetsPlanetIdRoute
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
func easyjson5a671d6eEncodeGithubComAntihaxGoesiEsi(out *jwriter.Writer, in GetCharactersCharacterIdPlanetsPlanetIdRouteList) {
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
func (v GetCharactersCharacterIdPlanetsPlanetIdRouteList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a671d6eEncodeGithubComAntihaxGoesiEsi(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetCharactersCharacterIdPlanetsPlanetIdRouteList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a671d6eEncodeGithubComAntihaxGoesiEsi(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetCharactersCharacterIdPlanetsPlanetIdRouteList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a671d6eDecodeGithubComAntihaxGoesiEsi(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetCharactersCharacterIdPlanetsPlanetIdRouteList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a671d6eDecodeGithubComAntihaxGoesiEsi(l, v)
}
func easyjson5a671d6eDecodeGithubComAntihaxGoesiEsi1(in *jlexer.Lexer, out *GetCharactersCharacterIdPlanetsPlanetIdRoute) {
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
		case "route_id":
			out.RouteId = int64(in.Int64())
		case "source_pin_id":
			out.SourcePinId = int64(in.Int64())
		case "destination_pin_id":
			out.DestinationPinId = int64(in.Int64())
		case "content_type_id":
			out.ContentTypeId = int32(in.Int32())
		case "quantity":
			out.Quantity = float32(in.Float32())
		case "waypoints":
			if in.IsNull() {
				in.Skip()
				out.Waypoints = nil
			} else {
				in.Delim('[')
				if out.Waypoints == nil {
					if !in.IsDelim(']') {
						out.Waypoints = make([]int64, 0, 8)
					} else {
						out.Waypoints = []int64{}
					}
				} else {
					out.Waypoints = (out.Waypoints)[:0]
				}
				for !in.IsDelim(']') {
					var v4 int64
					v4 = int64(in.Int64())
					out.Waypoints = append(out.Waypoints, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson5a671d6eEncodeGithubComAntihaxGoesiEsi1(out *jwriter.Writer, in GetCharactersCharacterIdPlanetsPlanetIdRoute) {
	out.RawByte('{')
	first := true
	_ = first
	if in.RouteId != 0 {
		const prefix string = ",\"route_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.RouteId))
	}
	if in.SourcePinId != 0 {
		const prefix string = ",\"source_pin_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.SourcePinId))
	}
	if in.DestinationPinId != 0 {
		const prefix string = ",\"destination_pin_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.DestinationPinId))
	}
	if in.ContentTypeId != 0 {
		const prefix string = ",\"content_type_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.ContentTypeId))
	}
	if in.Quantity != 0 {
		const prefix string = ",\"quantity\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Float32(float32(in.Quantity))
	}
	if len(in.Waypoints) != 0 {
		const prefix string = ",\"waypoints\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v5, v6 := range in.Waypoints {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.Int64(int64(v6))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetCharactersCharacterIdPlanetsPlanetIdRoute) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a671d6eEncodeGithubComAntihaxGoesiEsi1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetCharactersCharacterIdPlanetsPlanetIdRoute) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a671d6eEncodeGithubComAntihaxGoesiEsi1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetCharactersCharacterIdPlanetsPlanetIdRoute) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a671d6eDecodeGithubComAntihaxGoesiEsi1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetCharactersCharacterIdPlanetsPlanetIdRoute) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a671d6eDecodeGithubComAntihaxGoesiEsi1(l, v)
}
