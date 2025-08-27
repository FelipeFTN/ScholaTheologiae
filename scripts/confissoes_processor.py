import os
import re
import sqlite3

from utils import sqlite_connection, roman_to_int, save_content

def parse_markdown(markdown_content):
    lines = markdown_content.split('\n')

    parts = []
    content = {
        'part_title': '',
        'chapter_number': 0,
        'chapter_title': '',
        'content': []
    }

    for line in lines:
        # Check for part "LIVRO" detected
        if line.startswith('# LIVRO'):
            if content['part_title'] and content['chapter_title']:
                parts.append(content)
                content = {
                    'part_title': '',
                    'chapter_number': 0,
                    'chapter_title': '',
                    'content': []
                }
            content['part_title'] = line.replace('# ', '').strip().lower().replace(' ', '_')
            print("[-] Processing Part: " + str(content['part_title']))
            continue

        # Check for chapter "CAPÍTULO I..." detected
        elif line.startswith('## CAPÍTULO'):
            chapter_match = re.match(r'## CAPÍTULO (\w+)', line)
            if chapter_match:
                roman_numeral = chapter_match.group(1)
                content['chapter_number'] = roman_to_int(roman_numeral)

        # Check for chapter title detected
        elif line.startswith('### '):
            content['chapter_title'] = line.replace('### ', '').strip()
            print("[-] Processing Chapter: " + str(content['chapter_number']) + ": " + str(content['chapter_title']))
            if content['part_title'] and content['chapter_number']:
                parts.append(content)
                content = {
                    'part_title': content['part_title'],
                    'chapter_number': content['chapter_number'],
                    'chapter_title': '',
                    'content': []
                }

        # Add line to part content
        content['content'].append(line)

    if content['chapter_title']:
        parts.append(content)

    return parts

def insert_into_database(cursor, parts):
    for part in parts:
        part_title = part['part_title']
        chapter_number = part['chapter_number']
        chapter_title = part['chapter_title']

        # Insert chapter
        cursor.execute('''
            INSERT INTO confissoes (part_title, chapter_number, chapter_title)
            VALUES (?, ?, ?)
        ''', (part_title, chapter_number, chapter_title))

def main():
    db_path = 'data/confissoes/confissoes.db'
    if os.path.exists(db_path):
        os.remove(db_path)

    with open('data/confissoes/confissoes.md', 'r', encoding='utf-8') as file:
        markdown_content = file.read()

    # Create database & table
    conn, cursor = sqlite_connection(db_path, "confissoes")

    # Parse markdown content
    parts = parse_markdown(markdown_content)

    # Insert data into database
    insert_into_database(cursor, parts)

    # Save parts to markdown files
    for part in parts:
        save_content(part, 'confissoes')

    # Commit and close
    conn.commit()
    conn.close()
    print("[+] Data successfully processed into data/confissoes/...!")

if __name__ == "__main__":
    main()
