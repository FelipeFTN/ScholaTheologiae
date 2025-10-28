import sqlite3
import os

# function to convert roman numeral into an integer
def roman_to_int(roman):

    # dictionary to store the roman-integer values
    map_symbols = {"I": 1, "V": 5, "X": 10, "L": 50, "C": 100, "D": 500, "M": 1000}

    ans = 0
    for i in range(len(roman)):
        if i > 0 and map_symbols[roman[i]] > map_symbols[roman[i - 1]]:
            ans += map_symbols[roman[i]] - 2 * map_symbols[roman[i - 1]]
        else:
            ans += map_symbols[roman[i]]
    return ans

def sqlite_connection(db_path, db_name):
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    cursor.execute('SELECT sqlite_version();')

    result = cursor.fetchall()
    print(f'[!] SQLite version is {result[0][0]}')

    cursor.execute(f'''
        CREATE TABLE IF NOT EXISTS {db_name} (
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

def ensure_directory_exists(part_title, content_name):
    # Replace spaces with underscores for directory name
    dir_name = "content/" + part_title.replace(' ', '_')
    os.makedirs(f'data/{content_name}/{dir_name}', exist_ok=True)
    return dir_name

def save_content(content, content_name):
    # Define the output directory and file path
    dir_name = ensure_directory_exists(content['part_title'], content_name)
    output_path = f'data/{content_name}/{dir_name}/content.md'
    
    # Join content lines, preserving empty lines
    content_text = '\n'.join(content['content']).strip()
    
    # Write to file
    with open(output_path, 'w', encoding='utf-8') as file:
        file.write(content_text)

def save_chapter_content(content, content_name):
    # Define the output directory and file path
    part_key = 'part' if 'part' in content else 'part_title'
    dir_name = ensure_directory_exists(content[part_key], content_name)
    chapter_number = content['chapter_number']
    output_path = f'data/{content_name}/{dir_name}/chapter_{chapter_number}.md'
    
    # Join chapter's content lines, preserving empty lines
    content_text = '\n'.join(content['chapter_content']).strip()
    
    # Write to file
    with open(output_path, 'w', encoding='utf-8') as file:
        file.write(content_text)
