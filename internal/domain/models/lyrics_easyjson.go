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

func easyjson6549bd74DecodeGithubComFlutterDizasterMusicLibraryInternalModels(in *jlexer.Lexer, out *Lyrics) {
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
		case "lyrics":
			out.Lyrics = string(in.String())
		case "pagination":
			easyjson6549bd74DecodeGithubComFlutterDizasterMusicLibraryInternalModels1(in, &out.Pagination)
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
func easyjson6549bd74EncodeGithubComFlutterDizasterMusicLibraryInternalModels(out *jwriter.Writer, in Lyrics) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"lyrics\":"
		out.RawString(prefix[1:])
		out.String(string(in.Lyrics))
	}
	{
		const prefix string = ",\"pagination\":"
		out.RawString(prefix)
		easyjson6549bd74EncodeGithubComFlutterDizasterMusicLibraryInternalModels1(out, in.Pagination)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Lyrics) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6549bd74EncodeGithubComFlutterDizasterMusicLibraryInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Lyrics) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6549bd74EncodeGithubComFlutterDizasterMusicLibraryInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Lyrics) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6549bd74DecodeGithubComFlutterDizasterMusicLibraryInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Lyrics) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6549bd74DecodeGithubComFlutterDizasterMusicLibraryInternalModels(l, v)
}
func easyjson6549bd74DecodeGithubComFlutterDizasterMusicLibraryInternalModels1(in *jlexer.Lexer, out *Pagination) {
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
		case "limit":
			out.Limit = int(in.Int())
		case "offset":
			out.Offset = int(in.Int())
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
func easyjson6549bd74EncodeGithubComFlutterDizasterMusicLibraryInternalModels1(out *jwriter.Writer, in Pagination) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"limit\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Limit))
	}
	{
		const prefix string = ",\"offset\":"
		out.RawString(prefix)
		out.Int(int(in.Offset))
	}
	{
		const prefix string = ",\"total\":"
		out.RawString(prefix)
		out.Int(int(in.Total))
	}
	out.RawByte('}')
}