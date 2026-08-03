package quiz

// Question is one multiple-choice item. Answer is the 0-based index into Choices.
// Explanation, if set, is shown when the question is missed.
type Question struct {
	Question    string   `json:"question"`
	Choices     []string `json:"choices"`
	Answer      int      `json:"answer"`
	Explanation string   `json:"explanation,omitempty"`
	ModuleID    int      `json:"-"`
}

// Module is a named, numbered set of questions (e.g. "Module 1").
type Module struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Questions []Question `json:"questions"`
}