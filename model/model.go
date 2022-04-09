package model

type User struct {
	// Basic
	Email       string `json:"email" firestore:"email"`
	Name        string `json:"name" firestore:"name"`
	LinkedInURL string `json:"linkedin_url" firestore:"linkedin_url"`
	Location    string `json:"location" firestore:"location"`
	Timezone    string `json:"timezone" firestore:"timezone"`
	Phone       string `json:"phone" firestore:"phone"`

	// Extras
	Representation   string `json:"representation" firestore:"representation"`
	Gender           string `json:"gender" firestore:"gender"`
	AgeGroup         string `json:"age_group" firestore:"age_group"`
	EmploymentStatus string `json:"employment_status" firestore:"employment_status"`
	HighestSchool    string `json:"highest_school" firestore:"highest_school"`
	CanWorkInUSA     string `json:"can_work_in_usa" firestore:"can_work_in_usa"`
	LearningTrack    string `json:"learning_track" firestore:"learning_track"`
	TechExperience   string `json:"tech_experience" firestore:"tech_experience"`
	HoursPerWeek     string `json:"hours_per_week" firestore:"hours_per_week"`
	Referral         string `json:"referral" firestore:"referral"`
}
