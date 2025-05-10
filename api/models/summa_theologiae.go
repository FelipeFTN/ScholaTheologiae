package models

import "scholatheologiae-api/constants"

type SummaTheologiaeRequest struct {
	Part     string
	Question string
	Article  string
	Type     string
}

func (s *SummaTheologiaeRequest) Validate() {
	if s.Part == "" {
		s.Type = constants.SUMMA_LIST_PARTS
		return
	}

	if s.Question == "" {
		s.Type = constants.SUMMA_LIST_QUESTIONS
		return
	}

	if s.Article == "" {
		s.Type = constants.SUMMA_GET_QUESTION
		return
	}

	s.Type = constants.SUMMA_GET_ARTICLE
	return
}
