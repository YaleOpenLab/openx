// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package database

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

func easyjson9e1087fdDecodeGithubComYaleOpenLabOpenxDatabase(in *jlexer.Lexer, out *User) {
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
		case "Index":
			out.Index = int(in.Int())
		case "EncryptedSeed":
			if in.IsNull() {
				in.Skip()
				out.EncryptedSeed = nil
			} else {
				out.EncryptedSeed = in.Bytes()
			}
		case "Name":
			out.Name = string(in.String())
		case "PublicKey":
			out.PublicKey = string(in.String())
		case "Username":
			out.Username = string(in.String())
		case "Pwhash":
			out.Pwhash = string(in.String())
		case "Address":
			out.Address = string(in.String())
		case "Description":
			out.Description = string(in.String())
		case "Image":
			out.Image = string(in.String())
		case "FirstSignedUp":
			out.FirstSignedUp = string(in.String())
		case "Kyc":
			out.Kyc = bool(in.Bool())
		case "Inspector":
			out.Inspector = bool(in.Bool())
		case "Email":
			out.Email = string(in.String())
		case "Notification":
			out.Notification = bool(in.Bool())
		case "Reputation":
			out.Reputation = float64(in.Float64())
		case "LocalAssets":
			if in.IsNull() {
				in.Skip()
				out.LocalAssets = nil
			} else {
				in.Delim('[')
				if out.LocalAssets == nil {
					if !in.IsDelim(']') {
						out.LocalAssets = make([]string, 0, 4)
					} else {
						out.LocalAssets = []string{}
					}
				} else {
					out.LocalAssets = (out.LocalAssets)[:0]
				}
				for !in.IsDelim(']') {
					var v2 string
					v2 = string(in.String())
					out.LocalAssets = append(out.LocalAssets, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "RecoveryShares":
			if in.IsNull() {
				in.Skip()
				out.RecoveryShares = nil
			} else {
				in.Delim('[')
				if out.RecoveryShares == nil {
					if !in.IsDelim(']') {
						out.RecoveryShares = make([]string, 0, 4)
					} else {
						out.RecoveryShares = []string{}
					}
				} else {
					out.RecoveryShares = (out.RecoveryShares)[:0]
				}
				for !in.IsDelim(']') {
					var v3 string
					v3 = string(in.String())
					out.RecoveryShares = append(out.RecoveryShares, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "PwdResetCode":
			out.PwdResetCode = string(in.String())
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
func easyjson9e1087fdEncodeGithubComYaleOpenLabOpenxDatabase(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Index\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Index))
	}
	{
		const prefix string = ",\"EncryptedSeed\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Base64Bytes(in.EncryptedSeed)
	}
	{
		const prefix string = ",\"Name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"PublicKey\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.PublicKey))
	}
	{
		const prefix string = ",\"Username\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"Pwhash\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Pwhash))
	}
	{
		const prefix string = ",\"Address\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Address))
	}
	{
		const prefix string = ",\"Description\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"Image\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Image))
	}
	{
		const prefix string = ",\"FirstSignedUp\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.FirstSignedUp))
	}
	{
		const prefix string = ",\"Kyc\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Kyc))
	}
	{
		const prefix string = ",\"Inspector\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Inspector))
	}
	{
		const prefix string = ",\"Email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"Notification\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Notification))
	}
	{
		const prefix string = ",\"Reputation\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Float64(float64(in.Reputation))
	}
	{
		const prefix string = ",\"LocalAssets\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.LocalAssets == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v6, v7 := range in.LocalAssets {
				if v6 > 0 {
					out.RawByte(',')
				}
				out.String(string(v7))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"RecoveryShares\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.RecoveryShares == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.RecoveryShares {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.String(string(v9))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"PwdResetCode\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.PwdResetCode))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeGithubComYaleOpenLabOpenxDatabase(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeGithubComYaleOpenLabOpenxDatabase(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeGithubComYaleOpenLabOpenxDatabase(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeGithubComYaleOpenLabOpenxDatabase(l, v)
}
