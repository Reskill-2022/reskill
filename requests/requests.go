package requests

type (
	CreateUserRequest struct {
		AuthCode    string `json:"code"`
		RedirectURI string `json:"redirect_uri"`
	}

	UpdateUserRequest struct {
		LinkedInURL            string `json:"linkedin_url"`
		Timezone               string `json:"timezone"`
		Phone                  string `json:"phone"`
		Representation         string `json:"representation"`
		Gender                 string `json:"gender"`
		AgeGroup               string `json:"age_group"`
		EmploymentStatus       string `json:"employment_status"`
		HighestSchool          string `json:"highest_school"`
		OptionalMajor          string `json:"field_of_study"`
		CanWorkInUSA           string `json:"can_work_in_usa"`
		LearningTrack          string `json:"learning_track"`
		TechExperience         string `json:"tech_experience"`
		HoursPerWeek           string `json:"hours_per_week"`
		Referral               string `json:"referral"`
		ReferralOther          string `json:"referral_other"`
		Photo                  string `json:"photo"`
		City                   string `json:"city"`
		State                  string `json:"state"`
		ProfessionalExperience string `json:"professional_experience"`
		Industries             string `json:"industries"`
		WillChangeJob          string `json:"will_change_job"`
		WillChangeJobRole      string `json:"will_change_job_role"`
		OpenToMeet             string `json:"open_to_meet"`
		RacialDemographic      string `json:"racial_demographic"`
		PriorKnowledge         string `json:"prior_knowledge"`
	}
)
