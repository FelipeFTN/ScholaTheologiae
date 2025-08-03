package data

import (
	"sort"
	"strconv"
	"strings"
)

func (d *Data) GetSummaTheologiae(part, question, article string) (string, error) {
	// Prepare the statement
	query := "SELECT article_text FROM summa_theologiae WHERE part = ? AND question_num = ? AND article = ?"
	stmt, err := d.SQLite.databases[SUMMA_THEOLOGIAE].db.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	// Execute the statement
	row := stmt.QueryRow(part, question, article)

	// Scan the result into a variable
	var text string
	err = row.Scan(&text)
	if err != nil {
		return "", err
	}

	return text, nil
}

func (d *Data) GetSummaTheologiaeParts() ([]string, error) {
	// Prepare the statement
	query := "SELECT DISTINCT part FROM summa_theologiae"
	stmt, err := d.SQLite.databases[SUMMA_THEOLOGIAE].db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the statement
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the results into a slice
	var parts []string
	for rows.Next() {
		var part string
		err = rows.Scan(&part)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}

	return parts, nil
}

func (d *Data) GetSummaTheologiaeQuestions(part string) ([]string, error) {
	// Prepare the statement
	query := "SELECT question_title, question_num FROM summa_theologiae WHERE part = ?"
	stmt, err := d.SQLite.databases[SUMMA_THEOLOGIAE].db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the statement
	rows, err := stmt.Query(part)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the results into a slice
	questions := make(map[string]string, 0)
	for rows.Next() {
		var question, questionNum string
		err = rows.Scan(&question, &questionNum)
		if err != nil {
			return nil, err
		}
		questions[questionNum] = question
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Convert the map to a slice
	// "questionNum: QuestionTitle"
	questionsSlice := make([]string, 0, len(questions))
	for questionNum, questionTitle := range questions {
		questionsSlice = append(questionsSlice, questionNum+": "+questionTitle)
	}

	// Sort the slice ascendingly
	sort.Slice(questionsSlice, func(i, j int) bool {
		// Split the string by ": " and compare the first part (question number)
		a, _ := strconv.Atoi(strings.Split(questionsSlice[i], ":")[0])
		b, _ := strconv.Atoi(strings.Split(questionsSlice[j], ":")[0])
		return a < b
	})

	return questionsSlice, nil
}

func (d *Data) GetSummaTheologiaeQuestionNums(part string) ([]string, error) {
	// Prepare the statement
	query := "SELECT DISTINCT question_num FROM summa_theologiae WHERE part = ? ORDER BY question_num ASC"
	stmt, err := d.SQLite.databases[SUMMA_THEOLOGIAE].db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the statement
	rows, err := stmt.Query(part)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the results into a slice
	var questions []string
	for rows.Next() {
		var question string
		err = rows.Scan(&question)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil
}

func (d *Data) GetSummaTheologiaeArticles(part, question string) ([]string, error) {
	// Prepare the statement
	query := "SELECT DISTINCT article FROM summa_theologiae WHERE part = ? AND question_num = ?"
	stmt, err := d.SQLite.databases[SUMMA_THEOLOGIAE].db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the statement
	rows, err := stmt.Query(part, question)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the results into a slice
	var articles []string
	for rows.Next() {
		var article string
		err = rows.Scan(&article)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (d *Data) GetSummaTheologiaeArticle(part, question, article string) (string, error) {
	// Prepare the statement
	query := "SELECT article_text FROM summa_theologiae WHERE part = ? AND question_num = ? AND article = ?"
	stmt, err := d.SQLite.databases[SUMMA_THEOLOGIAE].db.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	// Execute the statement
	row := stmt.QueryRow(part, question, article)

	// Scan the result into a variable
	var text string
	err = row.Scan(&text)
	if err != nil {
		return "", err
	}

	return text, nil
}
