package model

type User struct {
	// Basic
	Email       string `json:"email"`
	Name        string `json:"name"`
	LinkedInURL string `json:"linkedin_url"`
	Location    string `json:"location"`
	Timezone    string `json:"timezone"`
	Phone       string `json:"phone"`

	// Extras
	Representation   string `json:"representation"`
	Gender           string `json:"gender"`
	AgeGroup         string `json:"age_group"`
	EmploymentStatus string `json:"employment_status"`
	HighestSchool    string `json:"highest_school"`
	CanWorkInUSA     string `json:"can_work_in_usa"`
	LearningTrack    string `json:"learning_track"`
	TechExperience   string `json:"tech_experience"`
	HoursPerWeek     string `json:"hours_per_week"`
	Referral         string `json:"referral"`
}
