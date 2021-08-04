package messages

type Emoji struct {
	Code string
	Text string
}

var emojiMap = map[string]string{
	"smile": "🙂",
	"laugh": "😀",
}

func GetEmojiByCode(code string) *Emoji {
	val, ok := emojiMap[code]
	if !ok {
		return nil
	}
	return &Emoji{
		Code: code,
		Text: val,
	}
}
