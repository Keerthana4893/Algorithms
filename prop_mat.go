package main

type Matched_req struct {
	Req
	Match_score float32
}

func NewMatched_req(re Req, score float32) Matched_req {
	return Matched_req{
		Req: re,
		Match_score:  score,
	}
}

type Prop_m interface {
	Match(Prop_list, []Req, ReqMargins) []Matched_req
}

type Prop_m1 struct{}

func NewProp_m_algo() Prop_m1_algo {
	return Prop_m1_algo{}
}

func (ar Prop_m1_algo) Match(pro Prop_list, requirements []Req, rMargins ReqMargins) []Matched_req {
	matchedReqs := []Matched_req{}
	scoring := make(chan bool)
	defer close(scoring)

	scores := ar.createReqScores(req)

	
	go ar.distanceMatching(pro.Latitude, pro.Longitude, scores, scoring)
	go ar.budgetMatching(pro.Price, requirements, scores, rMargins, scoring)
	go ar.bedroomsMatching(pro.Bedrooms, requirements, scores, rMargins, scoring)
	go ar.bathroomsMatching(pro.Bathrooms, requirements, scores, rMargins, scoring)
//The above is the go routines for budget,bathroom, bedroom and distance matching...
	
	for i := 0; i < 4; i++ {
		<-scoring
	}
	SortScores(scores)

	for i, _ := range scores {
		if scores[i].Total < 40.0 {
			continue
		}
		matchedReqs = append(matchedReqs, NewMatched_req(req[scores[i].Index].Req, scores[i].Total))
	}
	return matchedReqs
}

func (ar Prop_m1_algo) createReqScores(re []Req) []Score {
	scores := make([]Score, len(re))
	for i, _ := range re {
		scores[i] = NewScore(i, re[i].Dist)
	}
	return scores
}

func (ar Prop_m1_algo) dist_matching(lat, lon float32, scores []Score, scoring chan bool) {
	base_dist := float32(2)
	max_dist := float32(10)

	for i, _ := range scores {
		scores[i].Dist_sc= GetDist_sc(scores[i].Dist, base_dist, max_dist)
	}
	scoring <- true
}

func (ar Prop_m1_algo) budgetMatching(price float32, re []Req, scores []Score, rMargins ReqMargins, scoring chan bool) {
	for i, _ := range scores {
		scores[i].BudgetScore = GetBudgetScore(re[i].MinBudget, re[i].MaxBudget, price, rMargins.MinPrice, rMargins.MaxPrice)
	}
	scoring <- true
}

func (ar Prop_m1_algo) bedroomsMatching(bedrooms int16, re []Req, scores []Score, rMargins ReqMargins, scoring chan bool) {
	for i, _ := range scores {
		scores[i].BedroomScore = GetBedroomScore(re[i].MinBedrooms, re[i].MaxBedrooms, bedrooms, rMargins.MinBeds, rMargins.MaxBeds)
	}
	scoring <- true
}

func (ar Prop_m1_algo) bathroomsMatching(bathrooms int16, re []Req, scores []Score, rMargins ReqMargins, scoring chan bool) {
	for i, _ := range scores {
scores[i].BathroomScore = GetBathroomScore(re[i].MinBathrooms, re[i].MaxBathrooms, bathrooms, rMargins.MinBaths, rMargins.MaxBaths)
	}
	scoring <- true
}
