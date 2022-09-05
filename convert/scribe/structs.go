package scribe

type EntriesHeader struct {
	Unk         []byte `json:",omitempty"`
	UnkB64      string `json:"unk"`
	EntriesSize int32  `json:",omitempty"`
	UnkTwo      int32  `json:"unk_two"`
}

type Entry struct {
	Size      int32  `json:",omitempty"`
	Unk       int32  `json:"unk"`
	UnkTwo    []byte `json:",omitempty"`
	UnkTwoB64 string `json:"unk_two"`
	// 12 null
	TypeStringLen int32  `json:",omitempty"`
	TypeString    string `json:"type_string"`
	// 2-8 null, seems arbitrary
	TextLen int32 `json:",omitempty"`
	// int32, 8192
	Text string `json:"text"`
	// 2-8 null, seems arbitrary
}

type Scribe struct {
	EntriesHeader *EntriesHeader `json:"entries_header"`
	Entries       []*Entry       `json:"entries"`
	EndData       []byte         `json:",omitempty"`
	EndDataB64    string         `json:"end_data"`
}
