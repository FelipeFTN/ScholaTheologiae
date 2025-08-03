package controllers

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"scholatheologiae-api/constants"
	"scholatheologiae-api/models"
)

func (c *Controllers) SummaTheologiae(request models.SummaTheologiaeRequest) (any, error) {
	switch request.Type {
	case constants.SUMMA_LIST_PARTS:
		return c.summaListParts()
	case constants.SUMMA_LIST_QUESTIONS:
		return c.summaListQuestions(request.Part)
	case constants.SUMMA_GET_QUESTION:
		return c.summaGetQuestion(request.Part, request.Question)
	}

	return nil, errors.New("Unexpected error in Summa Controller")
}

func (c *Controllers) summaListParts() ([]string, error) {
	return c.data.GetSummaTheologiaeParts()
}

func (c *Controllers) summaListQuestions(part string) ([]string, error) {
	available_parts, err := c.summaListParts()
	if err != nil {
		return nil, err
	}

	if !slices.Contains(available_parts, part) {
		return nil, errors.New("Invalid Summa Part")
	}

	return c.data.GetSummaTheologiaeQuestions(part)
}

func (c *Controllers) summaListQuestionNums(part string) ([]string, error) {
	available_parts, err := c.summaListParts()
	if err != nil {
		return nil, err
	}

	if !slices.Contains(available_parts, part) {
		return nil, errors.New("Invalid Summa Part")
	}

	return c.data.GetSummaTheologiaeQuestionNums(part)
}

func (c *Controllers) summaGetQuestion(part string, question string) (any, error) {
	available_parts, err := c.summaListParts()
	if err != nil {
		return nil, err
	}

	if !slices.Contains(available_parts, part) {
		return nil, errors.New("Invalid Summa Part")
	}

	available_questions, err := c.summaListQuestionNums(part)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(available_questions, question) {
		return nil, errors.New("Invalid Summa Question")
	}

	file_path := fmt.Sprintf("./data/library/summa_theologiae/content/%s/questions/question_%s.md", part, question)
	file_bytes, err := os.ReadFile(file_path)
	if err != nil {
		return nil, err
	}

	return string(file_bytes), nil
}
