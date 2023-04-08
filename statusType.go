package main

// StatusType is the type of a status message
type StatusType int

const (
	// StatusTypeNeutral is a neutral status message (ex: a loading message)
	StatusTypeNeutral = iota
	// StatusTypeNegative is a negative status message (ex: an error)
	StatusTypeNegative
	// StatusTypePositive is a positive status message (ex: an action success)
	StatusTypePositive
)

func (s StatusType) String() string {
	switch s {
	case StatusTypeNeutral:
		return "neutral"
	case StatusTypeNegative:
		return "negative"
	case StatusTypePositive:
		return "positive"
	default:
		return "unknown"
	}
}
