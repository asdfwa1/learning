package leaderbord

import (
	"sort"
	"time"
)

type Records struct {
	SecretNumber int    `json:"secret_number"`
	Attempts     int    `json:"attempts_used"`
	Result       string `json:"result"`
	Timestamp    string `json:"time"`
}

type Leaderboard struct {
	Records []Records `json:"records"`
}

func (l *Leaderboard) AddRecord(secNum, attempts int, result string) {
	l.Records = append(l.Records, Records{
		SecretNumber: secNum,
		Attempts:     attempts,
		Result:       result,
		Timestamp:    time.Now().Format(time.RFC3339),
	})
	sort.Slice(l.Records, func(i, j int) bool {
		if l.Records[i].Result == "WIN" && l.Records[j].Result == "WIN" {
			return l.Records[i].Attempts < l.Records[j].Attempts
		}
		return l.Records[i].Result == "WIN"
	})

	if len(l.Records) > 5 {
		l.Records = l.Records[:5]
	}
}
