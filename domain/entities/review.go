package entities

type Review struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	Stars       int64  `json:"stars"`
	User        User   `json:"reviewer"`
	ReviewDate  string `json:"review_date"`
}

type AddReviewRequest struct {
	Description string `json:"description"`
	Stars       int64  `json:"stars"`
	User        User   `json:"reviewer"`
}

type AddReviewResponse struct {
	ID int64 `json:"id"`
}
