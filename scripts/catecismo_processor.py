import os
import re
import sqlite3

def sqlite_connection(db_path):
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    cursor.execute('SELECT sqlite_version();')

    result = cursor.fetchall()
    print(f'[!] SQLite version is {result[0][0]}')

    cursor.execute('''
        CREATE TABLE IF NOT EXISTS catecismo_pio_x (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            part_title TEXT,
            chapter_number INTEGER NOT NULL,
            chapter_title TEXT NOT NULL,
            article_title TEXT,
            article_number INTEGER
        );
    ''')

    conn.commit()
    return conn, cursor

def insert_into_database(cursor, chapters):
    for item in chapters:
        # Insert items in the database
        cursor.execute('''
            INSERT INTO catecismo_pio_x (
                part_title, part_subtitle, chapter_number, chapter_title
            ) VALUES (?, ?, ?, ?);
        ''', (item['part_title'], item['part_subtitle'], item['chapter_number'], item['chapter_title']))

def ensure_directory_exists(part_title):
    # Replace spaces with underscores for directory name
    dir_name = "content/" + part_title.replace(' ', '_')
    os.makedirs(f'data/catecismo_pio_x/{dir_name}', exist_ok=True)
    return dir_name

def get_part_subtitle(text):
    for line in text:
        match = re.search(r'<aside>(.*?)</aside>', line)

        if match:
            return match.group(1)

    return ''

def save_content(content):
    # Define the output directory and file path
    dir_name = ensure_directory_exists(content['part'])
    output_path = f'data/catecismo_pio_x/{dir_name}/content.md'
    
    # Join content lines, preserving empty lines
    content_text = '\n'.join(content['content']).strip()
    
    # Write to file
    with open(output_path, 'w', encoding='utf-8') as file:
        file.write(content_text)

def save_chapter_content(content):
    # Define the output directory and file path
    dir_name = ensure_directory_exists(content['part'])
    chapter_number = content['chapter_number']
    output_path = f'data/catecismo_pio_x/{dir_name}/chapter_{chapter_number}.md'
    
    # Join chapter's content lines, preserving empty lines
    content_text = '\n'.join(content['chapter_content']).strip()
    
    # Write to file
    with open(output_path, 'w', encoding='utf-8') as file:
        file.write(content_text)

def parse_markdown(file_content):
    lines = file_content.split('\n')

    chapters = []
    content = {'part': '', 'content': [], 'chapter_content': [], 'chapter_title': '', 'chapter_number': -1} # -1 because we have a preliminary lesson


    for line in lines:
        # Check for part detected
        if line.startswith("# "):
            print("[-] Processing Part: " + line[2:])
            if content['part'] != '':
                # Part already inserted, so it's the end of a part and the beginning of another
                save_content(content)
                save_chapter_content(content)
                # Clear content for next stuff to be inserted
                content = {'part': '', 'content': [], 'chapter_content': [], 'chapter_title': '', 'chapter_number': 0}
            # Add new part to content
            content['part'] = line[2:].strip().replace(' ', '_').lower()
            content['part'] = content['part'].replace('ã', 'a').replace('ç', 'c').replace('ê', 'e') # Normalize special characters

        if line.startswith("## "):
            if content['chapter_title'] != '':
                save_chapter_content(content)
                # Clear chapter_content for next stuff to be inserted
                content['chapter_content'] = []

            content['chapter_number'] += 1
            content['chapter_title'] = line[3:]
            chapter_title = content['chapter_title'].replace('<em>', '').replace('</em>', '')
            chapters.append(
                {
                    'part_title': content['part'], 'part_subtitle': get_part_subtitle(content['content']),
                    'chapter_number': content['chapter_number'], 'chapter_title': chapter_title
                }
            )


        # include chapter content
        if content['chapter_title'] == '':
            # include it in content
            content['content'].append(line)
        else:
            content['chapter_content'].append(line)

    save_content(content)
    save_chapter_content(content)

    return chapters

def main():
    db_path = 'data/catecismo_pio_x/catecismo_pio_x.db'
    if os.path.exists(db_path):
        os.remove(db_path)

    with open('data/catecismo_pio_x/catecismo_pio_x.md', 'r', encoding='utf-8') as file:
        markdown_content = file.read()

    # Create database & table
    conn, cursor = sqlite_connection(db_path)
    
    # Parse markdown content
    chapters = parse_markdown(markdown_content)

    # Insert data into database
    insert_into_database(cursor, chapters)

    # Commit and close
    conn.commit()
    conn.close()
    print("[+] Data successfully processed into data/catecismo_pio_x/...!")

if __name__ == "__main__":
    main()
