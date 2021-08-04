package messages

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

func ParseMessageContent(content string) string {
	t := html.NewTokenizer(strings.NewReader(content))
	outputString := ""

	var blockQuoteCount = 0
	var emojiBlockCount = 0

	for {
		tokenType := t.Next()
		if tokenType == html.ErrorToken {
			return outputString
		}

		tagName, hasAttr := t.TagName()
		text := t.Text()

		switch tokenType {
		case html.StartTagToken:
			fmt.Printf("Start %s: hasAttr=%v\n", string(tagName), hasAttr)
			switch string(tagName){
			case "span":
				if hasAttr {
					// Check if it's an emoji
					attributes := getAttributesByToken(t)
					emoticon := getEmoticon(attributes)
					if emoticon == nil {
						continue
					} else {
						outputString += emoticon.Text
						emojiBlockCount++
					}
				}
			case "blockquote":
				fmt.Printf("blockquote=%s\n", text)
				blockQuoteCount++
			}
		case html.TextToken:
			if emojiBlockCount > 0 {
				continue
			}
			if blockQuoteCount > 0 {
				outputString += strings.Repeat("> ", blockQuoteCount)
			}
			outputString += string(text)

		case html.EndTagToken:
			switch string(tagName) {
			case "blockquote":
				blockQuoteCount--
			case "span":
				if emojiBlockCount > 0 {
					emojiBlockCount--
				}
			}
		}
	}
}

func getAttributesByToken(t *html.Tokenizer) map[string]string {
	attributes := map[string]string{}
	var hasMoreAttr = true
	for hasMoreAttr {
		var key, val []byte
		key, val, hasMoreAttr = t.TagAttr()
		attributes[string(key)] = string(val)
	}
	return attributes
}

func getEmoticon(attributes map[string]string) *Emoji {
	emojiCode := ""
	for k, v := range attributes {
		switch k {
		case "class":
			if !strings.HasPrefix(v, "animated-emoticon-") {
				return nil
			}
		case "type":
			emojiCode = v
		}
	}

	return GetEmojiByCode(emojiCode)
}
