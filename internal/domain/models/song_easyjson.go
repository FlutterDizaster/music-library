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

func easyjson366630dfDecodeGithubComFlutterDizasterMusicLibraryInternalModels(in *jlexer.Lexer, out *Song) {
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
		case "id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
			}
		case "releaseDate":
			out.ReleaseDate = string(in.String())
		case "text":
			out.Text = string(in.String())
		case "link":
			out.Link = string(in.String())
		case "group":
			out.Group = string(in.String())
		case "song":
			out.Song = string(in.String())
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
func easyjson366630dfEncodeGithubComFlutterDizasterMusicLibraryInternalModels(out *jwriter.Writer, in Song) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.RawText((in.ID).MarshalText())
	}
	{
		const prefix string = ",\"releaseDate\":"
		out.RawString(prefix)
		out.String(string(in.ReleaseDate))
	}
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix)
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"link\":"
		out.RawString(prefix)
		out.String(string(in.Link))
	}
	{
		const prefix string = ",\"group\":"
		out.RawString(prefix)
		out.String(string(in.Group))
	}
	{
		const prefix string = ",\"song\":"
		out.RawString(prefix)
		out.String(string(in.Song))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Song) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson366630dfEncodeGithubComFlutterDizasterMusicLibraryInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Song) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson366630dfEncodeGithubComFlutterDizasterMusicLibraryInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Song) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson366630dfDecodeGithubComFlutterDizasterMusicLibraryInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Song) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson366630dfDecodeGithubComFlutterDizasterMusicLibraryInternalModels(l, v)
}
