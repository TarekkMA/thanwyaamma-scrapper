package score2018

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"gitlab.com/tarekkma/thanwyaamma-scrapper/score"
)

//Scrapper2018 struct
type scrapper2018 struct{}

//NewScrepper returns new scrapper for the year 2018
func NewScrepper() score.Scrapper {
	return new(scrapper2018)
}

func (s *scrapper2018) Get(seatingNumber int32) (*score.Result, *score.Error) {
	log.Info("==> Getting scores for ", seatingNumber)

	client := new(http.Client)

	homeReq, err := http.NewRequest("GET", "http://natega.thanwya.emis.gov.eg/", nil)

	if err != nil {
		log.Errorf("<== Error constructing main page request for %d : %v", seatingNumber, err)
		return nil, &score.Error{SeatingNumber: seatingNumber, Err: err}
	}

	homeRes, err := client.Do(homeReq)

	if err != nil {
		log.Errorf("<== Error while retrieveing main page for %d : %v", seatingNumber, err)
		return nil, &score.Error{SeatingNumber: seatingNumber, Err: err}
	}

	defer homeRes.Body.Close()
	if homeRes.StatusCode != 200 {
		log.Errorf("<== Status code for home page was't 200 was %d for %d : %v", homeRes.StatusCode, seatingNumber, err)
		return nil, &score.Error{SeatingNumber: seatingNumber, Err: fmt.Errorf("statuscode:%d", homeRes.StatusCode)}
	}

	doc, err := goquery.NewDocumentFromReader(homeRes.Body)
	if err != nil {
		log.Errorf("<== Can't init goquery for %d : %v", seatingNumber, err)
		return nil, &score.Error{SeatingNumber: seatingNumber, Err: err}
	}

	hiddenInput := make(map[string]string)

	hiddenNodes := doc.Find("input")

	hiddenNodes.Each(func(index int, sel *goquery.Selection) {
		name, _ := sel.Attr("name")
		value, _ := sel.Attr("value")
		hiddenInput[name] = value
	})

	form := url.Values{}

	for key, val := range hiddenInput {
		form.Add(key, val)
	}
	form.Set("TextBox1", fmt.Sprintf("%v", seatingNumber))

	scoreReq, err := http.NewRequest("POST", "http://natega.thanwya.emis.gov.eg/", strings.NewReader(form.Encode()))

	scoreReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		log.Errorf("<== Error constructing score page request for %d : %v", seatingNumber, err)
		return nil, &score.Error{SeatingNumber: seatingNumber, Err: err}
	}

	scoreRes, err := client.Do(scoreReq)

	if err != nil {
		log.Errorf("<== Error while retrieveing main page for %d : %v", seatingNumber, err)
		return nil, &score.Error{SeatingNumber: seatingNumber, Err: err}
	}

	defer scoreRes.Body.Close()
	if scoreRes.StatusCode != 200 {
		log.Errorf("<== Status code for score page was't 200 was %d for %d : %v", scoreRes.StatusCode, seatingNumber, err)
		return nil, &score.Error{SeatingNumber: seatingNumber, Err: fmt.Errorf("statuscode:%d", scoreRes.StatusCode)}
	}

	scoreDoc, err := goquery.NewDocumentFromReader(scoreRes.Body)
	if err != nil {
		log.Errorf("<== Can't init goquery for %d : %v", seatingNumber, err)
		return nil, &score.Error{SeatingNumber: seatingNumber, Err: err}
	}

	result := new(score.Result)

	result.SeatingNumber = seatingNumber

	result.Name = scoreDoc.Find("#std_name").Text()
	result.School = scoreDoc.Find("#school_name").Text()
	result.SeatingNumberStr = scoreDoc.Find("#seating_no").Text()
	result.ElModorya = scoreDoc.Find("#mud_name").Text()
	result.ElEdara = scoreDoc.Find("#edara_name").Text()

	result.ArabicScore = parseScore(scoreDoc.Find("#s1").Text())
	result.Lang1Score = parseScore(scoreDoc.Find("#s2").Text())
	result.Lang2Score = parseScore(scoreDoc.Find("#s3").Text())

	result.HistoryScore = parseScore(scoreDoc.Find("#s17").Text())
	result.GeographyScore = parseScore(scoreDoc.Find("#s8").Text())
	result.PhilosopheScore = parseScore(scoreDoc.Find("#s18").Text())
	result.PsychologyScore = parseScore(scoreDoc.Find("#s9").Text())

	result.BiologyScore = parseScore(scoreDoc.Find("#s5").Text())
	result.GeologyScore = parseScore(scoreDoc.Find("#s7").Text())
	result.ChemistryScore = parseScore(scoreDoc.Find("#s4").Text())
	result.PsychicsScore = parseScore(scoreDoc.Find("#s15").Text())
	result.PureMathematicsScore = parseScore(scoreDoc.Find("#s6").Text())
	result.AppliedMathematicsScore = parseScore(scoreDoc.Find("#s16").Text())

	result.TotalScore = parseScore(scoreDoc.Find("#total").Text())

	result.ReligionScore = parseScore(scoreDoc.Find("#s10").Text())
	result.CitizenshipScore = parseScore(scoreDoc.Find("#s14").Text())
	result.StatisticsScore = parseScore(scoreDoc.Find("#s19").Text())

	result.LagnaName = scoreDoc.Find("#CONTROL_NAME").Text()
	result.LagnaAddress = scoreDoc.Find("#CONTROL_ADDRESS").Text()
	result.LagnaPhoneNumber = scoreDoc.Find("#CONTROL_PHONE").Text()

	log.Info("<== Done for ", seatingNumber)

	return result, nil
}

func parseScore(text string) float64 {
	switch text {
	case "ــ":
		return -1
	case "ـــ":
		return -1
	case "غ":
		return -2
	}
	if num, err := strconv.ParseFloat(text, 64); err != nil {
		return -3
	} else {
		return num
	}
}
