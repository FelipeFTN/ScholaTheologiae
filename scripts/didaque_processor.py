import os
import re
import sqlite3

from utils import sqlite_connection, roman_to_int, save_content, save_chapter_content

def parse_markdown(markdown_content):
    lines = markdown_content.split('\n')

    chapters = []
    current_content = {
        'part_title': '',
        'chapter_number': 0,
        'chapter_title': '',
        'chapter_content': [],
        'content': []
    }

    for line in lines:
        # Check for book part
        if line.startswith('## '):
            # Save previous content if exists
            if current_content['part_title'] and current_content['chapter_title']:
                save_chapter_content(current_content, 'didaque')
                save_content(current_content, 'didaque')
                chapters.append({
                    'part_title': current_content['part_title'],
                    'chapter_number': current_content['chapter_number'],
                    'chapter_title': current_content['chapter_title']
                })

            # Start new part
            current_content = {
                'part_title': line.replace('## ', '').strip().lower().replace(' ', '_'),
                'chapter_number': 0,
                'chapter_title': '',
                'chapter_content': [],
                'article_title': '',
                'article_number': 0,
                'content': []
            }
            print("[-] Processing Part: " + current_content['part_title'])
            current_content['content'].append(line)

        # Check for chapter
        elif line.startswith('### '):
            # Save previous chapter if exists
            if current_content['chapter_title']:
                if current_content['chapter_title'] == "Introdução":
                    print("    [-] Saving Introduction")
                save_chapter_content(current_content, 'didaque')
                chapters.append({
                    'part_title': current_content['part_title'],
                    'chapter_number': current_content['chapter_number'],
                    'chapter_title': current_content['chapter_title']
                })

            chapter_match = re.match(r'### Capítulo (.+) - (.+)', line)
            if chapter_match:
                roman_numeral = chapter_match.group(1)
                current_content['chapter_title'] = chapter_match.group(2).strip()
                current_content['chapter_number'] = roman_to_int(roman_numeral)
                current_content['chapter_content'] = [line]  # Start fresh chapter content
            # else it's introduction
            elif line.startswith('### Introdução'):
                current_content['chapter_title'] = 'Introdução'
                current_content['chapter_number'] = 0
                current_content['chapter_content'] = [line]
            current_content['content'].append(line)

        # Add line to content
        else: 
            current_content['chapter_content'].append(line)
            current_content['content'].append(line)

    # Save final content
    if current_content['chapter_title']:
        save_chapter_content(current_content, 'didaque')
        chapters.append({
            'part_title': current_content['part_title'],
            'chapter_number': current_content['chapter_number'],
            'chapter_title': current_content['chapter_title']
        })

    if current_content['part_title']:
        save_content(current_content, 'didaque')

    return chapters

def insert_into_database(cursor, chapters):
    for chapter in chapters:
        part_title = chapter['part_title']
        chapter_number = chapter['chapter_number']
        chapter_title = chapter['chapter_title']

        # Insert chapter
        cursor.execute('''
            INSERT INTO didaque (part_title, chapter_number, chapter_title)
            VALUES (?, ?, ?)
        ''', (part_title, chapter_number, chapter_title))

def main():
    db_path = 'data/didaque/didaque.db'
    if os.path.exists(db_path):
        os.remove(db_path)

    with open('data/didaque/didaque.md', 'r', encoding='utf-8') as file:
        markdown_content = file.read()

    # Create database & table
    conn, cursor = sqlite_connection(db_path, "didaque")

    # Parse markdown content
    chapters = parse_markdown(markdown_content)

    # Insert data into database
    insert_into_database(cursor, chapters)

    # Commit and close
    conn.commit()
    conn.close()
    print("[+] Data successfully processed into data/didaque/...!")

if __name__ == "__main__":
    main()
