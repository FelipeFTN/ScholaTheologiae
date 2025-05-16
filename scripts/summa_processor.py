import pdfplumber
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

def init_dependencies(base_path, book_name):
    """
    Create the base output directory
    Args:
        base_path (str): Path to the output directory (e.g., 'web_app').
        book_name (str): Name of the book (e.g., 'summa_theologiae').
    Returns:
        None
    """
    # Create the base directory
    os.makedirs(base_path, exist_ok=True)
    os.makedirs(os.path.join(base_path, book_name), exist_ok=True)

    base_path = os.path.join(base_path, book_name)

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

def process_pdf(book_name, pdf_path, output_dir):
    """
    Main function to process the PDF and extract Questions and Articles.
    Args:
        book_name (str): Name of the book (e.g., 'summa_theologiae').
        pdf_path (str): Path to the PDF file.
        output_dir (str): Path to the output directory (e.g., 'web_app').
    """
    # Determine the Part based on the PDF filename
    part = map_volume_to_part(os.path.basename(pdf_path))
    print(f"Processing PDF as Part: {part}")
    
    # Set up the base directory
    base_dir = os.path.join(output_dir, book_name)

    # Initialize variables
    current_question = None
    question_content = []

    try:
        with pdfplumber.open(pdf_path) as pdf:
            for page in pdf.pages:
                # Extract text from the page
                text = page.extract_text()
                if not text:
                    print(f"Empty text on page {page.page_number}")
                    continue
                
                # Process each line on the page
                lines = text.split('\n')
                i = 0
                while i < len(lines):
                    line = lines[i].strip()
                    if not line:
                        print(f"Empty line on page {page.page_number}")
                        print(f"Line: {line}")
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
                        question_title = question_match.group(2)
                        question_content = [f"# Questão {current_question}: {question_title}"]
                        print(f"Found Question {current_question}: {question_title} on page {page.page_number}")
                        i += 1
                        continue

                    # Check for a new Article - First Line
                    article_match = re.match(r'^Art\.\s(\d+).?\s[-─―–—]\s(.+\.)$', line)
                    ### There is a case where it is not matching the article because the article title
                    ### is not in the same line as the article number. So we need to check if the line starts with "Art."
                    if article_match:
                        article_num = article_match.group(1)
                        article_title = article_match.group(2)
                        # Add Article as a subsection in the Question content
                        question_content.append(f"\n## Art. {article_num} — {article_title}")
                        print(f"Found Article {article_num}: {article_title}")
                        i += 1
                        continue
                    # Found and Article but not matching the regex
                    elif line.startswith("Art."):
                        if "único" in line:
                            line = line.replace("único", "1")
                            question_content.append(line)   
                            i += 1
                            continue
                        # This is a special case where the article title is not in the same line as the article number
                        # We need to get the next line and add it to the article title
                        invalid_lines_pattern = [
                            "discute-se", "Parece que",
                            "(", ")", "Sent.", "dist.", "Verit.", "De Trin.", "Cont. Gent."
                        ]
                        used_lines = 1
                        first_line = lines[i].strip()
                        second_line = lines[i + 1].strip()
                        third_line = lines[i + 2].strip()
                        # set fourth line without getting out of index error
                        fourth_line = lines[i + 3].strip() if i + 3 < len(lines) else ""
                        new_line = first_line
                        if not matches_invalid(second_line, invalid_lines_pattern):
                            print(f"2. Article {second_line}")
                            new_line = f"{first_line} {second_line}"
                            used_lines = 2
                            if not matches_invalid(third_line, invalid_lines_pattern):
                                print(f"3. Article {third_line}")
                                new_line = f"{first_line} {second_line} {third_line}"
                                used_lines = 3
                                if 0 < len(fourth_line.split()) < 5 and not matches_invalid(fourth_line):
                                    print(f"4. Article {fourth_line}")
                                    new_line = f"{first_line} {second_line} {third_line} {fourth_line}"
                                    used_lines = 4
                        # Check if the new line does not ends with "." and add it
                        if not new_line.endswith("."):
                            new_line += "."
                        print(f"Generated Article {new_line}")
                        question_content.append(f"\n## {new_line}")
                        i += used_lines
                        continue
                    
                    # Add the line to the current Question content
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
        print(f"Error processing PDF: {str(e)}")
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
    init_dependencies(output_dir, book_name)
    
    # Process each PDF volume
    for pdf_volume in pdf_volumes:
        pdf_path = os.path.join("../books", pdf_volume)
        print(f"\nProcessing {pdf_path}...")
        process_pdf(book_name, pdf_path, output_dir)

if __name__ == "__main__":
    main()
