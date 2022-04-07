package requests

type (
	CreateUserRequest struct {
		Email string `json:"email"`
	}

	UpdateUserRequest struct {
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
)
