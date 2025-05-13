from utils import init_dependencies, process_pdf
import os

# if __name__ == "__main__":
#     book_name = "suma_teologica"
#     pdf_path = "../books/suma_teologica_Vol_I.pdf"
#     output_dir = "../articles"
#     db = init_dependencies(output_dir, book_name)
#     process_pdf(book_name, pdf_path, output_dir, db)

def main():
    """
    Main function to process all PDF volumes.
    """
    book_name = "summa_teologica"
    output_dir = "../books"
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
    
    # Close the database connection
    db.close()

if __name__ == "__main__":
    main()
