package queststrg

type questionRow struct {
	Problem      string   `json:"problem"`
	CorrectIndex int      `json:"correct_index"`
	Answers      []string `json:"answers"`
}
