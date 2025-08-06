import pdfplumber
import logging
import sqlite3
import sys
import re
import os

# Suppress warnings from pdfplumber and its dependency pdfminer to keep the output clean
logging.getLogger("pdfminer").setLevel(logging.ERROR)

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

    db_path = f'{base_path}/{book_name}.db'
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    # Create summa_theologiae table with the question_title column if it doesn't exist
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS summa_theologiae (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            part_title TEXT,
            chapter_number INTEGER NOT NULL,
            chapter_title TEXT NOT NULL,
            article_title TEXT,
            article_number INTEGER
        );
    ''')
    conn.commit()
    return conn

def save_db_metadata(db, part, question_num, question_title, article_num, article_title):
    """
    Update the SQLite database with metadata for an Article.
    Args:
        db (sqlite3.Connection): Connection to the SQLite database.
        part_title (str): The Part name.
        chapter_number (str): The Question number.
        chapter_title (str): The Question title.
        article_number (str): The Article number.
        article_title (str): The Article title.
    """
    cursor = db.cursor()
    cursor.execute('''
        INSERT OR REPLACE INTO summa_theologiae (part_title, chapter_number, chapter_title, article_number, article_title)
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
        question_number (str): The Question number.
        content (list): List of lines for the Question content.
    """
    if not part or not question_num or not content:
        return
    
    # Create the questions directory for this Part
    question_dir = os.path.join(base_dir, part)
    os.makedirs(question_dir, exist_ok=True)
    
    # Write the Question file
    filename = f"chapter_{question_num}.md"
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

def process_pdf(book_name, pdf_path, output_dir, db):
    """
    Main function to process the PDF and extract Questions and Articles.
    Args:
        book_name (str): Name of the book (e.g., 'summa_theologiae').
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
    question_title = None
    question_content = []

    try:
        with pdfplumber.open(pdf_path) as pdf:
            for page in pdf.pages:
                # Extract text from the page
                text = page.extract_text()
                if not text:
                    continue
                
                # Process each line on the page
                lines = text.split('\n')
                i = 0
                while i < len(lines):
                    line = lines[i].strip()
                    if not line:
                        i += 1
                        continue

                    # Check for a new Question
                    question_match = re.match(r'^Questão\s(\d+):\s(.+)$', line)
                    if question_match:
                        # Save the previous Question if it exists
                        if current_question:
                            save_question(base_dir, part, current_question, question_content)
                            # Clear memory
                            question_content.clear()
                        
                        # Start a new Question
                        current_question = question_match.group(1)
                        question_title = question_match.group(2).replace(" - Traduzir", "")
                        question_content = [f"# Questão {current_question}: {question_title}"]
                        print(f"Question {current_question}: {question_title} on page {page.page_number}")
                        i += 1
                        continue

                    # Check for a new Article - First Line
                    used_lines = 1
                    if line.startswith("Art."):
                        invalid_lines_pattern = [
                            "discute-se", "Parece que",
                            "(", ")", "Sent.", "dist.", "Verit.", "De Trin.", "Cont. Gent."
                        ]
                        first_line = lines[i].strip().replace("único", "1")
                        second_line = lines[i + 1].strip() if i + 1 < len(lines) else ""
                        third_line = lines[i + 2].strip() if i + 2 < len(lines) else ""
                        fourth_line = lines[i + 3].strip() if i + 3 < len(lines) else ""

                        new_line = first_line

                        if not matches_invalid(second_line, invalid_lines_pattern):
                            new_line = f"{first_line} {second_line}"
                            used_lines = 2
                            if not matches_invalid(third_line, invalid_lines_pattern):
                                new_line = f"{first_line} {second_line} {third_line}"
                                used_lines = 3
                                if 0 < len(fourth_line.split()) < 5 and not matches_invalid(fourth_line, invalid_lines_pattern):
                                    new_line = f"{first_line} {second_line} {third_line} {fourth_line}"
                                    used_lines = 4
                        # Check if the new line does not ends with "." and add it
                        if not new_line.endswith("."):
                            new_line += "."

                        article_match = re.match(r'^Art\.\s(\d+).?\s[-─―–—]\s(.+\.)$', new_line)
                        if article_match:
                            article_num = article_match.group(1)
                            article_title = article_match.group(2)
                            # Add Article as a subsection in the Question content
                            question_content.append(f"\n## Art. {article_num} — {article_title}")
                            save_db_metadata(db, part, current_question, question_title, article_num, article_title)
                            i += used_lines
                            continue
                        else:
                            print(f"[ERROR] Article {new_line} could not be matched")
                    
                    # Add the line to the current Question content
                    # Breaking lines in the content
                    if current_question:
                        if re.match(r'^\d\.', line):
                            # This is a line that starts with a number, let's add a break line
                            line = '\n' + line
                        elif re.match(r'^\b[A-ZÀ-ÖØ-Ý]{2,}\b', line):
                            line = '\n' + line

                        question_content.append(line)
                    i += 1
        # Save the last Question
        if current_question:
            save_question(base_dir, part, current_question, question_content)

    except FileNotFoundError:
        print(f"Error: PDF file '{pdf_path}' not found.")
    except Exception as e:
        exc_type, exc_obj, exc_tb = sys.exc_info()
        fname = os.path.split(exc_tb.tb_frame.f_code.co_filename)[1]
        print(f"Error processing PDF: {str(e)}")
        print(exc_type, fname, exc_tb.tb_lineno)
    finally:
        # Clear memory
        question_content.clear()
        print(f"Sections extracted to {output_dir}")

def main():
    """
    Main function to process all PDF volumes.
    """
    book_name = "summa_theologiae"
    output_dir = "data"
    pdf_volumes = [
        "suma_teologica_Vol_I.pdf",
        "suma_teologica_Vol_II.pdf",
        "suma_teologica_Vol_III.pdf",
        "suma_teologica_Vol_IV.pdf",
        "suma_teologica_Vol_V.pdf",
        "suma_teologica_Vol_V_Apendice.pdf"
    ]
    
    # Initialize dependencies
    db = init_dependencies(output_dir, book_name)
    
    # Process each PDF volume
    for pdf_volume in pdf_volumes:
        pdf_path = os.path.join("../books", pdf_volume)
        print(f"\nProcessing {pdf_path}...")
        process_pdf(book_name, pdf_path, output_dir, db)

if __name__ == "__main__":
    main()
