package internal

/**
 * Compare new and old maps with adds.
 * Since we do not care abound values (url links) we only compare keys
 */
func Compare(old, new map[string]string) (newEntries []string) {
	for keyNew, valNew := range new {
		found := false

		for keyOld, _ := range old {
			if keyOld == keyNew {
				found = true
				break
			}
		}

		if !found {
			newEntries = append(newEntries, "www.list.am"+valNew)
		}
	}

	return newEntries
}
