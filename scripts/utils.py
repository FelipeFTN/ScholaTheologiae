import pdfplumber
import sqlite3
import logging
import os
import re

# Suppress warnings from pdfplumber and its dependency pdfminer to keep the output clean
logging.getLogger("pdfminer").setLevel(logging.ERROR)

def map_volume_to_part(pdf_filename):
    """
    Map a PDF volume filename to its corresponding Part in the Summa Theologica.
    Args:
        pdf_filename (str): Name of the PDF file (e.g., 'suma_teologica_Vol_I.pdf').
    Returns:
        str: Normalized Part name (e.g., 'prima_pars').
    """
    # Extract the volume part from the filename (e.g., 'Vol_I' from 'suma_teologica_Vol_I.pdf')
    volume_match = re.search(r'Vol_([IV]+)(?:_Apendice)?\.pdf$', pdf_filename, re.IGNORECASE)
    if not volume_match:
        return "unknown_part"
    
    volume = volume_match.group(1)
    # Map volumes to Parts
    part_mapping = {
        'I': 'prima_pars',
        'II': 'prima_pars_secundae',
        'III': 'secundae_secundae',
        'IV': 'tertia_pars',
        'V': 'supplementum'
    }
    part = part_mapping.get(volume, 'unknown_part')
    
    # Handle the appendix case (e.g., 'suma_teologica_Vol_V_Apendice.pdf')
    if 'Apendice' in pdf_filename:
        part = f"{part}_appendix"
    
    return part

def init_dependencies(base_path, book_name):
    """
    Create the base output directory and set up the SQLite database.
    Args:
        base_path (str): Path to the output directory (e.g., 'web_app').
        book_name (str): Name of the book (e.g., 'summa_teologica').
    Returns:
        sqlite3.Connection: Connection to the SQLite database.
    """
    # Create the base directory
    os.makedirs(base_path, exist_ok=True)
    os.makedirs(os.path.join(base_path, book_name), exist_ok=True)

    base_path = os.path.join(base_path, book_name)

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

def save_question(base_dir, part, question_num, content, metadata):
    """
    Save a Question to a Markdown file with Jekyll-compatible front matter.
    Args:
        base_dir (str): Base directory (e.g., 'web_app/summa_teologica').
        part (str): The Part name (e.g., 'prima_pars').
        question_num (str): The Question number.
        content (list): List of lines for the Question content.
        metadata (dict): Metadata for the Question (e.g., title, page_start).
    """
    if not part or not question_num or not content:
        return
    
    # Add Jekyll-specific metadata
    metadata['part'] = part
    metadata['layout'] = 'question'
    metadata['permalink'] = f'/{part}/question/{question_num}'
    
    # Create the questions directory for this Part
    question_dir = os.path.join(base_dir, part, 'questions')
    os.makedirs(question_dir, exist_ok=True)
    
    # Write the Question file
    filename = f"question_{question_num}.md"
    filepath = os.path.join(question_dir, filename)
    content_text = '\n'.join(content)
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write("---\n")
        for key, value in metadata.items():
            f.write(f"{key}: {value}\n")
        f.write(f"---\n\n{content_text}")
    print(f"Saved Question {question_num} to {filepath}")

def update_article_metadata(db, part, question_num, article_num, article_title):
    """
    Update the SQLite database with metadata for an Article.
    Args:
        db (sqlite3.Connection): Connection to the SQLite database.
        part (str): The Part name.
        question_num (str): The Question number.
        article_num (str): The Article number.
        article_title (str): The Article title.
    """
    cursor = db.cursor()
    cursor.execute('''
        INSERT OR REPLACE INTO articles (part, question_num, article_num, article_title)
        VALUES (?, ?, ?, ?)
    ''', (part, int(question_num), int(article_num), article_title))
    db.commit()

def process_pdf(book_name, pdf_path, output_dir, db):
    """
    Main function to process the PDF and extract Questions and Articles.
    Args:
        book_name (str): Name of the book (e.g., 'summa_teologica').
        pdf_path (str): Path to the PDF file.
        output_dir (str): Path to the output directory (e.g., 'web_app').
        db (sqlite3.Connection): SQLite database connection.
    """
    # Determine the Part based on the PDF filename
    part = map_volume_to_part(os.path.basename(pdf_path))
    print(f"Processing PDF as Part: {part}")
    
    # Set up the base directory
    base_dir = os.path.join(output_dir, book_name)
    
    # Initialize variables
    current_question = None
    question_content = []
    metadata = {}
    
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

                    # Check for a new Question
                    question_match = re.match(r'^Questão\s(\d+):\s(.+)$', line)
                    if question_match:
                        # Save the previous Question if it exists
                        if current_question:
                            save_question(base_dir, part, current_question, question_content, metadata)
                            # Clear memory
                            question_content.clear()
                            metadata.clear()
                        
                        # Start a new Question
                        current_question = question_match.group(1)
                        question_title = question_match.group(2)
                        question_content = [f"# Questão {current_question}: {question_title}"]
                        metadata = {
                            'question': current_question,
                            'title': question_title,
                            'page_start': page.page_number
                        }
                        print(f"Found Question {current_question}: {question_title} on page {page.page_number}")
                        continue

                    # Check for a new Article
                    article_match = re.match(r'^Art\.\s(\d+)\s[—-]\s(.+)$', line)
                    if article_match:
                        article_num = article_match.group(1)
                        article_title = article_match.group(2)
                        # Add Article as a subsection in the Question content
                        question_content.append(f"\n## Art. {article_num} — {article_title}")
                        # Update SQLite database with Article metadata
                        if current_question:
                            update_article_metadata(db, part, current_question, article_num, article_title)
                        print(f"Found Article {article_num}: {article_title}")
                        continue
                    
                    # Add the line to the current Question content
                    if current_question:
                        question_content.append(line)
        
        # Save the last Question
        if current_question:
            save_question(base_dir, part, current_question, question_content, metadata)
    
    except FileNotFoundError:
        print(f"Error: PDF file '{pdf_path}' not found.")
    except Exception as e:
        print(f"Error processing PDF: {str(e)}")
    finally:
        # Clear memory
        question_content.clear()
        metadata.clear()
        print(f"Sections extracted to {output_dir}")
