package controller

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"sort"
	"strconv"

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

func (c *Controller) summaGetQuestion(part string, question string) (any, error) {
	available_parts, err := c.summaListParts()
	if err != nil {
		return nil, err
	}

	if !slices.Contains(available_parts, part) {
		return nil, errors.New("Invalid Summa Part")
	}

	part_directory := fmt.Sprintf("../articles/suma_teologica/%s/questions/%s", part, question)
	dir_questions, err := os.ReadDir(part_directory)
	if err != nil {
		return nil, err
	}

	var articles []string
	for _, q := range dir_questions {
		if q.IsDir() {
			continue
		}

		match_question_num := regexp.MustCompile(`\d+`)

		articles = append(articles, match_question_num.FindString(q.Name()))
	}

	var nums []int = make([]int, len(articles))
	for i, s := range articles {
		num, err := strconv.Atoi(s) // Convert string to int
		if err != nil {
			return nil, err
		}

		nums[i] = num
	}

	sort.Ints(nums)
	articles = make([]string, 0)
	for _, num := range nums {
		articles = append(articles, strconv.Itoa(num))
	}

	return articles, nil
}

func (c *Controller) summaGetArticle(part string, question string, article string) (any, error) {
	available_parts, err := c.summaListParts()
	if err != nil {
		return nil, err
	}

	if !slices.Contains(available_parts, part) {
		return nil, errors.New("Invalid Summa Part")
	}

	part_directory := fmt.Sprintf("../articles/suma_teologica/%s/questions/%s/articles", part, question)
	dir_questions, err := os.ReadDir(part_directory)
	if err != nil {
		return nil, err
	}

	var articles []string
	for _, q := range dir_questions {
		if q.IsDir() {
			continue
		}

		match_question_num := regexp.MustCompile(`\d+`)

		articles = append(articles, match_question_num.FindString(q.Name()))
	}

	if !slices.Contains(articles, article) {
		return nil, errors.New("Invalid Summa Article")
	}

	return c.db.GetSummaTheologiaeArticle(part, question, article)
}
