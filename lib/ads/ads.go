package ads

import "sort"

// Representation of the advertisement on list.am
type Ad struct {
	Link  string
	Price string
	Dscr  string
	At    string
}

// Remove duplicates in ads slice. Use At field as a unique id
func RemoveDuplicatesInPlace(ads []Ad) []Ad {
	// if there are 0 or 1 items we return the slice itself
	if len(ads) < 2 {
		return ads
	}

	sort.SliceStable(ads, func(i, j int) bool { return ads[i].At < ads[j].At })

	uniqPointer := 0
	for i := 1; i < len(ads); i++ {
		if ads[uniqPointer] != ads[i] {
			uniqPointer++
			ads[uniqPointer] = ads[i]
		}
	}

	return ads[:uniqPointer+1]
}
