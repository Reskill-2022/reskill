package model

type User struct {
	// Basic
	Email       string `json:"email" firestore:"email"`
	Name        string `json:"name" firestore:"name"`
	LinkedInURL string `json:"linkedin_url" firestore:"linkedin_url"`
	Location    string `json:"location" firestore:"location"`
	Timezone    string `json:"timezone" firestore:"timezone"`
	Phone       string `json:"phone" firestore:"phone"`
	FirstName   string `json:"first_name" firestore:"first_name"`
	LastName    string `json:"last_name" firestore:"last_name"`
	Photo       string `json:"photo" firestore:"photo"`

	// Extras
	Representation         string `json:"representation" firestore:"representation"`
	Gender                 string `json:"gender" firestore:"gender"`
	AgeGroup               string `json:"age_group" firestore:"age_group"`
	EmploymentStatus       string `json:"employment_status" firestore:"employment_status"`
	HighestSchool          string `json:"highest_school" firestore:"highest_school"`
	OptionalMajor          string `json:"optional_major" firestore:"optional_major"`
	CanWorkInUSA           string `json:"can_work_in_usa" firestore:"can_work_in_usa"`
	LearningTrack          string `json:"learning_track" firestore:"learning_track"`
	TechExperience         string `json:"tech_experience" firestore:"tech_experience"`
	HoursPerWeek           string `json:"hours_per_week" firestore:"hours_per_week"`
	Referral               string `json:"referral" firestore:"referral"`
	ReferralOther          string `json:"referral_other" firestore:"referral_other"`
	City                   string `json:"city" firestore:"city"`
	State                  string `json:"state" firestore:"state"`
	ProfessionalExperience string `json:"professional_experience" firestore:"professional_experience"`
	Industries             string `json:"industries" firestore:"industries"`
	WillChangeJob          string `json:"will_change_job" firestore:"will_change_job"`
	WillChangeJobRole      string `json:"will_change_job_role" firestore:"will_change_job_role"`
	OpenToMeet             string `json:"open_to_meet" firestore:"open_to_meet"`
	RacialDemographic      string `json:"racial_demographic" firestore:"racial_demographic"`
	PriorKnowledge         string `json:"prior_knowledge" firestore:"prior_knowledge"`

	// Meta
	Enrolled     bool   `json:"enrolled" firestore:"enrolled"`
	CreatedAt    string `json:"created_at" firestore:"created_at"`
	GitAccount   string `json:"gitaccount" firestore:"gitaccount"`
	FigmaAccount string `json:"figmaaccount" firestore:"figmaaccount"`
	GitYes       string `json:"git_yes" firestore:"git_yes"`
	FigmaYes     string `json:"figma_yes" firestore:"figma_yes"`
}
