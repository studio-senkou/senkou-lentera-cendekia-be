package requests

type CreateTestimonyRequest struct {
	TestimonerName             string `json:"testimoner_name" validate:"required"`
	TestimonerCurrentPosition  string `json:"testimoner_current_position" validate:"required"`
	TestimonerPreviousPosition string `json:"testimoner_previous_position" validate:"required"`
	TestimonyText              string `json:"testimony_text" validate:"required"`
}

type UpdateTestimonyRequest struct {
	TestimonerName             string `json:"testimoner_name" validate:"required"`
	TestimonerCurrentPosition  string `json:"testimoner_current_position" validate:"required"`
	TestimonerPreviousPosition string `json:"testimoner_previous_position" validate:"required"`
	TestimonyText              string `json:"testimony_text" validate:"required"`
}
