import os
import re
import sqlite3

def init_dependencies(base_path, book_name):
    """
    Create the base output directory
    Args:
        base_path (str): Path to the output directory (e.g., 'web_app').
        book_name (str): Name of the book (e.g., 'summa_theologiae').
    Returns:
        sqlite3.Connection: Connection to the SQLite database.
        None
    """
    # Create the base directory
    os.makedirs(base_path, exist_ok=True)
    os.makedirs(os.path.join(base_path, book_name), exist_ok=True)

    base_path = os.path.join(base_path, book_name)

    db_path = f'data/{book_name}.db'
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    # Create summa_theologiae table with the question_title column if it doesn't exist
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS summa_theologiae (
            part TEXT,
            question_num INTEGER,
            question_title TEXT,
            article_num INTEGER,
            article_title TEXT,
            PRIMARY KEY (part, question_num, article_num)
        )
    ''')
    conn.commit()
    return conn


def save_db_metadata(db, part, question_num, question_title, article_num, article_title):
    """
    Update the SQLite database with metadata for an Article.
    Args:
        db (sqlite3.Connection): Connection to the SQLite database.
        part (str): The Part name.
        question_num (str): The Question number.
        question_title (str): The Question title.
        article_num (str): The Article number.
        article_title (str): The Article title.
    """
    cursor = db.cursor()
    cursor.execute('''
        INSERT OR REPLACE INTO summa_theologiae (part, question_num, question_title, article_num, article_title)
        VALUES (?, ?, ?, ?, ?)
    ''', (part, int(question_num), question_title, int(article_num), article_title))
    db.commit()

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
        'I': '1-prima_pars',
        'II': '2.1-prima_pars_secundae',
        'III': '2.2-secundae_secundae',
        'IV': '3-tertia_pars',
        'V': '4-supplementum'
    }
    part = part_mapping.get(volume, 'unknown_part')
    
    # Handle the appendix case (e.g., 'suma_teologica_Vol_V_Apendice.pdf')
    if 'Apendice' in pdf_filename:
        part = f"4.1-supplementum_appendix"
    
    return part


def save_question(base_dir, part, question_num, content):
    """
    Save a Question to a Markdown file.
    Args:
        base_dir (str): Base directory (e.g., 'data/summa_theologiae').
        part (str): The Part name (e.g., 'prima_pars').
        question_num (str): The Question number.
        content (list): List of lines for the Question content.
    """
    if not part or not question_num or not content:
        return
    
    # Create the questions directory for this Part
    question_dir = os.path.join(base_dir, part, 'questions')
    os.makedirs(question_dir, exist_ok=True)
    
    # Write the Question file
    filename = f"question_{question_num}.md"
    filepath = os.path.join(question_dir, filename)
    content_text = '\n'.join(content)

    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content_text)

def matches_invalid(line: str, patterns: list[str]) -> bool:
    for pattern in patterns:
        # If the pattern is purely alphanumeric (without dots or symbols), use word-boundary matching
        if re.match(r'^[\w]+$', pattern):
            if re.search(rf'\b{re.escape(pattern)}\b', line):
                # print(f"[DEBUG] Matched word pattern '{pattern}' in line: {line}")
                return True
        else:
            # For symbols or abbreviations (e.g., '(', 'Sent.'), do a simple substring match
            if pattern in line:
                # print(f"[DEBUG] Matched symbol pattern '{pattern}' in line: {line}")
                return True
    return False
