package controller

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"scholatheologiae-api/constants"
	"scholatheologiae-api/models"
)

func (c *Controller) SummaTheologiae(request models.SummaTheologiaeRequest) (any, error) {
	switch request.Type {
	case constants.SUMMA_LIST_PARTS:
		return c.summaListParts()
	case constants.SUMMA_LIST_QUESTIONS:
		return c.summaListQuestions(request.Part)
	case constants.SUMMA_GET_QUESTION:
		return c.summaGetQuestion(request.Part, request.Question)
	case constants.SUMMA_GET_ARTICLE:
		return c.summaGetArticle(request.Part, request.Question, request.Article)
	}

	return nil, errors.New("Unexpected error in Summa Controller")
}

func (c *Controller) summaListParts() ([]string, error) {
	return c.db.GetSummaTheologiaeParts()
}

func (c *Controller) summaListQuestions(part string) ([]string, error) {
	available_parts, err := c.summaListParts()
	if err != nil {
		return nil, err
	}

	if !slices.Contains(available_parts, part) {
		return nil, errors.New("Invalid Summa Part")
	}

	return c.db.GetSummaTheologiaeQuestions(part)
}

func (c *Controller) summaListQuestionNums(part string) ([]string, error) {
	available_parts, err := c.summaListParts()
	if err != nil {
		return nil, err
	}

	if !slices.Contains(available_parts, part) {
		return nil, errors.New("Invalid Summa Part")
	}

	return c.db.GetSummaTheologiaeQuestionNums(part)
}

func (c *Controller) summaGetQuestion(part string, question string) (any, error) {
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

	file_path := fmt.Sprintf("./data/summa_theologiae/%s/questions/question_%s.md", part, question)
	file_bytes, err := os.ReadFile(file_path)
	if err != nil {
		return nil, err
	}

	return string(file_bytes), nil
}

func (c *Controller) summaGetArticle(part string, question string, article string) (any, error) {
	// available_parts, err := c.summaListParts()
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if !slices.Contains(available_parts, part) {
	// 	return nil, errors.New("Invalid Summa Part")
	// }
	//
	// available_questions, err := c.summaListQuestions(part)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if !slices.Contains(available_questions, question) {
	// 	return nil, errors.New("Invalid Summa Question")
	// }
	//
	// available_articles, err := c.db.GetSummaTheologiaeArticles(part, question)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if !slices.Contains(available_articles, article) {
	// 	return nil, errors.New("Invalid Summa Article")
	// }

	return nil, errors.New("Not yet implemented")
}
