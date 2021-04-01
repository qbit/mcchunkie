package plugins

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/caneroj1/stemmer"
	"github.com/matrix-org/gomatrix"
)

// Toki responds to toki pona word queries
type Toki struct {
	POS          string
	Meanings     []string
	MeaningStems []string
	Alt          string
	Principle    string
}

// Print prints the definition
func (t *Toki) Print(w string) string {
	return fmt.Sprintf("**%s**: (_%s_) %s", w, t.POS, strings.Join(t.Meanings, ", "))
}

// Words prints the definition
func (t *Toki) Words() []string {
	s := strings.Join(t.Meanings, " ")
	w := strings.Split(s, " ")
	contains := make(map[string]bool)
	var result []string

	for _, x := range w {
		if ok := contains[x]; !ok {
			contains[x] = true
			result = append(result, x)
		}
	}

	return result
}

// TokiLang is our full representation of toki pona
var TokiLang = map[string][]Toki{
	"telo (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to wash with water",
				"to put water to",
				"to melt",
				"to liquify",
				"to water",
			},
		},
	},
	"nanpa (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to reckon",
				"to number",
				"to count",
			},
		},
	},
	"kasi (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to grow",
				"to plant",
			},
		},
	},
	"ken (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to enable",
				"to allow",
				"to permit",
				"to make possible",
			},
		},
	},
	"kiwen": {
		{
			POS: "adjective",
			Meanings: []string{
				"solid",
				"stone-like",
				"made of stone or metal",
				"hard",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"solid",
				"stone-like",
				"made of stone or metal",
				"hard",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"rock",
				"stone",
				"metal",
				"mineral",
				"clay",
				"hard thing",
			},
		},
	},
	"weka": {
		{
			POS: "adjective",
			Meanings: []string{
				"away",
				"ignored",
				"absent",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"absence",
			},
		},
	},
	"open (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to start",
				"to begin",
				"to turn on",
				"to open",
			},
		},
	},
	"kama (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to summon",
				"to bring about",
			},
		},
	},
	"walo": {
		{
			POS: "adjective",
			Meanings: []string{
				"whitish",
				"light-coloured",
				"pale",
				"white",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"whiteness",
				"lightness",
				"white thing or part",
			},
		},
	},
	"mi ' pona, tan ni": {},
	"anpa": {
		{
			POS: "adjective",
			Meanings: []string{
				"lower",
				"bottom",
				"down",
				"low",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"below",
				"deep",
				"low",
				"deeply",
				"downstairs",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"lower part",
				"under",
				"below",
				"floor",
				"beneath",
				"bottom",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to prostrate oneself",
			},
		},
	},
	"lukin": {
		{
			POS: "adjective",
			Meanings: []string{
				"visual(ly)",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"visual(ly)",
			},
		},
		{
			POS: "auxiliary verb",
			Meanings: []string{
				"try to",
				"look for",
				"to seek to",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"look",
				"glance",
				"sight",
				"gaze",
				"glimpse",
				"seeing",
				"vision",
				"view",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to watch out",
				"to pay attention",
				"to look",
			},
		},
	},
	"pali (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to make",
				"to build",
				"to create",
				"to do",
			},
		},
	},
	"musi (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to entertain",
				"to amuse",
			},
		},
	},
	"mu!": {
		{
			POS: "interjection",
			Meanings: []string{
				"woof! meow! moo! etc. (cute animal noise)",
			},
		},
	},
	"weka (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to eliminate",
				"to throw away",
				"to get rid of",
				"to remove",
			},
		},
	},
	"namako": {
		{
			POS: "adjective",
			Meanings: []string{
				"piquant",
				"spicy",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"something extra",
				"food additive",
				"accessory",
				"spice",
			},
		},
	},
	"pini": {
		{
			POS: "adjective",
			Meanings: []string{
				"finished",
				"past",
				"done",
				"completed",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"past",
				"perfectly",
				"ago",
			},
		},
		{
			POS: "auxiliary verb",
			Meanings: []string{
				"to finish",
				"to end",
				"to interrupt",
				"to stop",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"tip",
				"end",
			},
		},
	},
	"pakala!": {
		{
			POS: "interjection",
			Meanings: []string{
				"damn! fuck!",
			},
		},
	},
	"kama moli": {
		{
			POS: "intransitives verb",
			Meanings: []string{
				"dieing",
			},
		},
	},
	"mi moku, tan ni": {},
	"moli (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to kill",
			},
		},
	},
	"mi wile e ni": {},
	"kalama": {
		{
			POS: "adjective",
			Meanings: []string{
				"loud",
				"rowdy",
				"noisy",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"noise",
				"voice",
				"sound",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to make noise",
			},
		},
	},
	"linja": {
		{
			POS: "adjective",
			Meanings: []string{
				"oblong",
				"long",
				"elongated",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"rope",
				"hair",
				"thread",
				"cord",
				"chain",
				"line",
				"yarn",
				"long and flexible thing; string",
			},
		},
	},
	"lape": {
		{
			POS: "adjective",
			Meanings: []string{
				"of sleep",
				"dormant",
				"sleeping",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"asleep",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"rest",
				"sleep",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to rest",
				"to sleep",
			},
		},
	},
	"tenpo": {
		{
			POS: "adjective",
			Meanings: []string{
				"chronological",
				"chronologic",
				"temporal",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"chronologically",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"period of time",
				"moment",
				"duration",
				"situation",
				"occasion",
				"time",
			},
		},
	},
	"sewi": {
		{
			POS: "adjective",
			Meanings: []string{
				"elevated",
				"religious",
				"formal",
				"superior",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"elevated",
				"religious",
				"formal",
				"superior",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"up",
				"above",
				"top",
				"over",
				"on",
				"high",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to get up",
			},
		},
	},
	"kon (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to puff away something",
				"to blow away something",
			},
		},
	},
	"waso": {
		{
			POS: "adjective",
			Meanings: []string{
				"bird-",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"bat; flying creature",
				"winged animal",
				"bird",
			},
		},
	},
	"sitelen": {
		{
			POS: "adjective",
			Meanings: []string{
				"pictorial",
				"metaphorical",
				"metaphorisch",
				"figurative",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"pictorially",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"image",
				"representation",
				"symbol",
				"mark",
				"writing",
				"picture",
			},
		},
	},
	"sin (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to renovate",
				"to freshen",
				"to renew",
			},
		},
	},
	"sike (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to circle",
				"to revolve",
				"to circle around",
				"to rotate",
				"to orbit",
			},
		},
	},
	"unpa": {
		{
			POS: "adjective",
			Meanings: []string{
				"sexual",
				"erotic",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"sexual",
				"erotic",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"sexuality",
				"sex",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to have sex",
			},
		},
	},
	"sijelo (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to heal up",
				"to cure",
				"to heal",
			},
		},
	},
	"palisa (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to beat",
				"to poke",
				"to stab",
				"to sexually arouse",
				"to stretch",
			},
		},
	},
	"pakala (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to ruin",
				"to break",
				"to hurt",
				"to injure",
				"to damage",
				"to screw up",
			},
		},
	},
	"alasa (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to forage",
				"to hunt",
			},
		},
	},
	"insa": {
		{
			POS: "adjective",
			Meanings: []string{
				"internal",
				"inner",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"inner world",
				"centre",
				"stomach",
				"inside",
			},
		},
	},
	"ko": {
		{
			POS: "noun",
			Meanings: []string{
				"dough",
				"glue",
				"paste",
				"powder",
				"gum",
				"semi-solid or squishy substance; clay",
			},
		},
	},
	"len": {
		{
			POS: "adjective",
			Meanings: []string{
				"clothed",
				"costumed",
				"dressed up",
				"dressed",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"cloth",
				"fabric",
				"network",
				"internet",
				"clothing",
			},
		},
	},
	"lawa": {
		{
			POS: "adjective",
			Meanings: []string{
				"leading",
				"in charge",
				"main",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"leading",
				"in charge",
				"main",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"mind",
				"head",
			},
		},
	},
	"sitelen (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to write",
				"to draw",
			},
		},
	},
	"!": {
		{
			POS: "separator",
			Meanings: []string{
				"'.",
			},
		},
	},
	"\"": {
		{
			POS: "separator",
			Meanings: []string{
				"Quotation marks are used for words with original spelling or for quotes.",
			},
		},
	},
	"lupa": {
		{
			POS: "adjective",
			Meanings: []string{
				"holey",
				"full of holes",
				"hole-",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"orifice",
				"door",
				"window",
				"hole",
			},
		},
	},
	"#": {
		{
			POS: "unofficial",
			Meanings: []string{
				"Number sign",
			},
		},
	},
	"kama jo (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to get",
			},
		},
	},
	"sijelo": {
		{
			POS: "adjective",
			Meanings: []string{
				"bodily",
				"corporal",
				"corporeal",
				"material",
				"carnal",
				"physical",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"bodily",
				"physically",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"physical state",
				"torso",
				"body (of person or animal)",
			},
		},
	},
	"pimeja (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to darken",
			},
		},
	},
	"ona li wile e ni": {},
	"a a a!": {
		{
			POS: "interjection",
			Meanings: []string{
				"laugh",
			},
		},
	},
	"kulupu (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to call together",
				"to convene",
				"to assemble",
			},
		},
	},
	"'": {
		{
			POS: "unofficial",
			Meanings: []string{
				"An apostrophe can identify a predicate that does not contain a verb.",
			},
		},
	},
	"tawa (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to displace",
				"to move",
			},
		},
	},
	"soweli": {
		{
			POS: "adjective",
			Meanings: []string{
				"animal",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"especially land mammal",
				"lovable animal",
				"beast",
				"animal",
			},
		},
	},
	"en": {
		{
			POS: "conjunction",
			Meanings: []string{
				"and (used to coordinate head nouns)",
			},
		},
	},
	"jo (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to contain",
				"to have",
			},
		},
	},
	"wile (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"need",
				"wish",
				"have to",
				"must",
				"will",
				"should",
				"to want",
			},
		},
	},
	",": {
		{
			POS: "separator",
			Meanings: []string{
				"A comma is used after an 'o' to addressing people. Optional you can put a comma before a preposition. Don't use a comma before or after",
			},
		},
	},
	"palisa": {
		{
			POS: "adjective",
			Meanings: []string{
				"long",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"rod",
				"stick",
				"pointy thing",
				"long hard thing; branch",
			},
		},
	},
	"alasa": {
		{
			POS: "adjective",
			Meanings: []string{
				"-hunting",
				"hunting",
				"hunting-",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"hunting",
			},
		},
	},
	"la": {
		{
			POS: "separator",
			Meanings: []string{
				"half sentence or noun. Don't use 'la' before or after",
				"A 'la' is between a conditional phrases and the main sentence. A context phrase can be sentence",
			},
		},
	},
	".": {
		{
			POS: "separator",
			Meanings: []string{
				"'.",
			},
		},
	},
	"ike!": {
		{
			POS: "interjection",
			Meanings: []string{
				"oh dear! woe! alas!",
			},
		},
	},
	"suli (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to lengthen",
				"to enlarge",
			},
		},
	},
	"tomo (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to construct",
				"to engineer",
				"to build",
			},
		},
	},
	"toki": {
		{
			POS: "adjective",
			Meanings: []string{
				"eloquent",
				"linguistic",
				"verbal",
				"grammatical",
				"speaking",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"eloquent",
				"linguistic",
				"verbal",
				"grammatical",
				"speaking",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"speech",
				"tongue",
				"lingo",
				"jargon",
				"",
				"language",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to chat",
				"to communicate",
				"to talk",
			},
		},
	},
	"taso": {
		{
			POS: "adjective",
			Meanings: []string{
				"sole",
				"only",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"just",
				"merely",
				"simply",
				"solely",
				"singly",
				"only",
			},
		},
		{
			POS: "conjunction",
			Meanings: []string{
				"however",
				"but",
			},
		},
	},
	"li": {
		{
			POS: "separator",
			Meanings: []string{
				"'",
				"'.",
				"'",
			},
		},
	},
	"suli": {
		{
			POS: "adjective",
			Meanings: []string{
				"tall",
				"long",
				"adult",
				"important",
				"big",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"tall",
				"long",
				"adult",
				"important",
				"big",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"size",
			},
		},
	},
	"selo mi li wile e ni": {},
	"pan (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to sow",
			},
		},
	},
	"sewi (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to lift",
			},
		},
	},
	"sama": {
		{
			POS: "adjective",
			Meanings: []string{
				"similar",
				"equal",
				"of equal status or position",
				"same",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"equally",
				"exactly the same",
				"just the same",
				"similarly",
				"just as",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"parity",
				"equity",
				"identity",
				"par",
				"sameness",
				"equality",
			},
		},
		{
			POS: "preposition",
			Meanings: []string{
				"as",
				"seem",
				"like",
			},
		},
	},
	"pona la": {
		{
			POS: "noun",
			Meanings: []string{
				"if simplicity",
				"if positivity",
				"if good",
			},
		},
	},
	"ike (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to worsen",
				"to make bad",
			},
		},
	},
	"kule (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to color",
				"to paint",
			},
		},
	},
	"lili (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to shorten",
				"to shrink",
				"to lessen",
				"to reduce",
			},
		},
	},
	"pali": {
		{
			POS: "adjective",
			Meanings: []string{
				"work-related",
				"operating",
				"working",
				"active",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"briskly",
				"actively",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"work",
				"deed",
				"project",
				"activity",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to work",
				"to function",
				"to act",
			},
		},
	},
	"ala": {
		{
			POS: "adjective",
			Meanings: []string{
				"not",
				"none",
				"un-",
				"no",
			},
		},
		{
			POS: "adjective numeral",
			Meanings: []string{
				"0",
				"null",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"don't",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"negation",
				"zero",
				"nothing",
			},
		},
	},
	"?": {
		{
			POS: "separator",
			Meanings: []string{
				"'.",
			},
		},
	},
	"selo (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to protect",
				"to guard",
				"to shelter",
			},
		},
	},
	"ale": {
		{
			POS: "adjective",
			Meanings: []string{
				"every",
				"complete",
				"whole (ale = ali)",
				"(depreciated)",
				"all",
			},
		},
		{
			POS: "adjective numeral",
			Meanings: []string{
				"100 (official Toki Pona book)",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"forever",
				"evermore",
				"eternally (ale = ali)",
				"(depreciated)",
				"always",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"anything",
				"life",
				"the universe",
				"(depreciated)",
				"everything",
			},
		},
	},
	"jaki!": {
		{
			POS: "interjection",
			Meanings: []string{
				"ew! yuck!",
			},
		},
	},
	"ken": {
		{
			POS: "auxiliary verb",
			Meanings: []string{
				"may",
				"to can",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"ability",
				"power to do things",
				"permission",
				"possibility",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"is able to",
				"is allowed to",
				"may",
				"is possible",
				"can",
			},
		},
	},
	"kin la": {
		{
			POS: "noun",
			Meanings: []string{
				"if fact",
				"if reality",
			},
		},
	},
	"ali": {
		{
			POS: "adjective",
			Meanings: []string{
				"every",
				"complete",
				"whole (ale = ali)",
				"all",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"forever",
				"evermore",
				"eternally (ale = ali)",
				"always",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"anything",
				"life",
				"the universe",
				"everything",
			},
		},
	},
	"ante la": {
		{
			POS: "noun",
			Meanings: []string{
				"if variance",
				"if disagreement",
				"if difference",
			},
		},
	},
	"esun (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to sell",
				"to barter",
				"to swap",
				"to buy",
			},
		},
	},
	"anpa (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to beat",
				"to vanquish",
				"to conquer",
				"to enslave",
				"to defeat",
			},
		},
	},
	"pipi": {
		{
			POS: "noun",
			Meanings: []string{
				"insect",
				"spider",
				"bug",
			},
		},
	},
	"open": {
		{
			POS: "adjective",
			Meanings: []string{
				"starting",
				"opening",
				"initial",
			},
		},
		{
			POS: "auxiliary verb",
			Meanings: []string{
				"to start",
				"to begin",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"beginning",
				"opening",
				"start",
			},
		},
	},
	"o!": {
		{
			POS: "interjection",
			Meanings: []string{
				"hey! (calling somebody's attention)",
			},
		},
	},
	"wan": {
		{
			POS: "adjective numeral",
			Meanings: []string{
				"1",
				"one",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"element",
				"particle",
				"part",
				"piece",
				"unit",
			},
		},
	},
	"telo": {
		{
			POS: "adjective",
			Meanings: []string{
				"slobbery",
				"moist",
				"damp",
				"humid",
				"sticky",
				"sweaty",
				"dewy",
				"drizzly",
				"wett",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"slobbery",
				"moist",
				"damp",
				"humid",
				"sticky",
				"sweaty",
				"dewy",
				"drizzly",
				"wett",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"liquid",
				"juice",
				"sauce",
				"water",
			},
		},
	},
	"pona (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to fix",
				"to repair",
				"to make good",
				"to improve",
			},
		},
	},
	"ma": {
		{
			POS: "adjective",
			Meanings: []string{
				"outdoor",
				"alfresco",
				"open-air",
				"countrified",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"earth",
				"country",
				"(outdoor) area",
				"land",
			},
		},
	},
	"sinpin": {
		{
			POS: "adjective",
			Meanings: []string{
				"frontal",
				"anterior",
				"vertical",
				"facial",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"foremost",
				"front",
				"wall",
				"chest",
				"torso",
				"face",
			},
		},
	},
	"poka": {
		{
			POS: "adjective",
			Meanings: []string{
				"neighbouring",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"hip",
				"next to",
				"side",
			},
		},
	},
	"seli": {
		{
			POS: "adjective",
			Meanings: []string{
				"warm",
				"cooked",
				"hot",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"warm",
				"cooked",
				"hot",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"warmth",
				"heat",
				"fire",
			},
		},
	},
	"luka": {
		{
			POS: "adjective numeral",
			Meanings: []string{
				"5",
				"five",
			},
		},
		{
			POS: "adjective",
			Meanings: []string{
				"palpable",
				"tangible",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"hand",
				"tacticle organ",
				"arm",
			},
		},
	},
	"sin": {
		{
			POS: "adjective",
			Meanings: []string{
				"fresh",
				"another",
				"more",
				"new",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"regenerative",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"novelty",
				"innovation",
				"newness",
				"new release",
				"news",
			},
		},
	},
	"pimeja": {
		{
			POS: "adjective",
			Meanings: []string{
				"dark",
				"black",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"shadows",
				"darkness",
			},
		},
	},
	"wile": {
		{
			POS: "auxiliary verb",
			Meanings: []string{
				"need",
				"wish",
				"have to",
				"must",
				"will",
				"should",
				"to want",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"need",
				"will",
				"desire",
			},
		},
	},
	"olin (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to love (a person)",
			},
		},
	},
	"mi": {
		{
			POS: "personal pronoun",
			Meanings: []string{
				"we",
				"I",
			},
		},
		{
			POS: "possessive pronoun",
			Meanings: []string{
				"our",
				"my",
			},
		},
	},
	"selo": {
		{
			POS: "noun",
			Meanings: []string{
				"outer form",
				"bark",
				"peel",
				"shell",
				"skin",
				"boundary",
				"shape",
				"skin",
			},
		},
	},
	"poki": {
		{
			POS: "noun",
			Meanings: []string{
				"box",
				"bowl",
				"cup",
				"glass",
				"container",
			},
		},
	},
	"o,": {
		{
			POS: "interjection",
			Meanings: []string{
				"adressing people",
			},
		},
	},
	"mute (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to make many or much",
			},
		},
	},
	"jaki": {
		{
			POS: "adjective",
			Meanings: []string{
				"gross",
				"filthy",
				"obscene",
				"dirty",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"gross",
				"filthy",
				"dirty",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"pollution",
				"garbage",
				"filth",
				"feces",
				"dirt",
			},
		},
	},
	"mun": {
		{
			POS: "adjective",
			Meanings: []string{
				"lunar",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"lunar",
				"night sky object",
				"star",
				"moon",
			},
		},
	},
	"loje": {
		{
			POS: "adjective",
			Meanings: []string{
				"ruddy",
				"pink",
				"pinkish",
				"gingery",
				"reddish",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"red",
			},
		},
	},
	"sike": {
		{
			POS: "adjective",
			Meanings: []string{
				"cyclical",
				"of one year",
				"round",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"rotated",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"ball",
				"cycle",
				"sphere",
				"wheel; round or circular thing",
				"circle",
			},
		},
	},
	"ijo (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to objectify",
			},
		},
	},
	"nasa": {
		{
			POS: "adjective",
			Meanings: []string{
				"crazy",
				"foolish",
				"drunk",
				"strange",
				"stupid",
				"weird",
				"silly",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"crazy",
				"foolish",
				"drunk",
				"strange",
				"stupid",
				"weird",
				"silly",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"foolishness",
				"silliness",
				"nonsense",
				"idiocy",
				"obtuseness",
				"muddler",
				"stupidity",
			},
		},
	},
	"mi pilin e ni": {},
	"ike la": {
		{
			POS: "noun",
			Meanings: []string{
				"if badness",
				"if evil",
				"if negativity",
			},
		},
	},
	"kiwen (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to harden",
				"to petrify",
				"to fossilize",
				"to solidify",
			},
		},
	},
	"mu (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to make animal noise",
			},
		},
	},
	"noka": {
		{
			POS: "adjective",
			Meanings: []string{
				"lower",
				"bottom",
				"foot-",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"on foot",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"foot; organ of locomotion; bottom",
				"lower part",
				"leg",
			},
		},
	},
	"o !": {
		{
			POS: "separator",
			Meanings: []string{
				"'o' replace 'li'.",
			},
		},
		{
			POS: "subject",
			Meanings: []string{
				"An 'o' is used for imperative (commands). 'o' replace the subject.",
			},
		},
	},
	"mu": {
		{
			POS: "adjective",
			Meanings: []string{
				"animal nois-",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"animal nois-",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"animal noise",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to communicate animally",
			},
		},
	},
	"a": {
		{
			POS: "interjection",
			Meanings: []string{
				"ha",
				"uh",
				"oh",
				"ooh",
				"aw",
				"well (emotion word)",
				"ah",
			},
		},
	},
	"oko": {
		{
			POS: "adjective",
			Meanings: []string{
				"eye-",
				"optical",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"eye",
			},
		},
	},
	"kala": {
		{
			POS: "adjective",
			Meanings: []string{
				"fish-",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"marine animal",
				"sea creature",
				"fish",
			},
		},
	},
	"e sina": {
		{
			POS: "reflexive pronoun",
			Meanings: []string{
				"yourselves",
				"yourself",
			},
		},
	},
	"nasa (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to make weird",
				"to drive crazy",
			},
		},
	},
	"e": {
		{
			POS: "separator",
			Meanings: []string{
				"'",
				"'.",
				"'",
			},
		},
	},
	"ijo": {
		{
			POS: "adjective",
			Meanings: []string{
				"of something",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"of something",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"something",
				"stuff",
				"anything",
				"object",
				"thing",
			},
		},
	},
	"pona!": {
		{
			POS: "interjection",
			Meanings: []string{
				"great! good! thanks! OK! cool! yay!",
			},
		},
	},
	"ante (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to alter",
				"to modify",
				"to change",
			},
		},
	},
	"akesi": {
		{
			POS: "adjective",
			Meanings: []string{
				"reptilian-",
				"slimy",
				"amphibian-",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"amphibian; non-cute animal",
				"reptile",
			},
		},
	},
	"seme": {
		{
			POS: "question pronoun",
			Meanings: []string{
				"which",
				"wh- (question word)",
				"what",
			},
		},
	},
	"nimi (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to name",
			},
		},
	},
	"e ona": {
		{
			POS: "reflexive pronoun",
			Meanings: []string{
				"herself",
				"itself",
				"themselves",
				"himself",
			},
		},
	},
	"mije": {
		{
			POS: "adjective",
			Meanings: []string{
				"masculine",
				"manly",
				"male",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"male",
				"husband",
				"boyfriend",
				"man",
			},
		},
	},
	"mama": {
		{
			POS: "adjective",
			Meanings: []string{
				"parental",
				"maternal",
				"fatherly",
				"motherly",
				"mumsy",
				"of the parent",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"mother",
				"father",
				"parent",
			},
		},
	},
	"tu": {
		{
			POS: "adjective numeral",
			Meanings: []string{
				"2",
				"two",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"pair",
				"duo",
			},
		},
	},
	"jaki (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to dirty",
				"to pollute",
			},
		},
	},
	"wan (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to make one",
				"to unite",
			},
		},
	},
	"suno (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to illumine",
				"to light",
			},
		},
	},
	"ni": {
		{
			POS: "adjective demonstrative pronoun",
			Meanings: []string{
				"that",
				"this",
			},
		},
		{
			POS: "noun demonstrative pronoun",
			Meanings: []string{
				"that",
				"this",
			},
		},
	},
	"kute (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to listen",
				"",
				"to hear",
			},
		},
	},
	"pana": {
		{
			POS: "adjective",
			Meanings: []string{
				"generous",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"transfer",
				"exchange",
				"giving",
			},
		},
	},
	"nanpa": {
		{
			POS: "adjective numeral",
			Meanings: []string{
				"To build ordinal numbers.",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"numeral",
				"number",
			},
		},
	},
	"lupa (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to stab",
				"to perforate",
				"to pierce",
			},
		},
	},
	"tomo": {
		{
			POS: "adjective",
			Meanings: []string{
				"domestic",
				"household",
				"urban",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"domestic",
				"household",
				"urban",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"e.g. house",
				"home",
				"room",
				"building",
				"indoor constructed space",
			},
		},
	},
	"nasin": {
		{
			POS: "adjective",
			Meanings: []string{
				"habitual",
				"customary",
				"doctrinal",
				"systematic",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"manner",
				"custom",
				"road",
				"path",
				"doctrine",
				"system",
				"method",
				"way",
			},
		},
	},
	"kepeken": {
		{
			POS: "noun",
			Meanings: []string{
				"usage",
				"tool",
				"use",
			},
		},
		{
			POS: "preposition",
			Meanings: []string{
				"using",
				"with",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to use",
			},
		},
	},
	"laso": {
		{
			POS: "adjective",
			Meanings: []string{
				"bluey",
				"bluish",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"blue-green",
				"blue",
			},
		},
	},
	"ko (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to pulverize",
				"to squash",
			},
		},
	},
	"lipu": {
		{
			POS: "adjective",
			Meanings: []string{
				"paper-",
				"card-",
				"ticket-",
				"sheet-",
				"page",
				"-",
				"book-",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"book",
				"card",
				"ticket",
				"sheet",
				"(web-)page",
				"list ; flat and bendable thing",
				"paper",
			},
		},
	},
	"suno": {
		{
			POS: "adjective",
			Meanings: []string{
				"sunnily",
				"sunny",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"sunnily",
				"sunny",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"light",
				"sun",
			},
		},
	},
	"sina toki e ni, tawa mi": {},
	"open la": {
		{
			POS: "noun",
			Meanings: []string{
				"in the beginning",
				"at the opening",
			},
		},
	},
	"ike": {
		{
			POS: "adjective",
			Meanings: []string{
				"negative",
				"wrong",
				"evil",
				"overly complex",
				"bad",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"negative",
				"wrong",
				"evil",
				"overly complex",
				"bad",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"badness",
				"evil",
				"negativity",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to suck",
				"to be bad",
			},
		},
	},
	"kule": {
		{
			POS: "adjective",
			Meanings: []string{
				"pigmented",
				"painted",
				"colourful",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"colour",
				"paint",
				"ink",
				"dye",
				"hue",
				"color",
			},
		},
	},
	"moku": {
		{
			POS: "adjective",
			Meanings: []string{
				"eating",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"eating",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"meal",
				"food",
			},
		},
	},
	"kasi": {
		{
			POS: "adjective",
			Meanings: []string{
				"vegetal",
				"biological",
				"biologic",
				"leafy",
				"vegetable",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"vegetation",
				"herb",
				"leaf",
				"plant",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to grow",
			},
		},
	},
	"musi": {
		{
			POS: "adjective",
			Meanings: []string{
				"fun",
				"recreational",
				"artful",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"cheerfully",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"playing",
				"game",
				"recreation",
				"art",
				"entertainment",
				"fun",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to have fun",
				"to play",
			},
		},
	},
	"awen": {
		{
			POS: "adjective",
			Meanings: []string{
				"stationary",
				"permanent",
				"sedentary",
				"remaining",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"yet",
				"still",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"continuity",
				"continuum",
				"stay",
				"inertia",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to wait",
				"to remain",
				"to stay",
			},
		},
	},
	"e mi": {
		{
			POS: "reflexive pronoun",
			Meanings: []string{
				"ourselves",
				"myself",
			},
		},
	},
	"uta": {
		{
			POS: "adjective",
			Meanings: []string{
				"oral",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"orally",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"lips",
				"oral cavity",
				"jaw",
				"beak",
				"mouth",
			},
		},
	},
	"kama": {
		{
			POS: "adjective",
			Meanings: []string{
				"future",
				"coming",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"future",
				"coming",
			},
		},
		{
			POS: "auxiliary verb",
			Meanings: []string{
				"to mange to",
				"to become",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"happening",
				"chance",
				"arrival",
				"beginning",
				"event",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to become",
				"to arrive",
				"to happen",
				"to come",
			},
		},
	},
	"wawa": {
		{
			POS: "adjective",
			Meanings: []string{
				"strong",
				"fierce",
				"intense",
				"sure",
				"confident",
				"energetic",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"powerfully",
				"strongly",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"strength",
				"power",
				"energy",
			},
		},
	},
	"uta (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to osculate",
				"to oral stimulate",
				"to suck",
				"to kiss",
			},
		},
	},
	"meli": {
		{
			POS: "adjective",
			Meanings: []string{
				"feminine",
				"womanly",
				"female",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"female",
				"girl",
				"wife",
				"girlfriend",
				"woman",
			},
		},
	},
	"ante": {
		{
			POS: "adjective",
			Meanings: []string{
				"dissimilar",
				"changed",
				"other",
				"unequal",
				"differential",
				"different",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"distinction",
				"differential",
				"variation",
				"variance",
				"disagreement",
				"difference",
			},
		},
	},
	"utala (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to strike",
				"to attack",
				"to compete against",
				"to hit",
			},
		},
	},
	"jan (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to humanize",
				"to personalize",
				"to personify",
			},
		},
	},
	"pu (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to apply (the official Toki Pona book) to",
			},
		},
	},
	"poki (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to put in",
				"to can",
				"to bottle",
				"to box up",
			},
		},
	},
	"len (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to be dressed",
				"to dress",
				"to wear",
			},
		},
	},
	"lon": {
		{
			POS: "adjective",
			Meanings: []string{
				"existing",
				"correct",
				"real",
				"genuine",
				"true",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"being",
				"presence",
				"existence",
			},
		},
		{
			POS: "preposition",
			Meanings: []string{
				"be (located) in/at/on",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to be present",
				"to be real/true",
				"to exist",
				"to be there",
			},
		},
	},
	"sona": {
		{
			POS: "adjective",
			Meanings: []string{
				"cognizant",
				"shrewd",
				"knowing",
			},
		},
		{
			POS: "auxiliary verb",
			Meanings: []string{
				"to know how to",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"wisdom",
				"intelligence",
				"understanding",
				"knowledge",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to understand",
				"to know",
			},
		},
	},
	"moli": {
		{
			POS: "adjective",
			Meanings: []string{
				"dying",
				"fatal",
				"deadly",
				"lethal",
				"mortal",
				"deathly",
				"killing",
				"dead",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"mortally",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"decease",
				"death",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to be dead",
				"to die",
			},
		},
	},
	"esun": {
		{
			POS: "adjective",
			Meanings: []string{
				"trade",
				"marketable",
				"for sale",
				"salable",
				"deductible",
				"commercial",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"shop",
				"fair",
				"bazaar",
				"business",
				"transaction",
				"market",
			},
		},
	},
	"kalama (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to ring",
				"to play (an instrument)",
				"to sound",
			},
		},
	},
	"anu": {
		{
			POS: "conjunction",
			Meanings: []string{
				"or (used for decision questions)",
			},
		},
	},
	"wawa (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to energize",
				"to empower",
				"to strengthen",
			},
		},
	},
	"suwi (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to sweeten",
			},
		},
	},
	"moku (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to drink",
				"to swallow",
				"to ingest",
				"to consume",
				"to eat",
			},
		},
	},
	"namako (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to flavor",
				"to decorate",
				"to spice",
			},
		},
	},
	"pakala": {
		{
			POS: "adjective",
			Meanings: []string{
				"ruined",
				"demolished",
				"shattered",
				"wrecked",
				"destroyed",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"ruined",
				"demolished",
				"shattered",
				"wrecked",
				"destroyed",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"accident",
				"mistake",
				"destruction",
				"damage",
				"breaking",
				"blunder",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to fall apart",
				"to break",
				"to screw up",
			},
		},
	},
	"pan": {
		{
			POS: "noun",
			Meanings: []string{
				"grain; barley",
				"corn",
				"oat",
				"rice",
				"wheat; bread",
				"pasta",
				"cereal",
			},
		},
	},
	"mani": {
		{
			POS: "adjective",
			Meanings: []string{
				"financially",
				"monetary",
				"pecuniary",
				"financial",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"financially",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"material wealth",
				"currency",
				"dollar",
				"capital",
				"money",
			},
		},
	},
	"toki!": {
		{
			POS: "interjection",
			Meanings: []string{
				"hi",
				"good morning",
				"",
				"hello",
			},
		},
	},
	"mute": {
		{
			POS: "adjective",
			Meanings: []string{
				"very",
				"much",
				"several",
				"a lot",
				"abundant",
				"numerous",
				"more",
				"many",
			},
		},
		{
			POS: "adjective numeral",
			Meanings: []string{
				"20 (official Toki Pona book)",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"very",
				"much",
				"several",
				"a lot",
				"abundant",
				"numerous",
				"more",
				"many",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"quantity",
				"amount",
			},
		},
	},
	"olin": {
		{
			POS: "adjective",
			Meanings: []string{
				"love",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"love",
			},
		},
	},
	"ala!": {
		{
			POS: "interjection",
			Meanings: []string{
				"no!",
			},
		},
	},
	"lon (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to create",
				"to give birth",
			},
		},
	},
	"tawa": {
		{
			POS: "adjective",
			Meanings: []string{
				"mobile",
				"moving",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"mobile",
				"moving",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"transportation",
				"movement",
			},
		},
		{
			POS: "preposition",
			Meanings: []string{
				"in order to",
				"towards",
				"for",
				"until",
				"to",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to travel",
				"to move",
				"to leave",
				"to visit",
				"to walk",
			},
		},
	},
	"seli (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to warm up",
				"to cook",
				"to heat",
			},
		},
	},
	"lawa (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to control",
				"to rule",
				"to steer",
				"to lead",
			},
		},
	},
	"pilin": {
		{
			POS: "adjective",
			Meanings: []string{
				"feeling",
				"empathic",
				"sensitive",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"perceptively",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"emotion",
				"feel",
				"think",
				"sense",
				"touch",
				"",
				"feelings",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to sense",
				"to feel",
			},
		},
	},
	"walo (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to whitewash",
				"to whiten",
			},
		},
	},
	"pini (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to stop",
				"to turn off",
				"to finish",
				"to close",
				"to end",
			},
		},
	},
	"pana (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to put",
				"to send",
				"to place",
				"to release",
				"to emit",
				"to cause",
				"to give",
			},
		},
	},
	"ilo": {
		{
			POS: "adjective",
			Meanings: []string{
				"useful",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"usefully",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"device",
				"machine",
				"thing used for a specific purpose",
				"tool",
			},
		},
	},
	"lete (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to chill",
				"to cool down",
			},
		},
	},
	"tu (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to double",
				"to separate",
				"to cut in two",
				"to divide",
			},
		},
	},
	"jelo": {
		{
			POS: "adjective",
			Meanings: []string{
				"yellowy",
				"yellowish",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"light green",
				"yellow",
			},
		},
	},
	"awen (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to keep",
			},
		},
	},
	"kulupu": {
		{
			POS: "adjective",
			Meanings: []string{
				"shared",
				"public",
				"of the society",
				"communal",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"community",
				"society",
				"company",
				"people",
				"group",
			},
		},
	},
	"pona": {
		{
			POS: "adjective",
			Meanings: []string{
				"simple",
				"positive",
				"nice",
				"correct",
				"right",
				"good",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"simple",
				"positive",
				"nice",
				"correct",
				"right",
				"good",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"simplicity",
				"positivity",
				"good",
			},
		},
	},
	"monsi": {
		{
			POS: "adjective",
			Meanings: []string{
				"rear",
				"back",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"rear end",
				"butt",
				"behind",
				"back",
			},
		},
	},
	"supa": {
		{
			POS: "adjective",
			Meanings: []string{
				"shallow",
				"flat-bottomed",
				"horizontal",
				"flat",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"e.g furniture",
				"table",
				"chair",
				"pillow",
				"floor",
				"horizontal surface",
			},
		},
	},
	"sama la": {
		{
			POS: "noun",
			Meanings: []string{
				"if parity",
				"on identity",
				"in case of equality",
			},
		},
	},
	"kute": {
		{
			POS: "adjective",
			Meanings: []string{
				"hearing",
				"auditory",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"ear",
				"hearing",
			},
		},
	},
	"mama (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to wet-nurse",
				"mothering",
				"to mother sb.",
			},
		},
	},
	"utala": {
		{
			POS: "adjective",
			Meanings: []string{
				"fighting",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"fighting",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"disharmony",
				"fight",
				"war",
				"battle",
				"attack",
				"violence",
				"conflict",
			},
		},
	},
	"sina": {
		{
			POS: "personal pronoun",
			Meanings: []string{
				"you",
			},
		},
		{
			POS: "possessive pronoun",
			Meanings: []string{
				"yours",
			},
		},
	},
	"pi": {
		{
			POS: "separator",
			Meanings: []string{
				"'pi' is used to build complex compound nouns. 'pi' separates a (pro)noun from another (pro)noun that has at least one adjective. After 'pi' could only be a noun or pronoun.",
			},
		},
	},
	"ona": {
		{
			POS: "personal pronoun",
			Meanings: []string{
				"he",
				"it",
				"they",
				"she",
			},
		},
		{
			POS: "possessive pronoun",
			Meanings: []string{
				"his",
				"its",
				"her",
			},
		},
	},
	"nena": {
		{
			POS: "adjective",
			Meanings: []string{
				"undulating",
				"mountainous",
				"hunchbacked",
				"humpbacked",
				"bumpy",
				"hilly",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"hill",
				"extrusion",
				"button",
				"mountain",
				"nose",
				"protuberance",
				"bump",
			},
		},
	},
	"kon": {
		{
			POS: "adjective",
			Meanings: []string{
				"ethereal",
				"gaseous",
				"air-like",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"ethereal",
				"gaseous",
				"air-like",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"wind",
				"smell",
				"soul",
				"air",
			},
		},
	},
	"unpa (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to sleep with",
				"to fuck",
				"to have sex with",
			},
		},
	},
	"toki (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to talk",
				"to say",
				"to pronounce",
				"to discourse",
				"to speak",
			},
		},
	},
	"sona (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to understand",
				"to know how to",
				"to know",
			},
		},
	},
	"jan": {
		{
			POS: "adjective",
			Meanings: []string{
				"somebody's",
				"personal",
				"of people",
				"human",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"somebody's",
				"personal",
				"of people",
				"human",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"people",
				"human",
				"being",
				"somebody",
				"anybody",
				"person",
			},
		},
	},
	"kama sona (e )": {
		{
			POS: "transitives verb",
			Meanings: []string{
				"to study",
				"to learn",
			},
		},
	},
	"ken la": {
		{
			POS: "noun",
			Meanings: []string{
				"if ability",
				"if permission",
				"if possibility",
			},
		},
	},
	"kin!": {
		{
			POS: "interjection",
			Meanings: []string{
				"really!",
			},
		},
	},
	"lete": {
		{
			POS: "adjective",
			Meanings: []string{
				"cool",
				"uncooked",
				"raw",
				"perishing",
				"cold",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"bleakly",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"chill",
				"bleakness",
				"cold",
			},
		},
	},
	"lili": {
		{
			POS: "adjective",
			Meanings: []string{
				"little",
				"young",
				"a bit",
				"short",
				"few",
				"less",
				"small",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"little",
				"young",
				"a bit",
				"short",
				"few",
				"less",
				"small",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"youth",
				"immaturity",
				"smallness",
			},
		},
	},
	"sama (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to make equal",
				"to make similar to",
				"to equate",
			},
		},
	},
	"lape (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to knock out",
			},
		},
	},
	"sina kama e ni": {},
	"pu": {
		{
			POS: "adjective",
			Meanings: []string{
				"buying and interacting with the official Toki Pona book",
			},
		},
		{
			POS: "auxiliary verb",
			Meanings: []string{
				"to buying and interacting with the official Toki Pona book",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"buying and interacting with the official Toki Pona book",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"to buy and to read (the official Toki Pona book)",
			},
		},
	},
	"pilin (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to think",
				"to touch",
				"to fumble",
				"to fiddle",
				"to feel",
			},
		},
	},
	"kin": {
		{
			POS: "adjective",
			Meanings: []string{
				"still",
				"too kin can be the very last word in an adjective group.",
				"indeed",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"indeed",
				"in fact",
				"really",
				"objectively",
				"kin can be the very last word in an adverb group.",
				"actually",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"fact",
				"reality",
			},
		},
	},
	"tan": {
		{
			POS: "adjective",
			Meanings: []string{
				"",
				"causal",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"cause",
				"origin",
			},
		},
		{
			POS: "preposition",
			Meanings: []string{
				"by",
				"because of",
				"since",
				"from",
			},
		},
		{
			POS: "verb intransitive",
			Meanings: []string{
				"originate from",
				"come out of",
				"to come from",
			},
		},
	},
	"kili": {
		{
			POS: "adjective",
			Meanings: []string{
				"fruity",
			},
		},
		{
			POS: "adverb",
			Meanings: []string{
				"fruity",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"pulpy vegetable",
				"mushroom",
				"fruit",
			},
		},
	},
	"suwi": {
		{
			POS: "adjective",
			Meanings: []string{
				"cute",
				"sweet",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"sweet food",
				"candy",
			},
		},
	},
	"jo": {
		{
			POS: "adjective",
			Meanings: []string{
				"personal",
				"private",
			},
		},
		{
			POS: "noun",
			Meanings: []string{
				"possessions",
				"content",
				"having",
			},
		},
	},
	"nimi": {
		{
			POS: "noun",
			Meanings: []string{
				"name",
				"word",
			},
		},
	},
	"lukin (e )": {
		{
			POS: "verb transitive",
			Meanings: []string{
				"to look at",
				"to watch",
				"to read",
				"to see",
			},
		},
	},
}

