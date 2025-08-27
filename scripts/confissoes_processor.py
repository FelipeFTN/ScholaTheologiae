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
        # Check for part "LIVRO" detected
        if line.startswith('# LIVRO'):
            # Save previous content if exists
            if current_content['part_title'] and current_content['chapter_title']:
                save_chapter_content(current_content, 'confissoes')
                save_content(current_content, 'confissoes')
                chapters.append({
                    'part_title': current_content['part_title'],
                    'chapter_number': current_content['chapter_number'],
                    'chapter_title': current_content['chapter_title']
                })
            
            # Start new part
            current_content = {
                'part_title': line.replace('# ', '').strip().lower().replace(' ', '_'),
                'chapter_number': 0,
                'chapter_title': '',
                'chapter_content': [],
                'content': []
            }
            print("[-] Processing Part: " + current_content['part_title'])
            current_content['content'].append(line)

        # Check for chapter "CAPÍTULO I..." detected
        elif line.startswith('## CAPÍTULO'):
            # Save previous chapter if exists
            if current_content['chapter_title']:
                save_chapter_content(current_content, 'confissoes')
                chapters.append({
                    'part_title': current_content['part_title'],
                    'chapter_number': current_content['chapter_number'],
                    'chapter_title': current_content['chapter_title']
                })

            chapter_match = re.match(r'## CAPÍTULO (\w+)', line)
            if chapter_match:
                roman_numeral = chapter_match.group(1)
                current_content['chapter_number'] = roman_to_int(roman_numeral)
                current_content['chapter_content'] = [line]  # Start fresh chapter content
            current_content['content'].append(line)

        # Check for chapter title detected
        elif line.startswith('### '):
            current_content['chapter_title'] = line.replace('### ', '').strip()
            print("[-] Processing Chapter: " + str(current_content['chapter_number']) + ": " + current_content['chapter_title'])
            current_content['chapter_content'].append(line)
            current_content['content'].append(line)

        # Add line to appropriate content
        else: 
            current_content['chapter_content'].append(line)
            current_content['content'].append(line)

    # Save final content
    if current_content['chapter_title']:
        save_chapter_content(current_content, 'confissoes')
        chapters.append({
            'part_title': current_content['part_title'],
            'chapter_number': current_content['chapter_number'],
            'chapter_title': current_content['chapter_title']
        })
    if current_content['part_title']:
        save_content(current_content, 'confissoes')

    return chapters

def insert_into_database(cursor, chapters):
    for chapter in chapters:
        part_title = chapter['part_title']
        chapter_number = chapter['chapter_number']
        chapter_title = chapter['chapter_title']

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
    chapters = parse_markdown(markdown_content)

    # Insert data into database
    insert_into_database(cursor, chapters)

    # Commit and close
    conn.commit()
    conn.close()
    print("[+] Data successfully processed into data/confissoes/...!")

if __name__ == "__main__":
    main()
