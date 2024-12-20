// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package dto

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

func easyjson56de76c1DecodeGithubComGoParkMailRu20242VKatunyMicroservicesNotificationsNotificationsDto(in *jlexer.Lexer, out *JSONResponse) {
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
		case "statusCode":
			out.HTTPStatus = int(in.Int())
		case "body":
			if m, ok := out.Body.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Body.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Body = in.Interface()
			}
		case "error":
			out.Error = string(in.String())
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
func easyjson56de76c1EncodeGithubComGoParkMailRu20242VKatunyMicroservicesNotificationsNotificationsDto(out *jwriter.Writer, in JSONResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"statusCode\":"
		out.RawString(prefix[1:])
		out.Int(int(in.HTTPStatus))
	}
	{
		const prefix string = ",\"body\":"
		out.RawString(prefix)
		if m, ok := in.Body.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Body.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Body))
		}
	}
	{
		const prefix string = ",\"error\":"
		out.RawString(prefix)
		out.String(string(in.Error))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v JSONResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComGoParkMailRu20242VKatunyMicroservicesNotificationsNotificationsDto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v JSONResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComGoParkMailRu20242VKatunyMicroservicesNotificationsNotificationsDto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *JSONResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComGoParkMailRu20242VKatunyMicroservicesNotificationsNotificationsDto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *JSONResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComGoParkMailRu20242VKatunyMicroservicesNotificationsNotificationsDto(l, v)
}
func easyjson56de76c1DecodeGithubComGoParkMailRu20242VKatunyMicroservicesNotificationsNotificationsDto1(in *jlexer.Lexer, out *EmployerNotification) {
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
			out.ID = uint64(in.Uint64())
		case "notificationText":
			out.NotificationText = string(in.String())
		case "applicantId":
			out.ApplicantID = uint64(in.Uint64())
		case "employerId":
			out.EmployerID = uint64(in.Uint64())
		case "vacancyId":
			out.VacancyID = uint64(in.Uint64())
		case "isRead":
			out.IsRead = bool(in.Bool())
		case "createdAt":
			out.CreatedAt = string(in.String())
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
func easyjson56de76c1EncodeGithubComGoParkMailRu20242VKatunyMicroservicesNotificationsNotificationsDto1(out *jwriter.Writer, in EmployerNotification) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"notificationText\":"
		out.RawString(prefix)
		out.String(string(in.NotificationText))
	}
	{
		const prefix string = ",\"applicantId\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.ApplicantID))
	}
	{
		const prefix string = ",\"employerId\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.EmployerID))
	}
	{
		const prefix string = ",\"vacancyId\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.VacancyID))
	}
	{
		const prefix string = ",\"isRead\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsRead))
	}
	{
		const prefix string = ",\"createdAt\":"
		out.RawString(prefix)
		out.String(string(in.CreatedAt))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EmployerNotification) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComGoParkMailRu20242VKatunyMicroservicesNotificationsNotificationsDto1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EmployerNotification) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComGoParkMailRu20242VKatunyMicroservicesNotificationsNotificationsDto1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EmployerNotification) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComGoParkMailRu20242VKatunyMicroservicesNotificationsNotificationsDto1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EmployerNotification) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComGoParkMailRu20242VKatunyMicroservicesNotificationsNotificationsDto1(l, v)
}
