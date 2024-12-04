// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

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

func easyjsonFf3e8803DecodeGithubComFlutterDizasterMusicLibraryInternalModels(in *jlexer.Lexer, out *Library) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "songs":
			if in.IsNull() {
				in.Skip()
				out.Songs = nil
			} else {
				in.Delim('[')
				if out.Songs == nil {
					if !in.IsDelim(']') {
						out.Songs = make([]Song, 0, 0)
					} else {
						out.Songs = []Song{}
					}
				} else {
					out.Songs = (out.Songs)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Song
					(v1).UnmarshalEasyJSON(in)
					out.Songs = append(out.Songs, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "pagination":
			easyjsonFf3e8803DecodeGithubComFlutterDizasterMusicLibraryInternalModels1(in, &out.Pagination)
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
func easyjsonFf3e8803EncodeGithubComFlutterDizasterMusicLibraryInternalModels(out *jwriter.Writer, in Library) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"songs\":"
		out.RawString(prefix[1:])
		if in.Songs == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Songs {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"pagination\":"
		out.RawString(prefix)
		easyjsonFf3e8803EncodeGithubComFlutterDizasterMusicLibraryInternalModels1(out, in.Pagination)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Library) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonFf3e8803EncodeGithubComFlutterDizasterMusicLibraryInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Library) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonFf3e8803EncodeGithubComFlutterDizasterMusicLibraryInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Library) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonFf3e8803DecodeGithubComFlutterDizasterMusicLibraryInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Library) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonFf3e8803DecodeGithubComFlutterDizasterMusicLibraryInternalModels(l, v)
}
func easyjsonFf3e8803DecodeGithubComFlutterDizasterMusicLibraryInternalModels1(in *jlexer.Lexer, out *Pagination) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "page":
			out.Page = int(in.Int())
		case "pageSize":
			out.PageSize = int(in.Int())
		case "total":
			out.Total = int(in.Int())
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
func easyjsonFf3e8803EncodeGithubComFlutterDizasterMusicLibraryInternalModels1(out *jwriter.Writer, in Pagination) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"page\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Page))
	}
	{
		const prefix string = ",\"pageSize\":"
		out.RawString(prefix)
		out.Int(int(in.PageSize))
	}
	{
		const prefix string = ",\"total\":"
		out.RawString(prefix)
		out.Int(int(in.Total))
	}
	out.RawByte('}')
}
