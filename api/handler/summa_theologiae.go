package handler

import "scholatheologiae-api/data"

func (h *SQLiteHandler) GetSummaTheologiae(part, question, article string) (string, error) {
	// Prepare the statement
	stmt, err := h.databases[data.SUMMA_THEOLOGIAE].db.Prepare("SELECT article_title FROM summa_theologiae WHERE part = ? AND question_num = ? AND article = ?")
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

func (h *SQLiteHandler) GetSummaTheologiaeParts() ([]string, error) {
	// Prepare the statement
	stmt, err := h.databases[data.SUMMA_THEOLOGIAE].db.Prepare("SELECT DISTINCT part FROM summa_theologiae")
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

func (h *SQLiteHandler) GetSummaTheologiaeQuestions(part string) ([]string, error) {
	// Prepare the statement
	stmt, err := h.databases[data.SUMMA_THEOLOGIAE].db.Prepare("SELECT DISTINCT question FROM summa_theologiae WHERE part = ?")
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

func (h *SQLiteHandler) GetSummaTheologiaeArticles(part, question string) ([]string, error) {
	// Prepare the statement
	stmt, err := h.databases[data.SUMMA_THEOLOGIAE].db.Prepare("SELECT DISTINCT article FROM summa_theologiae WHERE part = ? AND question = ?")
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

func (s *SQLiteHandler) GetSummaTheologiaeArticle(part, question, article string) (string, error) {
	// Prepare the statement
	stmt, err := s.databases[data.SUMMA_THEOLOGIAE].db.Prepare("SELECT article_text FROM summa_theologiae WHERE part = ? AND question = ? AND article = ?")
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
