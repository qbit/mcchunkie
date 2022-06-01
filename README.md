
# mcchunkie

[![mcchunkie's face](mcchunkie.png)](https://builds.sr.ht/~qbit/mcchunkie?)

[![builds.sr.ht status](https://builds.sr.ht/~qbit/mcchunkie.svg)](https://builds.sr.ht/~qbit/mcchunkie?)

A [Matrix](https://matrix.org) chat bot.

|Plugin Name|Match|Description|
|----|---|---|
|Beat|`(?i)^\.beat$\|^what time is it[\?!]+$\|^beat( )?time:?\??$`|Print the current [beat time](https://en.wikipedia.org/wiki/Swatch_Internet_Time).|
|Beer|`(?i)^beer: `|Queries [OpenDataSoft](https://public-us.opendatasoft.com/explore/dataset/open-beer-database/table/)'s beer database for a given beer.|
|BotSnack|`(?i)botsnack`|Consumes a botsnack. This pleases mcchunkie and brings balance to the universe.|
|Covid|`(?i)^covid: (.+)$`|Queries [thebigboard.cc](http://www.thebigboard.cc)'s api for information on COVID-19.|
|DMR|`(?i)^dmr (user\|repeater) (surname\|id\|callsign\|city\|county) (.+)$`|Queries radioid.net|
|Feder|`(?i)^(?:feder: \|tayshame: )(.*)$`|check the Matrix federation status of a given URL.|
|Groan|`(?i)^@groan$`|Ugh.|
|Ham|`(?i)^ham: (\w+)$`|queries HamDB.org for a given callsign.|
|HighFive|`o/\|\\o`|Everyone loves highfives.|
|Hi|`(?i)^hi\|hi$`|Friendly bots say hi.|
|LoveYou|`(?i)i love you`|Spreading love where ever we can by responding when someone shows us love.|
|OpenBSDMan|`(?i)^man: ([1-9][p]?)?\s?(.+)$`|Produces a link to man.openbsd.org.|
|PGP|`(?i)^pgp: (.+@.+\..+\|[a-f0-9]+)$`|Queries keys.openpgp.org|
|Palette|`(?i)^#[a-f0-9]{6}$`|Creates an solid 56x56 image of the color specified.|
|RFC|`(?i)^rfc\s?([0-9]+)$`|Produces a link to tools.ietf.org.|
|Salute|`o7`|Everyone loves salutes.|
|Snap|`(?i)^snap:$`|checks the current build date of OpenBSD snapshots.|
|Source|`(?i)where is your (source\|code)`|Tell people where they can find more information about myself.|
|Thanks|`(?i)^thank you\|thank you$\|^thanks\|thanks$\|^ty\|ty$`|Bots should be respectful. Respond to thanks.|
|Homestead|`(?i)^home:\|^homestead:\s?(\w+)?$`|Display weather information for the Homestead|
|Toki|`(?i)^(toki[\?]?):? (.+)$`|Toki Pona dictionary|
|Version|`(?i)version$`|Show a bit of information about what we are.|
|Wb|`(?i)^welcome back\|welcome back$\|^wb\|wb$`|Respond to welcome back messages.|
|Weather|`(?i)^weather: (\d+)$`|Produce weather information for a given ZIP code. Data comes from openweathermap.org.|
