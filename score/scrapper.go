package score

//Scrapper interfcae to abstract scrapping over the years
//this maybe not the best idea since new thanwaya will change everything
type Scrapper interface {
	Get(seatingNumber int32) (res *Result, err *Error)
}

//Error incase the scrapper faild
type Error struct {
	SeatingNumber int32
	Err           error
}

//Result of the student
type Result struct {
	SeatingNumber int32

	Name             string
	School           string
	SeatingNumberStr string
	ElModorya        string
	ElEdara          string

	ArabicScore             float64
	Lang1Score              float64
	Lang2Score              float64
	HistoryScore            float64
	GeographyScore          float64
	PhilosopheScore         float64
	PsychologyScore         float64
	BiologyScore            float64
	GeologyScore            float64
	ChemistryScore          float64
	PsychicsScore           float64
	PureMathematicsScore    float64
	AppliedMathematicsScore float64

	TotalScore float64

	ReligionScore    float64
	CitizenshipScore float64
	StatisticsScore  float64

	LagnaName        string
	LagnaAddress     string
	LagnaPhoneNumber string
}
