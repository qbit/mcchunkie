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
	result := []string{}

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
	"telo (e )": []Toki{
		Toki{
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
	"nanpa (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to reckon",
				"to number",
				"to count",
			},
		},
	},
	"kasi (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to grow",
				"to plant",
			},
		},
	},
	"ken (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to enable",
				"to allow",
				"to permit",
				"to make possible",
			},
		},
	},
	"kiwen": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"solid",
				"stone-like",
				"made of stone or metal",
				"hard",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"solid",
				"stone-like",
				"made of stone or metal",
				"hard",
			},
		},
		Toki{
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
	"weka": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"away",
				"ignored",
				"absent",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"absence",
			},
		},
	},
	"open (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to start",
				"to begin",
				"to turn on",
				"to open",
			},
		},
	},
	"kama (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to summon",
				"to bring about",
			},
		},
	},
	"walo": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"whitish",
				"light-coloured",
				"pale",
				"white",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"whiteness",
				"lightness",
				"white thing or part",
			},
		},
	},
	"mi ' pona, tan ni": []Toki{},
	"anpa": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"lower",
				"bottom",
				"down",
				"low",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"below",
				"deep",
				"low",
				"deeply",
				"downstairs",
			},
		},
		Toki{
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
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to prostrate oneself",
			},
		},
	},
	"lukin": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"visual(ly)",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"visual(ly)",
			},
		},
		Toki{
			POS: "auxiliary verb",
			Meanings: []string{
				"try to",
				"look for",
				"to seek to",
			},
		},
		Toki{
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
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to watch out",
				"to pay attention",
				"to look",
			},
		},
	},
	"pali (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to make",
				"to build",
				"to create",
				"to do",
			},
		},
	},
	"musi (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to entertain",
				"to amuse",
			},
		},
	},
	"mu!": []Toki{
		Toki{
			POS: "interjection",
			Meanings: []string{
				"woof! meow! moo! etc. (cute animal noise)",
			},
		},
	},
	"weka (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to eliminate",
				"to throw away",
				"to get rid of",
				"to remove",
			},
		},
	},
	"namako": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"piquant",
				"spicy",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"something extra",
				"food additive",
				"accessory",
				"spice",
			},
		},
	},
	"pini": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"finished",
				"past",
				"done",
				"completed",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"past",
				"perfectly",
				"ago",
			},
		},
		Toki{
			POS: "auxiliary verb",
			Meanings: []string{
				"to finish",
				"to end",
				"to interrupt",
				"to stop",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"tip",
				"end",
			},
		},
	},
	"pakala!": []Toki{
		Toki{
			POS: "interjection",
			Meanings: []string{
				"damn! fuck!",
			},
		},
	},
	"kama moli": []Toki{
		Toki{
			POS: "intransitives verb",
			Meanings: []string{
				"dieing",
			},
		},
	},
	"mi moku, tan ni": []Toki{},
	"moli (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to kill",
			},
		},
	},
	"mi wile e ni": []Toki{},
	"kalama": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"loud",
				"rowdy",
				"noisy",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"noise",
				"voice",
				"sound",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to make noise",
			},
		},
	},
	"linja": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"oblong",
				"long",
				"elongated",
			},
		},
		Toki{
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
	"lape": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"of sleep",
				"dormant",
				"sleeping",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"asleep",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"rest",
				"sleep",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to rest",
				"to sleep",
			},
		},
	},
	"tenpo": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"chronological",
				"chronologic",
				"temporal",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"chronologically",
			},
		},
		Toki{
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
	"sewi": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"elevated",
				"religious",
				"formal",
				"superior",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"elevated",
				"religious",
				"formal",
				"superior",
			},
		},
		Toki{
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
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to get up",
			},
		},
	},
	"kon (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to puff away something",
				"to blow away something",
			},
		},
	},
	"waso": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"bird-",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"bat; flying creature",
				"winged animal",
				"bird",
			},
		},
	},
	"sitelen": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"pictorial",
				"metaphorical",
				"metaphorisch",
				"figurative",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"pictorially",
			},
		},
		Toki{
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
	"sin (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to renovate",
				"to freshen",
				"to renew",
			},
		},
	},
	"sike (e )": []Toki{
		Toki{
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
	"unpa": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"sexual",
				"erotic",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"sexual",
				"erotic",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"sexuality",
				"sex",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to have sex",
			},
		},
	},
	"sijelo (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to heal up",
				"to cure",
				"to heal",
			},
		},
	},
	"palisa (e )": []Toki{
		Toki{
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
	"pakala (e )": []Toki{
		Toki{
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
	"alasa (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to forage",
				"to hunt",
			},
		},
	},
	"insa": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"internal",
				"inner",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"inner world",
				"centre",
				"stomach",
				"inside",
			},
		},
	},
	"ko": []Toki{
		Toki{
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
	"len": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"clothed",
				"costumed",
				"dressed up",
				"dressed",
			},
		},
		Toki{
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
	"lawa": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"leading",
				"in charge",
				"main",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"leading",
				"in charge",
				"main",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"mind",
				"head",
			},
		},
	},
	"sitelen (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to write",
				"to draw",
			},
		},
	},
	"!": []Toki{
		Toki{
			POS: "separator",
			Meanings: []string{
				"'.",
			},
		},
	},
	"\"": []Toki{
		Toki{
			POS: "separator",
			Meanings: []string{
				"Quotation marks are used for words with original spelling or for quotes.",
			},
		},
	},
	"lupa": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"holey",
				"full of holes",
				"hole-",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"orifice",
				"door",
				"window",
				"hole",
			},
		},
	},
	"#": []Toki{
		Toki{
			POS: "unofficial",
			Meanings: []string{
				"Number sign",
			},
		},
	},
	"kama jo (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to get",
			},
		},
	},
	"sijelo": []Toki{
		Toki{
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
		Toki{
			POS: "adverb",
			Meanings: []string{
				"bodily",
				"physically",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"physical state",
				"torso",
				"body (of person or animal)",
			},
		},
	},
	"pimeja (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to darken",
			},
		},
	},
	"ona li wile e ni": []Toki{},
	"a a a!": []Toki{
		Toki{
			POS: "interjection",
			Meanings: []string{
				"laugh",
			},
		},
	},
	"kulupu (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to call together",
				"to convene",
				"to assemble",
			},
		},
	},
	"'": []Toki{
		Toki{
			POS: "unofficial",
			Meanings: []string{
				"An apostrophe can identify a predicate that does not contain a verb.",
			},
		},
	},
	"tawa (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to displace",
				"to move",
			},
		},
	},
	"soweli": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"animal",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"especially land mammal",
				"lovable animal",
				"beast",
				"animal",
			},
		},
	},
	"en": []Toki{
		Toki{
			POS: "conjunction",
			Meanings: []string{
				"and (used to coordinate head nouns)",
			},
		},
	},
	"jo (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to contain",
				"to have",
			},
		},
	},
	"wile (e )": []Toki{
		Toki{
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
	",": []Toki{
		Toki{
			POS: "separator",
			Meanings: []string{
				"A comma is used after an 'o' to addressing people. Optional you can put a comma before a preposition. Don't use a comma before or after",
			},
		},
	},
	"palisa": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"long",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"rod",
				"stick",
				"pointy thing",
				"long hard thing; branch",
			},
		},
	},
	"alasa": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"-hunting",
				"hunting",
				"hunting-",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"hunting",
			},
		},
	},
	"la": []Toki{
		Toki{
			POS: "separator",
			Meanings: []string{
				"half sentence or noun. Don't use 'la' before or after",
				"A 'la' is between a conditional phrases and the main sentence. A context phrase can be sentence",
			},
		},
	},
	".": []Toki{
		Toki{
			POS: "separator",
			Meanings: []string{
				"'.",
			},
		},
	},
	"ike!": []Toki{
		Toki{
			POS: "interjection",
			Meanings: []string{
				"oh dear! woe! alas!",
			},
		},
	},
	"suli (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to lengthen",
				"to enlarge",
			},
		},
	},
	"tomo (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to construct",
				"to engineer",
				"to build",
			},
		},
	},
	"toki": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"eloquent",
				"linguistic",
				"verbal",
				"grammatical",
				"speaking",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"eloquent",
				"linguistic",
				"verbal",
				"grammatical",
				"speaking",
			},
		},
		Toki{
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
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to chat",
				"to communicate",
				"to talk",
			},
		},
	},
	"taso": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"sole",
				"only",
			},
		},
		Toki{
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
		Toki{
			POS: "conjunction",
			Meanings: []string{
				"however",
				"but",
			},
		},
	},
	"li": []Toki{
		Toki{
			POS: "separator",
			Meanings: []string{
				"'",
				"'.",
				"'",
			},
		},
	},
	"suli": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"tall",
				"long",
				"adult",
				"important",
				"big",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"tall",
				"long",
				"adult",
				"important",
				"big",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"size",
			},
		},
	},
	"selo mi li wile e ni": []Toki{},
	"pan (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to sow",
			},
		},
	},
	"sewi (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to lift",
			},
		},
	},
	"sama": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"similar",
				"equal",
				"of equal status or position",
				"same",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"equally",
				"exactly the same",
				"just the same",
				"similarly",
				"just as",
			},
		},
		Toki{
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
		Toki{
			POS: "preposition",
			Meanings: []string{
				"as",
				"seem",
				"like",
			},
		},
	},
	"pona la": []Toki{
		Toki{
			POS: "noun",
			Meanings: []string{
				"if simplicity",
				"if positivity",
				"if good",
			},
		},
	},
	"ike (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to worsen",
				"to make bad",
			},
		},
	},
	"kule (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to color",
				"to paint",
			},
		},
	},
	"lili (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to shorten",
				"to shrink",
				"to lessen",
				"to reduce",
			},
		},
	},
	"pali": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"work-related",
				"operating",
				"working",
				"active",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"briskly",
				"actively",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"work",
				"deed",
				"project",
				"activity",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to work",
				"to function",
				"to act",
			},
		},
	},
	"ala": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"not",
				"none",
				"un-",
				"no",
			},
		},
		Toki{
			POS: "adjective numeral",
			Meanings: []string{
				"0",
				"null",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"don't",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"negation",
				"zero",
				"nothing",
			},
		},
	},
	"?": []Toki{
		Toki{
			POS: "separator",
			Meanings: []string{
				"'.",
			},
		},
	},
	"selo (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to protect",
				"to guard",
				"to shelter",
			},
		},
	},
	"ale": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"every",
				"complete",
				"whole (ale = ali)",
				"(depreciated)",
				"all",
			},
		},
		Toki{
			POS: "adjective numeral",
			Meanings: []string{
				"100 (official Toki Pona book)",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"forever",
				"evermore",
				"eternally (ale = ali)",
				"(depreciated)",
				"always",
			},
		},
		Toki{
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
	"jaki!": []Toki{
		Toki{
			POS: "interjection",
			Meanings: []string{
				"ew! yuck!",
			},
		},
	},
	"ken": []Toki{
		Toki{
			POS: "auxiliary verb",
			Meanings: []string{
				"may",
				"to can",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"ability",
				"power to do things",
				"permission",
				"possibility",
			},
		},
		Toki{
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
	"kin la": []Toki{
		Toki{
			POS: "noun",
			Meanings: []string{
				"if fact",
				"if reality",
			},
		},
	},
	"ali": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"every",
				"complete",
				"whole (ale = ali)",
				"all",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"forever",
				"evermore",
				"eternally (ale = ali)",
				"always",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"anything",
				"life",
				"the universe",
				"everything",
			},
		},
	},
	"ante la": []Toki{
		Toki{
			POS: "noun",
			Meanings: []string{
				"if variance",
				"if disagreement",
				"if difference",
			},
		},
	},
	"esun (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to sell",
				"to barter",
				"to swap",
				"to buy",
			},
		},
	},
	"anpa (e )": []Toki{
		Toki{
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
	"pipi": []Toki{
		Toki{
			POS: "noun",
			Meanings: []string{
				"insect",
				"spider",
				"bug",
			},
		},
	},
	"open": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"starting",
				"opening",
				"initial",
			},
		},
		Toki{
			POS: "auxiliary verb",
			Meanings: []string{
				"to start",
				"to begin",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"beginning",
				"opening",
				"start",
			},
		},
	},
	"o!": []Toki{
		Toki{
			POS: "interjection",
			Meanings: []string{
				"hey! (calling somebody's attention)",
			},
		},
	},
	"wan": []Toki{
		Toki{
			POS: "adjective numeral",
			Meanings: []string{
				"1",
				"one",
			},
		},
		Toki{
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
	"telo": []Toki{
		Toki{
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
		Toki{
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
		Toki{
			POS: "noun",
			Meanings: []string{
				"liquid",
				"juice",
				"sauce",
				"water",
			},
		},
	},
	"pona (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to fix",
				"to repair",
				"to make good",
				"to improve",
			},
		},
	},
	"ma": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"outdoor",
				"alfresco",
				"open-air",
				"countrified",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"earth",
				"country",
				"(outdoor) area",
				"land",
			},
		},
	},
	"sinpin": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"frontal",
				"anterior",
				"vertical",
				"facial",
			},
		},
		Toki{
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
	"poka": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"neighbouring",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"hip",
				"next to",
				"side",
			},
		},
	},
	"seli": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"warm",
				"cooked",
				"hot",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"warm",
				"cooked",
				"hot",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"warmth",
				"heat",
				"fire",
			},
		},
	},
	"luka": []Toki{
		Toki{
			POS: "adjective numeral",
			Meanings: []string{
				"5",
				"five",
			},
		},
		Toki{
			POS: "adjective",
			Meanings: []string{
				"palpable",
				"tangible",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"hand",
				"tacticle organ",
				"arm",
			},
		},
	},
	"sin": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"fresh",
				"another",
				"more",
				"new",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"regenerative",
			},
		},
		Toki{
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
	"pimeja": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"dark",
				"black",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"shadows",
				"darkness",
			},
		},
	},
	"wile": []Toki{
		Toki{
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
		Toki{
			POS: "noun",
			Meanings: []string{
				"need",
				"will",
				"desire",
			},
		},
	},
	"olin (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to love (a person)",
			},
		},
	},
	"mi": []Toki{
		Toki{
			POS: "personal pronoun",
			Meanings: []string{
				"we",
				"I",
			},
		},
		Toki{
			POS: "possessive pronoun",
			Meanings: []string{
				"our",
				"my",
			},
		},
	},
	"selo": []Toki{
		Toki{
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
	"poki": []Toki{
		Toki{
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
	"o,": []Toki{
		Toki{
			POS: "interjection",
			Meanings: []string{
				"adressing people",
			},
		},
	},
	"mute (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to make many or much",
			},
		},
	},
	"jaki": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"gross",
				"filthy",
				"obscene",
				"dirty",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"gross",
				"filthy",
				"dirty",
			},
		},
		Toki{
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
	"mun": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"lunar",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"lunar",
				"night sky object",
				"star",
				"moon",
			},
		},
	},
	"loje": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"ruddy",
				"pink",
				"pinkish",
				"gingery",
				"reddish",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"red",
			},
		},
	},
	"sike": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"cyclical",
				"of one year",
				"round",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"rotated",
			},
		},
		Toki{
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
	"ijo (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to objectify",
			},
		},
	},
	"nasa": []Toki{
		Toki{
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
		Toki{
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
		Toki{
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
	"mi pilin e ni": []Toki{},
	"ike la": []Toki{
		Toki{
			POS: "noun",
			Meanings: []string{
				"if badness",
				"if evil",
				"if negativity",
			},
		},
	},
	"kiwen (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to harden",
				"to petrify",
				"to fossilize",
				"to solidify",
			},
		},
	},
	"mu (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to make animal noise",
			},
		},
	},
	"noka": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"lower",
				"bottom",
				"foot-",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"on foot",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"foot; organ of locomotion; bottom",
				"lower part",
				"leg",
			},
		},
	},
	"o !": []Toki{
		Toki{
			POS: "separator",
			Meanings: []string{
				"'o' replace 'li'.",
			},
		},
		Toki{
			POS: "subject",
			Meanings: []string{
				"An 'o' is used for imperative (commands). 'o' replace the subject.",
			},
		},
	},
	"mu": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"animal nois-",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"animal nois-",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"animal noise",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to communicate animally",
			},
		},
	},
	"a": []Toki{
		Toki{
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
	"oko": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"eye-",
				"optical",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"eye",
			},
		},
	},
	"kala": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"fish-",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"marine animal",
				"sea creature",
				"fish",
			},
		},
	},
	"e sina": []Toki{
		Toki{
			POS: "reflexive pronoun",
			Meanings: []string{
				"yourselves",
				"yourself",
			},
		},
	},
	"nasa (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to make weird",
				"to drive crazy",
			},
		},
	},
	"e": []Toki{
		Toki{
			POS: "separator",
			Meanings: []string{
				"'",
				"'.",
				"'",
			},
		},
	},
	"ijo": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"of something",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"of something",
			},
		},
		Toki{
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
	"pona!": []Toki{
		Toki{
			POS: "interjection",
			Meanings: []string{
				"great! good! thanks! OK! cool! yay!",
			},
		},
	},
	"ante (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to alter",
				"to modify",
				"to change",
			},
		},
	},
	"akesi": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"reptilian-",
				"slimy",
				"amphibian-",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"amphibian; non-cute animal",
				"reptile",
			},
		},
	},
	"seme": []Toki{
		Toki{
			POS: "question pronoun",
			Meanings: []string{
				"which",
				"wh- (question word)",
				"what",
			},
		},
	},
	"nimi (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to name",
			},
		},
	},
	"e ona": []Toki{
		Toki{
			POS: "reflexive pronoun",
			Meanings: []string{
				"herself",
				"itself",
				"themselves",
				"himself",
			},
		},
	},
	"mije": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"masculine",
				"manly",
				"male",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"male",
				"husband",
				"boyfriend",
				"man",
			},
		},
	},
	"mama": []Toki{
		Toki{
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
		Toki{
			POS: "noun",
			Meanings: []string{
				"mother",
				"father",
				"parent",
			},
		},
	},
	"tu": []Toki{
		Toki{
			POS: "adjective numeral",
			Meanings: []string{
				"2",
				"two",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"pair",
				"duo",
			},
		},
	},
	"jaki (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to dirty",
				"to pollute",
			},
		},
	},
	"wan (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to make one",
				"to unite",
			},
		},
	},
	"suno (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to illumine",
				"to light",
			},
		},
	},
	"ni": []Toki{
		Toki{
			POS: "adjective demonstrative pronoun",
			Meanings: []string{
				"that",
				"this",
			},
		},
		Toki{
			POS: "noun demonstrative pronoun",
			Meanings: []string{
				"that",
				"this",
			},
		},
	},
	"kute (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to listen",
				"",
				"to hear",
			},
		},
	},
	"pana": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"generous",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"transfer",
				"exchange",
				"giving",
			},
		},
	},
	"nanpa": []Toki{
		Toki{
			POS: "adjective numeral",
			Meanings: []string{
				"To build ordinal numbers.",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"numeral",
				"number",
			},
		},
	},
	"lupa (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to stab",
				"to perforate",
				"to pierce",
			},
		},
	},
	"tomo": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"domestic",
				"household",
				"urban",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"domestic",
				"household",
				"urban",
			},
		},
		Toki{
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
	"nasin": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"habitual",
				"customary",
				"doctrinal",
				"systematic",
			},
		},
		Toki{
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
	"kepeken": []Toki{
		Toki{
			POS: "noun",
			Meanings: []string{
				"usage",
				"tool",
				"use",
			},
		},
		Toki{
			POS: "preposition",
			Meanings: []string{
				"using",
				"with",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to use",
			},
		},
	},
	"laso": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"bluey",
				"bluish",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"blue-green",
				"blue",
			},
		},
	},
	"ko (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to pulverize",
				"to squash",
			},
		},
	},
	"lipu": []Toki{
		Toki{
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
		Toki{
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
	"suno": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"sunnily",
				"sunny",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"sunnily",
				"sunny",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"light",
				"sun",
			},
		},
	},
	"sina toki e ni, tawa mi": []Toki{},
	"open la": []Toki{
		Toki{
			POS: "noun",
			Meanings: []string{
				"in the beginning",
				"at the opening",
			},
		},
	},
	"ike": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"negative",
				"wrong",
				"evil",
				"overly complex",
				"bad",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"negative",
				"wrong",
				"evil",
				"overly complex",
				"bad",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"badness",
				"evil",
				"negativity",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to suck",
				"to be bad",
			},
		},
	},
	"kule": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"pigmented",
				"painted",
				"colourful",
			},
		},
		Toki{
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
	"moku": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"eating",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"eating",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"meal",
				"food",
			},
		},
	},
	"kasi": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"vegetal",
				"biological",
				"biologic",
				"leafy",
				"vegetable",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"vegetation",
				"herb",
				"leaf",
				"plant",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to grow",
			},
		},
	},
	"musi": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"fun",
				"recreational",
				"artful",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"cheerfully",
			},
		},
		Toki{
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
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to have fun",
				"to play",
			},
		},
	},
	"awen": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"stationary",
				"permanent",
				"sedentary",
				"remaining",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"yet",
				"still",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"continuity",
				"continuum",
				"stay",
				"inertia",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to wait",
				"to remain",
				"to stay",
			},
		},
	},
	"e mi": []Toki{
		Toki{
			POS: "reflexive pronoun",
			Meanings: []string{
				"ourselves",
				"myself",
			},
		},
	},
	"uta": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"oral",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"orally",
			},
		},
		Toki{
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
	"kama": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"future",
				"coming",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"future",
				"coming",
			},
		},
		Toki{
			POS: "auxiliary verb",
			Meanings: []string{
				"to mange to",
				"to become",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"happening",
				"chance",
				"arrival",
				"beginning",
				"event",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to become",
				"to arrive",
				"to happen",
				"to come",
			},
		},
	},
	"wawa": []Toki{
		Toki{
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
		Toki{
			POS: "adverb",
			Meanings: []string{
				"powerfully",
				"strongly",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"strength",
				"power",
				"energy",
			},
		},
	},
	"uta (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to osculate",
				"to oral stimulate",
				"to suck",
				"to kiss",
			},
		},
	},
	"meli": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"feminine",
				"womanly",
				"female",
			},
		},
		Toki{
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
	"ante": []Toki{
		Toki{
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
		Toki{
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
	"utala (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to strike",
				"to attack",
				"to compete against",
				"to hit",
			},
		},
	},
	"jan (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to humanize",
				"to personalize",
				"to personify",
			},
		},
	},
	"pu (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to apply (the official Toki Pona book) to",
			},
		},
	},
	"poki (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to put in",
				"to can",
				"to bottle",
				"to box up",
			},
		},
	},
	"len (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to be dressed",
				"to dress",
				"to wear",
			},
		},
	},
	"lon": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"existing",
				"correct",
				"real",
				"genuine",
				"true",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"being",
				"presence",
				"existence",
			},
		},
		Toki{
			POS: "preposition",
			Meanings: []string{
				"be (located) in/at/on",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to be present",
				"to be real/true",
				"to exist",
				"to be there",
			},
		},
	},
	"sona": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"cognizant",
				"shrewd",
				"knowing",
			},
		},
		Toki{
			POS: "auxiliary verb",
			Meanings: []string{
				"to know how to",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"wisdom",
				"intelligence",
				"understanding",
				"knowledge",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to understand",
				"to know",
			},
		},
	},
	"moli": []Toki{
		Toki{
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
		Toki{
			POS: "adverb",
			Meanings: []string{
				"mortally",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"decease",
				"death",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to be dead",
				"to die",
			},
		},
	},
	"esun": []Toki{
		Toki{
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
		Toki{
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
	"kalama (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to ring",
				"to play (an instrument)",
				"to sound",
			},
		},
	},
	"anu": []Toki{
		Toki{
			POS: "conjunction",
			Meanings: []string{
				"or (used for decision questions)",
			},
		},
	},
	"wawa (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to energize",
				"to empower",
				"to strengthen",
			},
		},
	},
	"suwi (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to sweeten",
			},
		},
	},
	"moku (e )": []Toki{
		Toki{
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
	"namako (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to flavor",
				"to decorate",
				"to spice",
			},
		},
	},
	"pakala": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"ruined",
				"demolished",
				"shattered",
				"wrecked",
				"destroyed",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"ruined",
				"demolished",
				"shattered",
				"wrecked",
				"destroyed",
			},
		},
		Toki{
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
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to fall apart",
				"to break",
				"to screw up",
			},
		},
	},
	"pan": []Toki{
		Toki{
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
	"mani": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"financially",
				"monetary",
				"pecuniary",
				"financial",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"financially",
			},
		},
		Toki{
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
	"toki!": []Toki{
		Toki{
			POS: "interjection",
			Meanings: []string{
				"hi",
				"good morning",
				"",
				"hello",
			},
		},
	},
	"mute": []Toki{
		Toki{
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
		Toki{
			POS: "adjective numeral",
			Meanings: []string{
				"20 (official Toki Pona book)",
			},
		},
		Toki{
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
		Toki{
			POS: "noun",
			Meanings: []string{
				"quantity",
				"amount",
			},
		},
	},
	"olin": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"love",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"love",
			},
		},
	},
	"ala!": []Toki{
		Toki{
			POS: "interjection",
			Meanings: []string{
				"no!",
			},
		},
	},
	"lon (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to create",
				"to give birth",
			},
		},
	},
	"tawa": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"mobile",
				"moving",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"mobile",
				"moving",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"transportation",
				"movement",
			},
		},
		Toki{
			POS: "preposition",
			Meanings: []string{
				"in order to",
				"towards",
				"for",
				"until",
				"to",
			},
		},
		Toki{
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
	"seli (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to warm up",
				"to cook",
				"to heat",
			},
		},
	},
	"lawa (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to control",
				"to rule",
				"to steer",
				"to lead",
			},
		},
	},
	"pilin": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"feeling",
				"empathic",
				"sensitive",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"perceptively",
			},
		},
		Toki{
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
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to sense",
				"to feel",
			},
		},
	},
	"walo (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to whitewash",
				"to whiten",
			},
		},
	},
	"pini (e )": []Toki{
		Toki{
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
	"pana (e )": []Toki{
		Toki{
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
	"ilo": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"useful",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"usefully",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"device",
				"machine",
				"thing used for a specific purpose",
				"tool",
			},
		},
	},
	"lete (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to chill",
				"to cool down",
			},
		},
	},
	"tu (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to double",
				"to separate",
				"to cut in two",
				"to divide",
			},
		},
	},
	"jelo": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"yellowy",
				"yellowish",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"light green",
				"yellow",
			},
		},
	},
	"awen (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to keep",
			},
		},
	},
	"kulupu": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"shared",
				"public",
				"of the society",
				"communal",
			},
		},
		Toki{
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
	"pona": []Toki{
		Toki{
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
		Toki{
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
		Toki{
			POS: "noun",
			Meanings: []string{
				"simplicity",
				"positivity",
				"good",
			},
		},
	},
	"monsi": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"rear",
				"back",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"rear end",
				"butt",
				"behind",
				"back",
			},
		},
	},
	"supa": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"shallow",
				"flat-bottomed",
				"horizontal",
				"flat",
			},
		},
		Toki{
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
	"sama la": []Toki{
		Toki{
			POS: "noun",
			Meanings: []string{
				"if parity",
				"on identity",
				"in case of equality",
			},
		},
	},
	"kute": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"hearing",
				"auditory",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"ear",
				"hearing",
			},
		},
	},
	"mama (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to wet-nurse",
				"mothering",
				"to mother sb.",
			},
		},
	},
	"utala": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"fighting",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"fighting",
			},
		},
		Toki{
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
	"sina": []Toki{
		Toki{
			POS: "personal pronoun",
			Meanings: []string{
				"you",
			},
		},
		Toki{
			POS: "possessive pronoun",
			Meanings: []string{
				"yours",
			},
		},
	},
	"pi": []Toki{
		Toki{
			POS: "separator",
			Meanings: []string{
				"'pi' is used to build complex compound nouns. 'pi' separates a (pro)noun from another (pro)noun that has at least one adjective. After 'pi' could only be a noun or pronoun.",
			},
		},
	},
	"ona": []Toki{
		Toki{
			POS: "personal pronoun",
			Meanings: []string{
				"he",
				"it",
				"they",
				"she",
			},
		},
		Toki{
			POS: "possessive pronoun",
			Meanings: []string{
				"his",
				"its",
				"her",
			},
		},
	},
	"nena": []Toki{
		Toki{
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
		Toki{
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
	"kon": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"ethereal",
				"gaseous",
				"air-like",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"ethereal",
				"gaseous",
				"air-like",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"wind",
				"smell",
				"soul",
				"air",
			},
		},
	},
	"unpa (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to sleep with",
				"to fuck",
				"to have sex with",
			},
		},
	},
	"toki (e )": []Toki{
		Toki{
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
	"sona (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to understand",
				"to know how to",
				"to know",
			},
		},
	},
	"jan": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"somebody's",
				"personal",
				"of people",
				"human",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"somebody's",
				"personal",
				"of people",
				"human",
			},
		},
		Toki{
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
	"kama sona (e )": []Toki{
		Toki{
			POS: "transitives verb",
			Meanings: []string{
				"to study",
				"to learn",
			},
		},
	},
	"ken la": []Toki{
		Toki{
			POS: "noun",
			Meanings: []string{
				"if ability",
				"if permission",
				"if possibility",
			},
		},
	},
	"kin!": []Toki{
		Toki{
			POS: "interjection",
			Meanings: []string{
				"really!",
			},
		},
	},
	"lete": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"cool",
				"uncooked",
				"raw",
				"perishing",
				"cold",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"bleakly",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"chill",
				"bleakness",
				"cold",
			},
		},
	},
	"lili": []Toki{
		Toki{
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
		Toki{
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
		Toki{
			POS: "noun",
			Meanings: []string{
				"youth",
				"immaturity",
				"smallness",
			},
		},
	},
	"sama (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to make equal",
				"to make similar to",
				"to equate",
			},
		},
	},
	"lape (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to knock out",
			},
		},
	},
	"sina kama e ni": []Toki{},
	"pu": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"buying and interacting with the official Toki Pona book",
			},
		},
		Toki{
			POS: "auxiliary verb",
			Meanings: []string{
				"to buying and interacting with the official Toki Pona book",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"buying and interacting with the official Toki Pona book",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"to buy and to read (the official Toki Pona book)",
			},
		},
	},
	"pilin (e )": []Toki{
		Toki{
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
	"kin": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"still",
				"too kin can be the very last word in an adjective group.",
				"indeed",
			},
		},
		Toki{
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
		Toki{
			POS: "noun",
			Meanings: []string{
				"fact",
				"reality",
			},
		},
	},
	"tan": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"",
				"causal",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"cause",
				"origin",
			},
		},
		Toki{
			POS: "preposition",
			Meanings: []string{
				"by",
				"because of",
				"since",
				"from",
			},
		},
		Toki{
			POS: "verb intransitive",
			Meanings: []string{
				"originate from",
				"come out of",
				"to come from",
			},
		},
	},
	"kili": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"fruity",
			},
		},
		Toki{
			POS: "adverb",
			Meanings: []string{
				"fruity",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"pulpy vegetable",
				"mushroom",
				"fruit",
			},
		},
	},
	"suwi": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"cute",
				"sweet",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"sweet food",
				"candy",
			},
		},
	},
	"jo": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"personal",
				"private",
			},
		},
		Toki{
			POS: "noun",
			Meanings: []string{
				"possessions",
				"content",
				"having",
			},
		},
	},
	"nimi": []Toki{
		Toki{
			POS: "noun",
			Meanings: []string{
				"name",
				"word",
			},
		},
	},
	"lukin (e )": []Toki{
		Toki{
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
func (t *Toki) Match(user, msg string) bool {
	re := regexp.MustCompile(t.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here
func (t *Toki) SetStore(s PluginStore) {}

func (t *Toki) fix(msg string) (string, string) {
	re := regexp.MustCompile(t.Re())
	return re.ReplaceAllString(msg, "$1"), re.ReplaceAllString(msg, "$2")
}

// RespondText to hi events
func (t *Toki) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	cmd, w := t.fix(post)
	switch cmd {
	case "toki":
		if word, ok := TokiLang[w]; ok {
			var defs []string
			for _, v := range word {
				defs = append(defs, v.Print(w))
			}
			SendMD(c, ev.RoomID, strings.Join(defs, "\n\n"))
		} else {
			SendText(c, ev.RoomID, "mi sona ala")
		}
	case "toki?":
		st := stemmer.Stem(w)
		words := []string{}
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
		SendMD(c, ev.RoomID, strings.Join(words, "\n\n"))
	}
}

// Name hi
func (t *Toki) Name() string {
	return "Toki"
}
