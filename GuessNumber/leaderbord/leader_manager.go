package leaderbord

import (
	"bufio"
	"encoding/json"
	"os"
	"sort"
)

const leaderboard = "./game.log"

func LoadLeaderBoard() (*Leaderboard, error) {
	file, err := os.Open(leaderboard)
	if err != nil {
		return &Leaderboard{}, nil
	}
	defer file.Close()

	lb := &Leaderboard{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var record Records
		line := scanner.Bytes()

		if err := json.Unmarshal(line, &record); err != nil {
			continue
		}
		if record.Result == "WIN" {
			lb.Records = append(lb.Records, record)
		}
	}
	sort.Slice(lb.Records, func(i, j int) bool {
		return lb.Records[i].Attempts < lb.Records[j].Attempts
	})

	if len(lb.Records) > 5 {
		lb.Records = lb.Records[:5]
	}

	return lb, scanner.Err()
}

func (lb *Leaderboard) GetTopWins() []Records {
	return lb.Records
}
