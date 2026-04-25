// Package scorecard tracks a per-host reliability score derived from
// successive scan outcomes.
//
// Each successful scan increases the score by a configurable gain value;
// each failed scan decreases it by a configurable decay value. Scores are
// clamped to the range [0, max].
//
// Typical usage:
//
//	sc := scorecard.New(
//		scorecard.WithMax(100),
//		scorecard.WithGain(5),
//		scorecard.WithDecay(10),
//	)
//	sc.RecordSuccess("192.168.1.1")
//	fmt.Println(sc.Get("192.168.1.1").Value) // 5
package scorecard