// Descr describes this plugin
func (t *Toki) Descr() string {
	return "Toki Pona dictionary"
}

// Re is the regex for matching hi messages.
func (t *Toki) Re() string {
	return `(?i)^(toki[\?]?):? (.+)$`
}

// Match determines if we are highfiving
func (t *Toki) Match(_, msg string) bool {
	re := regexp.MustCompile(t.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here
func (t *Toki) SetStore(_ PluginStore) {}

func (t *Toki) fix(msg string) (string, string) {
	re := regexp.MustCompile(t.Re())
	return re.ReplaceAllString(msg, "$1"), re.ReplaceAllString(msg, "$2")
}

// Process does the heavy lifting
func (t *Toki) Process(from, post string) string {
	cmd, w := t.fix(post)
	cmd = strings.ToLower(cmd)
	switch cmd {
	case "toki":
		if word, ok := TokiLang[w]; ok {
			var defs []string
			for _, v := range word {
				defs = append(defs, v.Print(w))
			}
			return strings.Join(defs, "\n\n")
		} else {
			return "mi sona ala"
		}
	case "toki?":
		st := stemmer.Stem(w)
		var words []string
		for i, ts := range TokiLang {
			for _, t := range ts {
				stems := stemmer.StemMultiple(t.Words())
				for _, x := range stems {
					if x == st {
						words = append(words, t.Print(i))
					}
				}
			}
		}
		return strings.Join(words, "\n\n")
	}
	return "mi sona ala"
}

// RespondText to hi events
func (t *Toki) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, post string) error {
	return SendMD(c, ev.RoomID, t.Process(ev.Sender, post))
}

// Name hi
func (t *Toki) Name() string {
	return "Toki"
}
