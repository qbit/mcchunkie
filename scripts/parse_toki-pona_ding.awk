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
		split(words[w], pnd, "#")
		print "\""w"\": []Toki{"
				print "\tToki{"
					split(pnd[2], a, "$")
					print "\t\tPOS: " "\""a[2]"\","
					split(a[3], b, ",")
					print "\t\tMeanings: []string{"
						for (x in b) {
							gsub(/ $/, "", b[x]);
							gsub(/^ /, "", b[x]);
							print "\t\t\t\""b[x]"\","
						}
					print "\t\t},"
				print "\t},"
		print "},"
	}
}
