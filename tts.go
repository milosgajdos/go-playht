package playht

type VoiceEngine string

const (
	PlayHTv1      VoiceEngine = "PlayHT1.0"
	PlayHTv2      VoiceEngine = "PlayHT2.0"
	PlayHTv2Turbo VoiceEngine = "PlayHT2.0-turbo"
)

func (v VoiceEngine) String() string {
	return string(v)
}

type OutputFormat string

const (
	Mp3   OutputFormat = "mp3"
	Wav   OutputFormat = "wav"
	Ogg   OutputFormat = "ogg"
	Flac  OutputFormat = "flac"
	Mulaw OutputFormat = "mulaw"
)

func (o OutputFormat) String() string {
	return string(o)
}

type Quality string

const (
	Draft   Quality = "draft"
	Low     Quality = "low"
	Medium  Quality = "medium"
	High    Quality = "high"
	Premium Quality = "premium"
)

func (q Quality) String() string {
	return string(q)
}

type Emotion string

const (
	FemaleHappy     Emotion = "female_happy"
	FemaleSad       Emotion = "female_sad"
	FemaleAngry     Emotion = "female_angry"
	FemaleFearful   Emotion = "female_fearful"
	FemaleDisgust   Emotion = "female_disgust"
	FemaleSurprised Emotion = "female_surprised"
	MaleHappy       Emotion = "male_happy"
	MaleSad         Emotion = "male_sad"
	MaleAngry       Emotion = "male_angry"
	MaleFearful     Emotion = "male_fearful"
	MaleDisgust     Emotion = "male_disgust"
	MaleSurprised   Emotion = "male_surprised"
)

func (e Emotion) String() string {
	return string(e)
}
