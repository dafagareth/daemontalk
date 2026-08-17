package templates

func reactionName(emoji string) string {
	switch emoji {
	case "like":
		return "Like"
	case "heart":
		return "Love"
	case "fire":
		return "Fire"
	case "mindblown":
		return "Mind Blown"
	case "insightful":
		return "Insightful"
	default:
		return "Reaction"
	}
}
