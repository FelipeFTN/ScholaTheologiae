from utils import init_dependencies, process_pdf

db = None

if __name__ == "__main__":
    book_name = "suma_teologica"
    pdf_path = "../books/suma_teologica_Vol_I.pdf"
    output_dir = "../articles"
    db = init_dependencies(output_dir, book_name)
    process_pdf(book_name, pdf_path, output_dir, db)
