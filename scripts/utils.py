import pdfplumber
import sqlite3
import logging
import os
import re

# Remove all logs from pdfplumber pdfminer (pdfminer is a dependency of pdfplumber)
logging.getLogger("pdfminer").setLevel(logging.ERROR)

def process_pdf(book_name, pdf_path, output_dir, db):
    """
    Main function to process the PDF and extract Parts, Questions, and Articles.
    Args:
        pdf_path (str): Path to the PDF file.
        output_dir (str): Path to the output directory (e.g., 'web_app').
        db (sqlite3.Connection): SQLite database connection.
    """
    # Set up the base directory and SQLite database
    base_dir = os.path.join(output_dir, book_name)
    
    try:
        with pdfplumber.open(pdf_path) as pdf:
            for page in pdf.pages:
                # Extract text from the page
                text = page.extract_text()
                if not text:
                    continue
                
                # Process each line on the page
                lines = text.split('\n')
                for line in lines:
                    line = line.strip()
                    if not line:
                        continue

                    # Check for Question, and Article titles
                    question_match = re.match(r'^Questão\s(\d):\s(.+)$', line)
                    if question_match:
                        print(f"Found Question {question_match.group(1)}: {question_match.group(2)}")
                        continue

                    article_match = re.match(r'^Art\.\s(\d)\s.\s(.+)$', line)
                    if article_match:
                        print(f"Found Article {article_match.group(1)}: {article_match.group(2)}")
                        continue

    except FileNotFoundError:
        print(f"Error: PDF file '{pdf_path}' not found.")
    except Exception as e:
        print(f"Error processing PDF: {str(e)}")
    finally:
        # Clear memory and close the database connection
        db.close()
        print(f"Sections extracted to {output_dir}")

def init_dependencies(base_path, book_name):
    """
    Create the base output directory and set up the SQLite database.
    Args:
        base_path (str): Path to the output directory.
        book_name (str): Name of the book.
    Returns:
        sqlite3.Connection: Connection to the SQLite database.
    """
    # Create the base directory
    os.makedirs(base_path, exist_ok=True)
    os.makedirs(os.path.join(base_path, book_name), exist_ok=True)

    base_path = os.path.join(base_path, book_name)

    # Create output directories
    os.makedirs(os.path.join(base_path, 'questions'), exist_ok=True)
    os.makedirs(os.path.join(base_path, 'articles'), exist_ok=True)
    
    # Set up SQLite database
    db_path = os.path.join(os.path.dirname(base_path), f'{book_name}.db')
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    # Create articles table if it doesn't exist
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS articles (
            part TEXT,
            question_num INTEGER,
            article_num INTEGER,
            article_title TEXT,
            PRIMARY KEY (part, question_num, article_num)
        )
    ''')
    conn.commit()
    return conn
