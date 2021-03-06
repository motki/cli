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

func easyjsonB5ebbfcbDecodeGithubComAntihaxGoesiEsi(in *jlexer.Lexer, out *GetKillmailsKillmailIdKillmailHashOkList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(GetKillmailsKillmailIdKillmailHashOkList, 0, 1)
			} else {
				*out = GetKillmailsKillmailIdKillmailHashOkList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 GetKillmailsKillmailIdKillmailHashOk
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
func easyjsonB5ebbfcbEncodeGithubComAntihaxGoesiEsi(out *jwriter.Writer, in GetKillmailsKillmailIdKillmailHashOkList) {
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
func (v GetKillmailsKillmailIdKillmailHashOkList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB5ebbfcbEncodeGithubComAntihaxGoesiEsi(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetKillmailsKillmailIdKillmailHashOkList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB5ebbfcbEncodeGithubComAntihaxGoesiEsi(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetKillmailsKillmailIdKillmailHashOkList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB5ebbfcbDecodeGithubComAntihaxGoesiEsi(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetKillmailsKillmailIdKillmailHashOkList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB5ebbfcbDecodeGithubComAntihaxGoesiEsi(l, v)
}
func easyjsonB5ebbfcbDecodeGithubComAntihaxGoesiEsi1(in *jlexer.Lexer, out *GetKillmailsKillmailIdKillmailHashOk) {
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
		case "killmail_id":
			out.KillmailId = int32(in.Int32())
		case "killmail_time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.KillmailTime).UnmarshalJSON(data))
			}
		case "victim":
			(out.Victim).UnmarshalEasyJSON(in)
		case "attackers":
			if in.IsNull() {
				in.Skip()
				out.Attackers = nil
			} else {
				in.Delim('[')
				if out.Attackers == nil {
					if !in.IsDelim(']') {
						out.Attackers = make([]GetKillmailsKillmailIdKillmailHashAttacker, 0, 1)
					} else {
						out.Attackers = []GetKillmailsKillmailIdKillmailHashAttacker{}
					}
				} else {
					out.Attackers = (out.Attackers)[:0]
				}
				for !in.IsDelim(']') {
					var v4 GetKillmailsKillmailIdKillmailHashAttacker
					easyjsonB5ebbfcbDecodeGithubComAntihaxGoesiEsi2(in, &v4)
					out.Attackers = append(out.Attackers, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "solar_system_id":
			out.SolarSystemId = int32(in.Int32())
		case "moon_id":
			out.MoonId = int32(in.Int32())
		case "war_id":
			out.WarId = int32(in.Int32())
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
func easyjsonB5ebbfcbEncodeGithubComAntihaxGoesiEsi1(out *jwriter.Writer, in GetKillmailsKillmailIdKillmailHashOk) {
	out.RawByte('{')
	first := true
	_ = first
	if in.KillmailId != 0 {
		const prefix string = ",\"killmail_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.KillmailId))
	}
	if true {
		const prefix string = ",\"killmail_time\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.KillmailTime).MarshalJSON())
	}
	if true {
		const prefix string = ",\"victim\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Victim).MarshalEasyJSON(out)
	}
	if len(in.Attackers) != 0 {
		const prefix string = ",\"attackers\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v5, v6 := range in.Attackers {
				if v5 > 0 {
					out.RawByte(',')
				}
				easyjsonB5ebbfcbEncodeGithubComAntihaxGoesiEsi2(out, v6)
			}
			out.RawByte(']')
		}
	}
	if in.SolarSystemId != 0 {
		const prefix string = ",\"solar_system_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.SolarSystemId))
	}
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
	if in.WarId != 0 {
		const prefix string = ",\"war_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.WarId))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetKillmailsKillmailIdKillmailHashOk) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB5ebbfcbEncodeGithubComAntihaxGoesiEsi1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetKillmailsKillmailIdKillmailHashOk) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB5ebbfcbEncodeGithubComAntihaxGoesiEsi1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetKillmailsKillmailIdKillmailHashOk) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB5ebbfcbDecodeGithubComAntihaxGoesiEsi1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetKillmailsKillmailIdKillmailHashOk) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB5ebbfcbDecodeGithubComAntihaxGoesiEsi1(l, v)
}
func easyjsonB5ebbfcbDecodeGithubComAntihaxGoesiEsi2(in *jlexer.Lexer, out *GetKillmailsKillmailIdKillmailHashAttacker) {
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
		case "character_id":
			out.CharacterId = int32(in.Int32())
		case "corporation_id":
			out.CorporationId = int32(in.Int32())
		case "alliance_id":
			out.AllianceId = int32(in.Int32())
		case "faction_id":
			out.FactionId = int32(in.Int32())
		case "security_status":
			out.SecurityStatus = float32(in.Float32())
		case "final_blow":
			out.FinalBlow = bool(in.Bool())
		case "damage_done":
			out.DamageDone = int32(in.Int32())
		case "ship_type_id":
			out.ShipTypeId = int32(in.Int32())
		case "weapon_type_id":
			out.WeaponTypeId = int32(in.Int32())
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
func easyjsonB5ebbfcbEncodeGithubComAntihaxGoesiEsi2(out *jwriter.Writer, in GetKillmailsKillmailIdKillmailHashAttacker) {
	out.RawByte('{')
	first := true
	_ = first
	if in.CharacterId != 0 {
		const prefix string = ",\"character_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.CharacterId))
	}
	if in.CorporationId != 0 {
		const prefix string = ",\"corporation_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.CorporationId))
	}
	if in.AllianceId != 0 {
		const prefix string = ",\"alliance_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.AllianceId))
	}
	if in.FactionId != 0 {
		const prefix string = ",\"faction_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.FactionId))
	}
	if in.SecurityStatus != 0 {
		const prefix string = ",\"security_status\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Float32(float32(in.SecurityStatus))
	}
	if in.FinalBlow {
		const prefix string = ",\"final_blow\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.FinalBlow))
	}
	if in.DamageDone != 0 {
		const prefix string = ",\"damage_done\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.DamageDone))
	}
	if in.ShipTypeId != 0 {
		const prefix string = ",\"ship_type_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.ShipTypeId))
	}
	if in.WeaponTypeId != 0 {
		const prefix string = ",\"weapon_type_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.WeaponTypeId))
	}
	out.RawByte('}')
}
