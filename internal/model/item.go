package model

type Item struct {
	Name         string `xml:"name,attr"`
	Lemma        string `xml:"lemma,attr"`
	MainWord     string `xml:"main_word,attr"`
	SyntType     string `xml:"synt_type,attr"`
	Poses        string `xml:"poses,attr"`
	Meaning      string `xml:"meaning,attr"`
	ID           string `xml:"id,attr"`
	SynsetID     string `xml:"synset_id,attr"`
	PartOfSpeech string `xml:"part_of_speech,attr"`
	ConceptID    string `xml:"concept_id,attr"`
	EntryID      string `xml:"entry_id,attr"`
	Text         string `xml:",chardata"`
}

type ItemSet struct {
	Items []Item `xml:"Item"`
}
