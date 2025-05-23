package plugins

import (
	"math/rand"
	"regexp"

	"github.com/matrix-org/gomatrix"
)

// ROA sends a random rule
type ROA struct {
}

// Descr describes this plugin
func (h *ROA) Descr() string {
	return "Rules of Acquisition."
}

// Re checks for roa:
func (h *ROA) Re() string {
	return `(?i)^roa:$`
}

// Match determines if we are asking for an roa
func (h *ROA) Match(user, msg string) bool {
	re := regexp.MustCompile(h.Re())
	return re.MatchString(msg)
}

// SetStore we don't need a store here
func (h *ROA) SetStore(_ PluginStore) {}

// Process
func (h *ROA) Process(from, post string) (string, func() string) {
	a := []string{
		`1	Once you have their money, you never give it back.`,
		`2	The best deal is the one that brings the most profit.`,
		`3	Never spend more for an acquisition than you have to.`,
		`4	A woman wearing clothes is like a man in the kitchen.`,
		`5	If you can't break a contract, bend it.`,
		`6	Never allow family to stand in the way of opportunity.`,
		`7	Keep your ears open.`,
		`9	Opportunity plus instinct equals profit.`,
		`10	Greed is eternal.`,
		`11	Latinum isn't the only thing that shines.`,
		`12	Anything worth selling is worth selling twice.`,
		`13	Anything worth doing is worth doing for money.`,
		`14	Anything stolen is pure profit.`,
		`15	Acting stupid is often smart.`,
		`16	A deal is a deal.`,
		`17	A contract is a contract is a contract ... but only between Ferengi.`,
		`18	A Ferengi without profit is no Ferengi at all.`,
		`19	Satisfaction is not guaranteed.`,
		`21	Never place friendship above profit.`,
		`22	A wise man can hear profit in the wind.`,
		`23	Nothing is more important than your health ... except for your money.`,
		`24	Latinum can't buy happiness, but you can sure have a blast renting it.`,
		`25	You can't make a deal if you're dead.`,
		`27	There's nothing more dangerous than an honest businessman.`,
		`31	Never make fun of a Ferengi's mother.`,
		`33	It never hurts to suck up to the boss.`,
		`34	War is good for business`,
		`35	Peace is good for business.`,
		`45	Expand or die.`,
		`47	Never trust a man wearing a better suit than your own.`,
		`48	The bigger the smile, the sharper the knife.`,
		`52	Never ask when you can take.`,
		`57	Good customers are as rare as latinum. Treasure them.`,
		`58	There is no substitute for success.`,
		`59	Free advice is seldom cheap.`,
		`60	Keep your lies consistent.`,
		`62	The riskier the road, the greater the profit.`,
		`65	Win or lose, there's always Hupyrian beetle snuff.`,
		`69	When she discusses money for "favors", charge her what she'll pay.`,
		`74	Knowledge equals profit.`,
		`75	Home is where the heart is, but the stars are made of latinum.`,
		`76	Every once in a while, declare peace. It confuses the hell out of your enemies.`,
		`79	Beware of the Vulcan greed for knowledge.`,
		`82	The flimsier the product, the higher the price.`,
		`85	Never let the competition know what you're thinking.`,
		`89	Ask not what your profits can do for you, but what you can do for your profits.`,
		`94	Females and finances don't mix.`,
		`95	Expand or die.`,
		`97	Enough...is never enough.`,
		`98	Every man has his price.`,
		`99	Trust is the biggest liability of all.`,
		`102	Nature decays, but latinum lasts forever.`,
		`103	Sleep can interfere with an opportunity.`,
		`104	Faith moves mountains...of inventory.`,
		`106	There is no honor in poverty.`,
		`109	Dignity and an empty sack is worth the sack.`,
		`111	Treat people in your debt like family ... exploit them.`,
		`112	Never have sex with the boss's sister.`,
		`113	Always have sex with the boss.`,
		`121	Everything is for sale, even friendship.`,
		`123	Even a blind man can recognize the glow of latinum.`,
		`125	You can't make a deal if you're dead.`,
		`139	Wives serve, brothers inherit.`,
		`141	Only fools pay retail.`,
		`144	There's nothing wrong with charity...as long as it winds up in your pocket.`,
		`162	Even in the worst of times, someone makes a profit.`,
		`168	Whisper your way to success.`,
		`177	Know your enemies...but do business with them always.`,
		`181	Not even dishonesty can tarnish the shine of profit.`,
		`189	Let others keep their reputation...you keep their latinum.`,
		` - 	A man is only worth the sum of his possessions.`,
		`190	Hear all, trust nothing.`,
		`192	Never cheat a Klingon...unless you can get away with it.`,
		`194	It's always good to know about new customers before they walk in your door.`,
		`202	The justification for profit is profit.`,
		`203	New customers are like razor-toothed gree-worms. They can be succulent, but sometimes they bite back.`,
		`208	Sometimes the only thing more dangerous than a question is an answer.`,
		`211	Employees are the rungs on the ladder of success. Don't hesitate to step on them.`,
		`214	Never begin a negotiation on an empty stomach.`,
		`217	You can't free a fish from water.`,
		`218	Always know what you're buying.`,
		`223	(incomplete, but presumably concerned the relationship between "keeping busy" and "being successful")`,
		`229	Latinum lasts longer than lust.`,
		`236	You can't buy fate.`,
		`239	Never be afraid to mislabel a product.`,
		`242	More is good...all is better.`,
		`255	A wife is a luxury...a smart accountant a necessity.`,
		`261	A wealthy man can afford anything except a conscience.`,
		`263	Never allow doubt to tarnish your lust for latinum.`,
		`266	When in doubt, lie.`,
		`284	Deep down, everyone's a Ferengi.`,
		`285	No good deed ever goes unpunished.`,
		`286 	When Morn leaves, it's all over.`,
		`299 	After you've exploited someone, it never hurts to thank them. That way, it's easier to exploit them next time.`,
		`-	Exploitation begins at home.`,
		`(The unwritten rule) 	When no appropriate rule applies, make one up.`,
		`-	When the messenger comes to appropriate your profits ... kill the messenger.`,
		`-	Time, like latinum, is a highly limited commodity.`,
		`-	Always inspect the merchandise before making a deal.`,
		`-	Money is money, but females are better.`,
		`-	Why ask, when you can take?`,
		`-	A good lie is easier to believe than the truth.`,
		`-	If that's what's written, then that's what's written.`,
	}

	return a[rand.Intn(len(a))], RespStub
}

// RespondText
func (h *ROA) RespondText(c *gomatrix.Client, ev *gomatrix.Event, _, _ string) error {
	resp, _ := h.Process(ev.Sender, "")
	return SendText(c, ev.RoomID, resp)
}

// Name ROA
func (h *ROA) Name() string {
	return "ROA"
}
