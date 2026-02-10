package utils

import (
	"strings"
)

// profanityWords contains a comprehensive list of common English profanity/obscenity words
var profanityWords = []string{
	// Strong profanity
	"fuck",
	"shit",
	"bitch",
	"asshole",
	"bastard",
	"damn",
	"cunt",
	"piss",
	"cock",
	"dick",
	"pussy",
	"whore",
	"slut",
	"fag",
	"faggot",
	"nigger",
	"nigga",
	"retard",
	"retarded",
	"motherfucker",
	"goddamn",
	"crap",
	"ass",
	"arse",
	"bollocks",
	"bugger",
	"bloody",
	"wanker",
	"twat",
	"prick",
	"douche",
	"douchebag",
	"jackass",
	"dipshit",
	"bullshit",
	"shitty",
	"fucked",
	"fucking",
	"fuckin",
	"bitching",
	"bitchy",
	"dumbass",
	"shithead",
	"dickhead",
	"asshat",
	"cumslut",
	"fuckface",
	"shitface",
	"butthole",
	"scumbag",
	"cocksucker",
	"motherfucking",
	"sonofabitch",
	"fuckboy",
	"bitchass",
	"horseshit",
	"chickenshit",
	"apeshit",
	"batshit",
	"ratass",
	"dumbfuck",
	"clusterfuck",
	"mindfuck",
	"asswipe",
	"shitbag",
	"fuckwit",
	"dickwad",
	"knobhead",
	"bellend",
	"tosser",
	"wankstain",
	"shitstain",
	"cockwomble",
	"thundercunt",
	"arsewipe",
	"shite",
	"fucknut",
	"numbnuts",
	"peckerhead",
	"dickface",
	"bumfuck",
	"assclown",
	"dickless",
	"shitless",
	"fucktard",
	"shitstorm",
	"pissed",
	"pissing",
	"pisser",
	"pissant",
	"goddammit",
	"godforsaken",
	"hellhole",
	"damnit",
	"dammit",
	"skullfuck",
	"buttfuck",
	"assfuck",
	"dickweed",
	"shitshow",
}

// ContainsProfanity checks if the given message contains any profanity words
// Returns true if profanity is found, false otherwise
// Uses case-insensitive matching
func ContainsProfanity(message string) bool {
	lowerMessage := strings.ToLower(message)

	for _, word := range profanityWords {
		if strings.Contains(lowerMessage, word) {
			return true
		}
	}

	return false
}
