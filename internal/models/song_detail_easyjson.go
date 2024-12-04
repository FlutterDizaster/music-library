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

func easyjson4257b4a3DecodeGithubComFlutterDizasterMusicLibraryInternalModels(in *jlexer.Lexer, out *SongDetail) {
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
		case "releaseDate":
			out.ReleaseDate = string(in.String())
		case "text":
			out.Text = string(in.String())
		case "link":
			out.Link = string(in.String())
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
func easyjson4257b4a3EncodeGithubComFlutterDizasterMusicLibraryInternalModels(out *jwriter.Writer, in SongDetail) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"releaseDate\":"
		out.RawString(prefix[1:])
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
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SongDetail) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4257b4a3EncodeGithubComFlutterDizasterMusicLibraryInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SongDetail) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4257b4a3EncodeGithubComFlutterDizasterMusicLibraryInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SongDetail) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4257b4a3DecodeGithubComFlutterDizasterMusicLibraryInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SongDetail) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4257b4a3DecodeGithubComFlutterDizasterMusicLibraryInternalModels(l, v)
}