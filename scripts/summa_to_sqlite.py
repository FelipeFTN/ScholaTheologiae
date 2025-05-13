import pdfplumber
import sqlite3
import logging
import re
import os

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

def init_dependencies(book_name):
    """
    Set up the SQLite database.
    Args:
        book_name (str): Name of the book (e.g., 'summa_teologica').
    Returns:
        sqlite3.Connection: Connection to the SQLite database.
    """
    # Set up SQLite database in the current directory
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

def update_article_metadata(db, part, question_num, question_title, article_num, article_title):
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

def process_pdf(book_name, pdf_path, db):
    """
    Main function to process the PDF and store Questions and summa_theologiae metadata in the SQLite database.
    Args:
        book_name (str): Name of the book (e.g., 'summa_teologica').
        pdf_path (str): Path to the PDF file.
        db (sqlite3.Connection): SQLite database connection.
    """
    # Determine the Part based on the PDF filename
    part = map_volume_to_part(os.path.basename(pdf_path))
    print(f"Processing PDF as Part: {part}")
    
    # Initialize variables
    current_question = None
    current_question_title = None
    
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
                        current_question = question_match.group(1)
                        current_question_title = question_match.group(2)
                        print(f"Found Question {current_question}: {current_question_title} on page {page.page_number}")
                        continue

                    # Check for a new Article
                    article_match = re.match(r'^Art\.\s(\d+)\s[—-]\s(.+)$', line)
                    if article_match and current_question and current_question_title:
                        article_num = article_match.group(1)
                        article_title = article_match.group(2)
                        # Update SQLite database with Article metadata
                        update_article_metadata(db, part, current_question, current_question_title, article_num, article_title)
                        print(f"Found Article {article_num}: {article_title}")
                        continue
        
        print("Metadata extracted to SQLite database")
    
    except FileNotFoundError:
        print(f"Error: PDF file '{pdf_path}' not found.")
    except Exception as e:
        print(f"Error processing PDF: {str(e)}")
    finally:
        # Clear memory
        current_question = None
        current_question_title = None

def main():
    """
    Main function to process all PDF volumes.
    """
    book_name = "summa_theologiae"
    pdf_volumes = [
        "suma_teologica_Vol_I.pdf",
        "suma_teologica_Vol_II.pdf",
        "suma_teologica_Vol_III.pdf",
        "suma_teologica_Vol_IV.pdf",
        "suma_teologica_Vol_V.pdf",
        "suma_teologica_Vol_V_Apendice.pdf"
    ]
    
    # Initialize dependencies
    db = init_dependencies(book_name)
    
    # Process each PDF volume
    for pdf_volume in pdf_volumes:
        pdf_path = os.path.join("../books", pdf_volume)
        print(f"\nProcessing {pdf_path}...")
        process_pdf(book_name, pdf_path, db)
    
    # Close the database connection
    db.close()

if __name__ == "__main__":
    main()
