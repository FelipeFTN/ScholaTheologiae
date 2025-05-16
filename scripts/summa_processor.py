import pdfplumber
import logging
import sys
import re
import os

from utils import map_volume_to_part, init_dependencies, save_question, matches_invalid, save_db_metadata

# Suppress warnings from pdfplumber and its dependency pdfminer to keep the output clean
logging.getLogger("pdfminer").setLevel(logging.ERROR)

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
