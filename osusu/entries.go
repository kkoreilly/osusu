package osusu

// Entries is a slice of multiple entries
type Entries []Entry

// // MissingData returns whether the given user is missing data in any of the given entries
// func (e Entries) MissingData(user User) bool {
// 	for _, entry := range e {
// 		if entry.MissingData(user) {
// 			return true
// 		}
// 	}
// 	return false
// }
