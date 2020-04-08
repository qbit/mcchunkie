# Script to parse https://raw.githubusercontent.com/jan-Lope/Toki_Pona_lessons_English/gh-pages/toki-pona_english.txt
# run with: grep "::" toki-pona_english.txt | grep "[a-z]: " | awk -F: -f parse_toki-pona_ding.awk
{
	gsub(/ $/, "", $1);
	gsub(/^ /, "", $1);

	gsub(/ $/, "", $3);
	gsub(/^ /, "", $3);

	gsub(/ $/, "", $NF);
	gsub(/^ /, "", $NF);

	words[$1] = words[$1] "#" "$" $3 "$" $NF;
}

END {
	for (w in words) {
		print "\""w"\": []Toki{"
		split(words[w], pnd, "#")

		for (p in pnd) {
			split(pnd[p], a, "$")
			pos = a[2]
			meanings = a[3]

			if (pos != "") {
				print "\tToki{"

				print "\t\tPOS: " "\"" pos "\","
				print "\t\tMeanings: []string{"
					split(meanings, b, ",")
					for (x in b) {
						gsub(/ $/, "", b[x]);
						gsub(/^ /, "", b[x]);
						print "\t\t\t\""b[x]"\","
					}
				print "\t\t},"
				print "\t},"
			}

		}
		print "},"
	}
}
