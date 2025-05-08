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

func SummaTheologiae(request models.SummaTheologiaeRequest) (any, error) {
	switch request.Type {
	case constants.SUMMA_LIST_PARTS:
		return summaListParts()
	case constants.SUMMA_LIST_QUESTIONS:
		return summaListQuestions(request.Part)
	case constants.SUMMA_GET_QUESTION:
		break
	case constants.SUMMA_GET_ARTICLE:
		break

	}

	return nil, errors.New("Unexpected error in Summa Controller")
}

func summaListParts() ([]string, error) {
	dir_parts, err := os.ReadDir("../articles/suma_teologica/")
	if err != nil {
		return nil, err
	}

	var parts []string
	for _, d := range dir_parts {
		if !d.IsDir() {
			continue
		}

		parts = append(parts, d.Name())
	}

	return parts, nil
}

func summaListQuestions(part string) ([]string, error) {
	available_parts, err := summaListParts()
	if err != nil {
		return nil, err
	}

	if !slices.Contains(available_parts, part) {
		return nil, errors.New("Invalid Summa Part")
	}

	part_directory := fmt.Sprintf("../articles/suma_teologica/%s/questions", part)
	dir_questions, err := os.ReadDir(part_directory)
	if err != nil {
		return nil, err
	}

	var questions []string
	for _, q := range dir_questions {
		if q.IsDir() {
			continue
		}

		match_question_num := regexp.MustCompile(`\d+`)

		questions = append(questions, match_question_num.FindString(q.Name()))
	}

	var nums []int = make([]int, len(questions))
	for i, s := range questions {
		num, err := strconv.Atoi(s) // Convert string to int
		if err != nil {
			return nil, err
		}

		nums[i] = num
	}

	sort.Ints(nums)
	questions = make([]string, 0)
	for _, num := range nums {
		questions = append(questions, strconv.Itoa(num))
	}

	return questions, nil
}
