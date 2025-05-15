import PyPDF2

# THIS TOOL WAS USING FOR BREAKING THE MASSIVE PDF INTO VOLUMES
def remove_first_pages(input_pdf_path, output_pdf_path, pages_to_remove=1):
    with open(input_pdf_path, "rb") as infile:
        reader = PyPDF2.PdfReader(infile)
        writer = PyPDF2.PdfWriter()

        total_pages = len(reader.pages)
        if pages_to_remove >= total_pages:
            raise ValueError("The PDF has fewer pages than the number to remove.")

        for i in range(pages_to_remove, total_pages):
            writer.add_page(reader.pages[i])

        with open(output_pdf_path, "wb") as outfile:
            writer.write(outfile)

def remove_last_pages(input_pdf_path, output_pdf_path, remove_pages_after=1):
    with open(input_pdf_path, "rb") as infile:
        reader = PyPDF2.PdfReader(infile)
        writer = PyPDF2.PdfWriter()

        total_pages = len(reader.pages)
        if remove_pages_after >= total_pages:
            raise ValueError("The PDF has fewer pages than the number to remove.")

        for i in range(remove_pages_after):
            writer.add_page(reader.pages[i])

        with open(output_pdf_path, "wb") as outfile:
            writer.write(outfile)

def remove_first_and_last_pages(input_pdf_path, output_pdf_path, first_pages_to_remove=1, pages_to_stop_removing=1):
    with open(input_pdf_path, "rb") as infile:
        reader = PyPDF2.PdfReader(infile)
        writer = PyPDF2.PdfWriter()

        if pages_to_stop_removing == 0:
            pages_to_stop_removing = len(reader.pages)

        for i in range(first_pages_to_remove, pages_to_stop_removing):
            writer.add_page(reader.pages[i])

        with open(output_pdf_path, "wb") as outfile:
            writer.write(outfile)

if __name__ == "__main__":
    book_name = "suma_teologica"
    remove_first_and_last_pages(f"../books/{book_name}.pdf", f"../books/{book_name}_Vol_V_Apendice.pdf", 4264, 0)
