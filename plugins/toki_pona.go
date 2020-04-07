package plugins

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/matrix-org/gomatrix"
)

// Toki responds to toki pona word queries
type Toki struct {
	POS       string
	Meanings  []string
	Alt       string
	Principle string
}

// Print prints the definition
func (t *Toki) Print(w string) string {
	return fmt.Sprintf("**%s**: (_%s_) %s", w, t.POS, strings.Join(t.Meanings, ", "))
}

// TokiLang is our full representation of toki pona
var TokiLang = map[string][]Toki{
	":": []Toki{
		Toki{
			POS: "separator",
			Meanings: []string{
				"A colon is between an hint sentences and a sentences. Before and after the colon has to be complete sentences. Don't use a colon before or after",
			},
		},
	},
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
	},
	"mi ' pona, tan ni": []Toki{
		Toki{
			POS: "",
			Meanings: []string{
				"I'm okay because I'm alive.",
			},
		},
	},
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
	},
	"lukin": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"visual(ly)",
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
	"mi moku, tan ni": []Toki{
		Toki{
			POS: "",
			Meanings: []string{
				"I eat because I'm hungry.",
			},
		},
	},
	"moli (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to kill",
			},
		},
	},
	"mi wile e ni": []Toki{
		Toki{
			POS: "",
			Meanings: []string{
				"I'm at home.",
			},
		},
	},
	"kalama": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"loud",
				"rowdy",
				"noisy",
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
	},
	"pimeja (e )": []Toki{
		Toki{
			POS: "verb transitive",
			Meanings: []string{
				"to darken",
			},
		},
	},
	"ona li wile e ni": []Toki{
		Toki{
			POS: "",
			Meanings: []string{
				"They don't want people to destroy the environment.",
			},
		},
	},
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
	},
	"taso": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"sole",
				"only",
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
	},
	"selo mi li wile e ni": []Toki{
		Toki{
			POS: "",
			Meanings: []string{
				"I touch it.",
			},
		},
	},
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
	},
	"poka": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"neighbouring",
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
	},
	"luka": []Toki{
		Toki{
			POS: "adjective numeral",
			Meanings: []string{
				"5",
				"five",
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
	},
	"pimeja": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"dark",
				"black",
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
	},
	"mun": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"lunar",
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
	},
	"mi pilin e ni": []Toki{
		Toki{
			POS: "",
			Meanings: []string{
				"I think that he doesn't have money.",
			},
		},
	},
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
	},
	"o !": []Toki{
		Toki{
			POS: "separator",
			Meanings: []string{
				"'o' replace 'li'.",
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
	},
	"kala": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"fish-",
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
	},
	"tu": []Toki{
		Toki{
			POS: "adjective numeral",
			Meanings: []string{
				"2",
				"two",
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
	},
	"nanpa": []Toki{
		Toki{
			POS: "adjective numeral",
			Meanings: []string{
				"To build ordinal numbers.",
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
	},
	"laso": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"bluey",
				"bluish",
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
	},
	"suno": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"sunnily",
				"sunny",
			},
		},
	},
	"sina toki e ni, tawa mi": []Toki{
		Toki{
			POS: "",
			Meanings: []string{
				"You told me that you are eating.",
			},
		},
	},
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
	},
	"moku": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"eating",
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
	},
	"kama": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"future",
				"coming",
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
	},
	"olin": []Toki{
		Toki{
			POS: "adjective",
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
	},
	"monsi": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"rear",
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
	},
	"sina": []Toki{
		Toki{
			POS: "personal pronoun",
			Meanings: []string{
				"you",
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
	"sina kama e ni": []Toki{
		Toki{
			POS: "",
			Meanings: []string{
				"I want to eat. You made me hungry.",
			},
		},
	},
	"pu": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"buying and interacting with the official Toki Pona book",
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
	},
	"tan": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"",
				"causal",
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
	},
	"suwi": []Toki{
		Toki{
			POS: "adjective",
			Meanings: []string{
				"cute",
				"sweet",
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
	return `(?i)^toki: (.+)$`
}

// Match determines if we are highfiving
func (t *Toki) Match(user, msg string) bool {
	re := regexp.MustCompile(t.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here
func (t *Toki) SetStore(s PluginStore) {}

func (t *Toki) fix(msg string) string {
	re := regexp.MustCompile(t.Re())
	return re.ReplaceAllString(msg, "$1")
}

// RespondText to hi events
func (t *Toki) RespondText(c *gomatrix.Client, ev *gomatrix.Event, user, post string) {
	w := t.fix(post)
	if word, ok := TokiLang[w]; ok {
		var defs []string
		for _, v := range word {
			defs = append(defs, v.Print(w))
		}
		SendMD(c, ev.RoomID, strings.Join(defs, "\n\n"))
	} else {
		SendText(c, ev.RoomID, "mi sona ala")
	}
}

// Name hi
func (t *Toki) Name() string {
	return "Toki"
}
