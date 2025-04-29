import pdfplumber
import re
import os
import json
from pathlib import Path

def create_output_dirs(base_path):
    """Create directory for output Markdown files."""
    os.makedirs(base_path, exist_ok=True)

def write_markdown_file(filepath, content, metadata):
    """Write content to a Markdown file with metadata."""
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write("---\n")
        for key, value in metadata.items():
            f.write(f"{key}: {value}\n")
        f.write(f"---\n\n{content}")

def get_questions(line):
    """Detect if a line starts a new Question and return its number."""
    question_pattern = re.compile(r'QUESTION (\d+)$', re.IGNORECASE)
    match = question_pattern.match(line.strip())
    if match:
        return match.group(1)
    return None

def get_articles(line):
    """Detect if a line starts a new Article and return its number."""
    article_pattern = re.compile(r'([A-Z]+) ARTICLE \[I, Q\. \d+, Art\. (\d+)\]$', re.IGNORECASE)
    match = article_pattern.match(line.strip())
    if match:
        return match.group(2)
    return None

def clean_line(line):
    """Remove unwanted text from a line."""
    return re.sub(r'www\.freecatholicebooks\.com', '', line).strip()

def save_question(output_dir, question_num, content, metadata, index_data):
    """Save a Question's content and metadata to a Markdown file."""
    if not question_num or not content:
        return
    # Add Jekyll-specific metadata
    metadata['layout'] = 'question'
    metadata['permalink'] = f'/question/{question_num}'
    filename = f"question_{question_num}.md"
    filepath = os.path.join(output_dir, 'questions', filename)
    content_text = '\n'.join(content)
    write_markdown_file(filepath, content_text, metadata)

def save_article(output_dir, question_num, article_num, content, metadata):
    """Save an Article's content and metadata to a Markdown file."""
    if not question_num or not article_num or not content:
        return
    # Create directory for the Question's Articles
    question_dir = os.path.join(output_dir, 'articles', f'question_{question_num}')
    create_output_dirs(question_dir)
    # Add Jekyll-specific metadata
    metadata['question'] = question_num
    metadata['article'] = article_num
    metadata['layout'] = 'article'
    metadata['permalink'] = f'/question/{question_num}/article/{article_num}'
    filename = f"article_{article_num}.md"
    filepath = os.path.join(question_dir, filename)
    content_text = '\n'.join(content)
    write_markdown_file(filepath, content_text, metadata)

def process_page(page, current_question, current_article, current_content, question_content, metadata, output_dir, index_data, awaiting_title, awaiting_article_title):
    """Process a single PDF page, detecting Questions and Articles."""
    text = page.extract_text()
    if not text:
        return current_question, current_article, current_content, question_content, metadata, awaiting_title, awaiting_article_title

    lines = text.split('\n')
    for line in lines:
        line = clean_line(line)
        if not line:
            continue

        # If we're awaiting a Question title
        if awaiting_title and current_question:
            metadata['title'] = line
            question_content.append(f"# Question {current_question}: {line}")
            awaiting_title = False
            continue

        # If we're awaiting an Article title
        if awaiting_article_title and current_article:
            metadata['article_title'] = line
            current_content.append(f"# Article {current_article}: {line}")
            # Add to the Question content with a subtitle (include the marker line)
            question_content.append(f"\n## {last_article_marker}")
            question_content.append(line)
            awaiting_article_title = False
            continue

        # Check for a new Question
        question_num = get_questions(line)
        if question_num:
            # Save the previous Question and Article if they exist
            save_article(output_dir, current_question, current_article, current_content, metadata.copy())
            if current_question:
                index_data[f"question_{current_question}"] = [f"article_{i}" for i in range(1, int(current_article or 0) + 1)]
            save_question(output_dir, current_question, question_content, metadata, index_data)
            
            # Start a new Question
            current_question = question_num
            current_article = None
            current_content = []
            question_content = []
            metadata = {
                'question': question_num,
                'page_start': page.page_number
            }
            awaiting_title = True
            awaiting_article_title = False
            print(f"CurrentQuestionMatch: {question_num}; PageNum: {page.page_number}")
            continue

        # Check for a new Article
        article_num = get_articles(line)
        if article_num:
            # Save the previous Article if it exists
            save_article(output_dir, current_question, current_article, current_content, metadata.copy())
            
            # Start a new Article
            current_article = article_num
            current_content = []
            last_article_marker = line  # Store the marker line for use in Question content
            awaiting_article_title = True
            continue

        # Collect content for the current Article and Question
        if current_question:
            if current_article and not awaiting_article_title:
                current_content.append(line)
            question_content.append(line)

    return current_question, current_article, current_content, question_content, metadata, awaiting_title, awaiting_article_title

def extract_questions(pdf_path, output_dir):
    """Extract Questions and Articles from a PDF and save as Markdown files."""
    base_dir = os.path.join(output_dir, 'summa_theologica')
    create_output_dirs(os.path.join(base_dir, 'questions'))
    create_output_dirs(os.path.join(base_dir, 'articles'))
    
    current_question = None
    current_article = None
    current_content = []
    question_content = []
    metadata = {}
    index_data = {}
    awaiting_title = False
    awaiting_article_title = False
    
    try:
        with pdfplumber.open(pdf_path) as pdf:
            for page in pdf.pages:
                current_question, current_article, current_content, question_content, metadata, awaiting_title, awaiting_article_title = process_page(
                    page, current_question, current_article, current_content, question_content, metadata, base_dir, index_data, awaiting_title, awaiting_article_title
                )
        
        # Save the last Question and Article
        save_article(base_dir, current_question, current_article, current_content, metadata.copy())
        if current_question:
            index_data[f"question_{current_question}"] = [f"article_{i}" for i in range(1, int(current_article or 0) + 1)]
        save_question(base_dir, current_question, question_content, metadata, index_data)
        
        # Save index.json
        with open(os.path.join(base_dir, 'index.json'), 'w') as f:
            json.dump(index_data, f, indent=2)
    
    except FileNotFoundError:
        print(f"Error: PDF file '{pdf_path}' not found.")
    except Exception as e:
        print(f"Error processing PDF: {str(e)}")

if __name__ == "__main__":
    pdf_path = "summa_theologica.pdf"  # Replace with your PDF file path
    output_dir = "web_app"
    extract_questions(pdf_path, output_dir)
    print(f"Sections extracted to {output_dir}")
